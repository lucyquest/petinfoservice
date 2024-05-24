package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type petInfoService struct {
	petinfoproto.UnimplementedPetInfoServiceServer

	queries database.Queries
	db      *sql.DB
}

func (p *petInfoService) Get(ctx context.Context, req *petinfoproto.PetGetRequest) (*petinfoproto.PetGetResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	pet, err := p.queries.GetPetByID(ctx, id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Error(codes.NotFound, "That pet does not exist")
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &petinfoproto.PetGetResponse{Pet: petAsProtoBufPet(pet)}, nil
}

func (p *petInfoService) GetMultiple(ctx context.Context, req *petinfoproto.PetGetMultipleRequest) (*petinfoproto.PetGetMultipleResponse, error) {
	// Convert protobuf ids to uuids
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for i := range req.IDs {
		id, err := uuid.Parse(req.IDs[i])
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid ID index (%v)", i)
		}
		ids = append(ids, id)
	}

	pets, err := p.queries.GetPetsByIDs(ctx, ids)
	switch {
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &petinfoproto.PetGetMultipleResponse{Pets: petsAsProtoBufPets(pets)}, nil
}

func (p *petInfoService) UpdateName(ctx context.Context, req *petinfoproto.PetUpdateNameRequest) (*petinfoproto.PetUpdateNameResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is empty")
	}

	err = p.queries.UpdatePetName(ctx, database.UpdatePetNameParams{
		ID:   id,
		Name: req.Name,
	})
	switch {
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return nil, nil
}

func (p *petInfoService) UpdateDateOfBirth(ctx context.Context, req *petinfoproto.PetUpdateDateOfBirthRequest) (*petinfoproto.PetUpdateDateOfBirthResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	if req.DateOfBirth.CheckValid() != nil {
		return nil, status.Error(codes.InvalidArgument, "date of birth not valid")
	}

	err = p.queries.UpdatePetDateOfBirth(ctx, database.UpdatePetDateOfBirthParams{
		ID:          id,
		DateOfBirth: req.DateOfBirth.AsTime(),
	})
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Error(codes.NotFound, "That pet does not exist")
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &petinfoproto.PetUpdateDateOfBirthResponse{}, nil
}

func (p *petInfoService) Add(ctx context.Context, req *petinfoproto.PetAddRequest) (*petinfoproto.PetAddResponse, error) {
	if req.IdempotencyKey == "" {
		return nil, status.Error(codes.InvalidArgument, "Please provide a IdempotencyKey")
	}

	if req.Pet == nil {
		return nil, status.Error(codes.InvalidArgument, "Pet not provided")
	}

	if req.Pet.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Must provide a Name")
	}

	if req.Pet.DateOfBirth == nil {
		return nil, status.Error(codes.InvalidArgument, "Must provide a Date Of Birth")
	}

	if req.Pet.ID != "" {
		return nil, status.Error(codes.InvalidArgument, "Pet ID must be empty")
	}

	requestProtobufBytes, err := proto.Marshal(req)
	if err != nil {
		slog.Error("Unknown error while marshaling request protobuf to bytes",
			"error", err.Error(),
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	tx, err := p.db.Begin()
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}
	defer func() {
		if err2 := tx.Rollback(); err2 != nil && !errors.Is(err2, sql.ErrTxDone) {
			slog.Error("Unknown error from database while rolling back transaction",
				"error", err2.Error(),
				"method", petinfoproto.PetInfoService_Add_FullMethodName,
			)
		}
	}()

	qtx := p.queries.WithTx(tx)

	userID := ctx.Value(UserID{}).(uuid.UUID)

	// Handle if idempotency key exists in the Database
	idemResponse, err := qtx.GetIdempotencyEntry(ctx, database.GetIdempotencyEntryParams{
		UserID:     userID,
		Key:        req.IdempotencyKey,
		MethodPath: petinfoproto.PetInfoService_Add_FullMethodName,
		Request:    requestProtobufBytes,
	})

	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		slog.Error("Unknown error from database while getting idempotency key",
			"error", err.Error(),
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	default:
		response := &petinfoproto.PetAddResponse{}
		if err = proto.Unmarshal(idemResponse, response); err != nil {
			slog.Error("Unknown error while unmarshaling protobuf message from idempotency entry",
				"error", err,
				"method", petinfoproto.PetInfoService_Add_FullMethodName,
			)
			return nil, status.Error(codes.Internal, "Internal error occurred")
		}
		return response, nil
	}

	petID, err := qtx.AddPet(ctx, database.AddPetParams{
		Name:        req.Pet.Name,
		DateOfBirth: req.Pet.DateOfBirth.AsTime(),
	})
	switch {
	case err != nil:
		slog.Error("Unknown error from database while adding pet",
			"error", err,
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	response := &petinfoproto.PetAddResponse{
		ID: petID.String(),
	}

	responseProtoBytes, err := proto.Marshal(response)
	if err != nil {
		slog.Error("Unknown error while marshaling responseProtobuf",
			"error", err,
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	err = qtx.AddIdempotencyEntry(ctx, database.AddIdempotencyEntryParams{
		UserID:     userID,
		Key:        req.IdempotencyKey,
		MethodPath: petinfoproto.PetInfoService_Add_FullMethodName,
		Request:    requestProtobufBytes,
		Response:   responseProtoBytes,
	})
	if err != nil {
		slog.Error("Unknown error from database while adding idempotency entry",
			"error", err,
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	if err = tx.Commit(); err != nil {
		slog.Error("Unknown error from database committing transaction",
			"error", err.Error(),
			"method", petinfoproto.PetInfoService_Add_FullMethodName,
		)
	}

	return response, nil
}

func (p *petInfoService) Remove(ctx context.Context, req *petinfoproto.PetRemoveRequest) (*petinfoproto.PetRemoveResponse, error) {
	return nil, nil
}
