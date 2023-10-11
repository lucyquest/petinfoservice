package service

import (
	"context"

	"github.com/lucyquest/petinfoservice/petinfoproto"
)

type petInfoDatabase interface {
}

type petInfoService struct {
	petinfoproto.UnimplementedPetInfoServiceServer

	db petInfoDatabase
}

func (p *petInfoService) Get(_ context.Context, _ *petinfoproto.PetGetRequest) (*petinfoproto.PetGetResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (p *petInfoService) GetMultiple(_ context.Context, _ *petinfoproto.PetGetMultipleRequest) (*petinfoproto.PetGetMultipleResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (p *petInfoService) EditName(_ context.Context, _ *petinfoproto.PetEditNameRequest) (*petinfoproto.PetEditNameResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (p *petInfoService) EditAge(_ context.Context, _ *petinfoproto.PetEditDateOfBirthRequest) (*petinfoproto.PetEditDateOfBirthResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (p *petInfoService) mustEmbedUnimplementedPetInfoServiceServer() {
	panic("not implemented") // TODO: Implement
}
