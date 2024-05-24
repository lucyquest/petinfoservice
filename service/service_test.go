package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	client petinfoproto.PetInfoServiceClient
)

func testMain(m *testing.M) (exitCode int, err error) {
	testDB := testDB{
		User: os.Getenv("POSTGRES_USER"),
		Pass: os.Getenv("POSTGRES_PASS"),
		Host: os.Getenv("POSTGRES_HOST"),
	}

	err = testDB.Open()
	if err != nil {
		return 1, err
	}
	defer func() {
		if err2 := testDB.Close(); err2 != nil {
			exitCode = 1
			err = errors.Join(err, fmt.Errorf("could not close testDB error (%w)", err2))
		}
	}()

	// Startup gRPC service.
	service := NewService(
		":9999",
		nil,
		database.New(testDB.DB),
		testDB.DB,
	)

	serviceErr := make(chan error, 1)
	go func() {
		serviceErr <- service.Open()
	}()

	defer func() {
		service.Close()

		err2 := <-serviceErr
		if err2 != nil {
			exitCode = 1
			err = errors.Join(err, fmt.Errorf("error while running grpc service (%w)", err2))
		}
	}()

	// Setup gRPC client.
	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 1, fmt.Errorf("could not dial test gRPC server error (%w)", err)
	}
	defer func() {
		if err2 := conn.Close(); err2 != nil {
			exitCode = 1
			err = errors.Join(err, fmt.Errorf("could not close grpc client connection (%w)", err2))
		}
	}()

	client = petinfoproto.NewPetInfoServiceClient(conn)

	return m.Run(), nil
}

func TestMain(m *testing.M) {
	exitCode, err := testMain(m)
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitCode)
}

func TestGetPetsInvalidUUID(t *testing.T) {
	// Create resp for "no error" test
	resp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
		IdempotencyKey: t.Name(),
		Pet: &petinfoproto.Pet{
			Name:        t.Name(),
			DateOfBirth: timestamppb.Now(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		// TODO: delete id
	}()

	tests := []struct {
		name string
		id   string
		code codes.Code
	}{
		{
			name: "no error",
			id:   resp.ID,
			code: codes.OK,
		},
		{
			name: "empty id",
			id:   "",
			code: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Get(context.Background(), &petinfoproto.PetGetRequest{ID: tt.id})
			if status.Code(err) != tt.code {
				t.Errorf("expected error to be status code (%v) got error (%v:%T)", tt.code, err, err)
			}
		})
	}
}

func getTestPets() []*petinfoproto.Pet {
	return []*petinfoproto.Pet{
		{
			Name:        "Lucy",
			DateOfBirth: timestamppb.New(time.Date(2007, 01, 10, 0, 0, 0, 0, time.UTC)),
		},
		{
			Name:        "Miles",
			DateOfBirth: timestamppb.New(time.Date(2023, 01, 20, 0, 0, 0, 0, time.UTC)),
		},
		{
			Name:        "Milo",
			DateOfBirth: timestamppb.New(time.Date(2023, 01, 20, 0, 0, 0, 0, time.UTC)),
		},
	}
}

func TestGetMultiple(t *testing.T) {
	t.Run("golden path", func(t *testing.T) {
		testPets := getTestPets()
		ids := make([]string, len(testPets))
		for i := range testPets {
			addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
				IdempotencyKey: testPets[i].Name + t.Name(),
				Pet:            testPets[i],
			})
			if err != nil {
				t.Fatalf("could not add pet (%v:T)", err)
			}

			ids[i] = addResp.ID
			testPets[i].ID = addResp.ID
			// TODO: cleanup after ourselves (or implement per a per test DB)
		}

		resp, err := client.GetMultiple(context.Background(), &petinfoproto.PetGetMultipleRequest{IDs: ids})
		if err != nil {
			t.Fatalf("could not get multiple pets error (%v)", err)
		}

		if len(resp.Pets) != len(testPets) {
			t.Fatalf("expected len(resp.Pets)(%v) to be len(testPets)(%v)", len(resp.Pets), len(testPets))
		}

		// do a search for the ID, since there might not be a guarantee results are ordered
		for _, p := range resp.Pets {
			i := slices.IndexFunc(testPets, func(pp *petinfoproto.Pet) bool {
				return p.ID == pp.ID
			})

			if !proto.Equal(testPets[i], p) {
				t.Fatalf("expected testPets(%v) to be resp.Pets(%v)", testPets[i], p)
			}
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := client.GetMultiple(context.Background(), &petinfoproto.PetGetMultipleRequest{
			IDs: []string{"invalid-uuid", "invalid-uuid2"},
		})

		if status.Code(err) != codes.InvalidArgument {
			t.Fatalf("expected status.Code(InvalidArgument) got (%v:%T)", err, err)
		}
	})
}

func TestUpdateDateOfBirth(t *testing.T) {
	testPets := getTestPets()
	for _, pet := range testPets {
		addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
			IdempotencyKey: t.Name() + pet.Name,
			Pet:            pet,
		})
		if err != nil {
			t.Fatalf("could not add pet error (%v)", err)
		}

		pet.ID = addResp.ID

		timeToSetTo := time.Now()
		_, err = client.UpdateDateOfBirth(context.Background(), &petinfoproto.PetUpdateDateOfBirthRequest{
			ID:          pet.ID,
			DateOfBirth: timestamppb.New(timeToSetTo),
		})
		if err != nil {
			t.Fatalf("could not update date of birth (%v)", err)
		}

		getPet, err := client.Get(context.Background(), &petinfoproto.PetGetRequest{ID: pet.ID})
		if err != nil {
			t.Fatalf("could not get pet (%v)", err)
		}

		pet.DateOfBirth = timestamppb.New(timeToSetTo)
		if !proto.Equal(pet, getPet.Pet) {
			t.Fatalf("expected getPet.Pet(%v) to equal pet(%v)", getPet.Pet, pet)
		}
	}
}

func TestAddGetPets(t *testing.T) {
	testPets := getTestPets()
	ids := make([]string, len(testPets))
	for i := range testPets {
		addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
			IdempotencyKey: testPets[i].Name + t.Name(),
			Pet: &petinfoproto.Pet{
				Name:        testPets[i].Name,
				DateOfBirth: testPets[i].DateOfBirth,
			},
		})
		if err != nil {
			t.Fatalf("could not add pet (%v)", err)
		}

		id, err := uuid.Parse(addResp.ID)
		if err != nil {
			t.Fatalf("expected a valid uuid for (%v) error (%v)", testPets[i].Name, err)
		}

		getResp, err := client.Get(context.Background(), &petinfoproto.PetGetRequest{ID: id.String()})
		if err != nil {
			t.Fatalf("could not get pet (%v) error (%v)", testPets[i].Name, err)
		}

		if getResp.Pet.ID != id.String() {
			t.Fatalf("expected pet (%v) id to be (%v) got (%v)", testPets[i].Name, id.String(), getResp.Pet.ID)
		}

		if getResp.Pet.Name != testPets[i].Name {
			t.Fatalf("expected pet (%v) name to be (%v) got (%v)", testPets[i].Name, testPets[i].Name, getResp.Pet.Name)
		}

		if !getResp.Pet.DateOfBirth.AsTime().Equal(testPets[i].DateOfBirth.AsTime()) {
			t.Fatalf("expected pet (%v) date of birth to be (%v) got (%v)", testPets[i].Name, testPets[i].DateOfBirth, getResp.Pet.DateOfBirth.AsTime())
		}

		ids[i] = id.String()
	}

	for i := range testPets {
		addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
			IdempotencyKey: testPets[i].Name + t.Name(),
			Pet: &petinfoproto.Pet{
				Name:        testPets[i].Name,
				DateOfBirth: testPets[i].DateOfBirth,
			},
		})
		if err != nil {
			t.Fatalf("could not add pet (%v)", err)
		}

		if ids[i] != addResp.ID {
			t.Fatalf("expected idempotent id (%v) to be (%v)", ids[i], addResp.ID)
		}
	}
}

func TestDeletePet(t *testing.T) {
	testPets := getTestPets()
	for _, pet := range testPets {
		addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
			IdempotencyKey: pet.Name + t.Name(),
			Pet:            pet,
		})
		if err != nil {
			t.Fatalf("could not add pet (%v)", err)
		}

		_, err = client.Get(context.Background(), &petinfoproto.PetGetRequest{
			ID: addResp.ID,
		})
		if err != nil {
			t.Fatalf("could not get pet (%v)", err)
		}

		_, err = client.Remove(context.Background(), &petinfoproto.PetRemoveRequest{
			IdempotencyKey: pet.Name + t.Name(),
			Id:             pet.ID,
		})
		if err != nil {
			t.Fatalf("could not delete pet (%v) error (%v)", pet.ID, err)
		}
	}
}
