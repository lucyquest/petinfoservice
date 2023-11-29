package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
		return exitCode, err
	}
	defer func() {
		if err2 := testDB.Close(); err2 != nil {
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
			err = errors.Join(err, fmt.Errorf("error while running grpc service (%w)", err2))
		}
	}()

	// Setup gRPC client.
	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return exitCode, fmt.Errorf("could not dial test gRPC server error (%w)", err)
	}
	defer func() {
		if err2 := conn.Close(); err2 != nil {
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

func TestAddGetPets(t *testing.T) {
	testPets := []database.Pet{
		{
			Name:        "Lucy",
			DateOfBirth: time.Date(2007, 01, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "Miles",
			DateOfBirth: time.Date(2023, 01, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "Milo",
			DateOfBirth: time.Date(2023, 01, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for i := range testPets {
		addResp, err := client.Add(context.Background(), &petinfoproto.PetAddRequest{
			IdempotencyKey: uuid.NewString(),
			Pet: &petinfoproto.Pet{
				Name:        testPets[i].Name,
				DateOfBirth: timestamppb.New(testPets[i].DateOfBirth),
			},
		})
		if err != nil {
			t.Fatalf("could not add pet (%v)", err)
		}

		id, err := uuid.Parse(addResp.ID)
		if err != nil {
			t.Fatalf("expected a valid uuid for (%v) error (%v)", testPets[i].Name, err)
		}

		testPets[i].ID = id

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

		if !getResp.Pet.DateOfBirth.AsTime().Equal(testPets[i].DateOfBirth) {
			t.Fatalf("expected pet (%v) date of birth to be (%v) got (%v)", testPets[i].Name, testPets[i].DateOfBirth, getResp.Pet.DateOfBirth.AsTime())
		}
	}
}
