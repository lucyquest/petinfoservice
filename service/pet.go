package service

import (
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func petAsProtoBufPet(p database.Pet) *petinfoproto.Pet {
	return &petinfoproto.Pet{
		ID:          p.ID.String(),
		Name:        p.Name,
		DateOfBirth: timestamppb.New(p.DateOfBirth),
	}
}

func petsAsProtoBufPets(pets []database.Pet) []*petinfoproto.Pet {
	pbPets := make([]*petinfoproto.Pet, 0, len(pets))

	for i := range pets {
		pbPets = append(pbPets, petAsProtoBufPet(pets[i]))
	}

	return pbPets
}
