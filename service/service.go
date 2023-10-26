package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type petInfoService struct {
	petinfoproto.UnimplementedPetInfoServiceServer

	db database.Queries
}

func (p *petInfoService) Get(ctx context.Context, req *petinfoproto.PetGetRequest) (*petinfoproto.PetGetResponse, error) {
	id, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ID")
	}

	pet, err := p.db.GetPetByID(ctx, id)
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

	pets, err := p.db.GetPetsByIDs(ctx, ids)
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

	err = p.db.UpdatePetName(ctx, database.UpdatePetNameParams{
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

	err = p.db.UpdatePetDateOfBirth(ctx, database.UpdatePetDateOfBirthParams{
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

func (p *petInfoService) Add(_ context.Context, _ *petinfoproto.PetAddRequest) (*petinfoproto.PetAddResponse, error) {
	panic("not implemented") // TODO: Implement
}
