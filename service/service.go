package service

import (
	"context"
	"database/sql"
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

	err = p.queries.UpdatePetDateOfBirth(ctx, database.UpdatePetDateOfBirthParams{
		ID:          id,
		DateOfBirth: req.DateOfBirth.AsTime(),
	})
	switch {
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return nil, nil
}

func (p *petInfoService) Add(ctx context.Context, req *petinfoproto.PetAddRequest) (*petinfoproto.PetAddResponse, error) {
	id, err := p.db.AddPet(ctx, database.AddPetParams{
		Name:        req.Pet.Name,
		DateOfBirth: req.Pet.DateOfBirth.AsTime(),
	})
	switch {
	case err != nil:
		slog.Error("Unknown error from database", "error", err)
		return nil, status.Error(codes.Internal, "Internal error occurred")
	}

	return &petinfoproto.PetAddResponse{ID: id.String()}, nil
}
