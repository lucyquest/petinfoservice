package service

import (
	"context"
	"database/sql"
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

func TestMain(m *testing.M) {
	// Force an non-zero exit code if a error happened anywhere
	var err error
	defer func() {
		if err != nil {
			os.Exit(1)
		}
	}()

	// Setup database and service for testing.
	// TODO: switch to a database per test?

	// Prepare a database to use in postgres
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresHost := os.Getenv("POSTGRES_HOST")

	connMainURI := fmt.Sprintf("postgres://%v:%v@%v?sslmode=disable", postgresUser, postgresPass, postgresHost)
	dbMain, err := sql.Open("postgres", connMainURI)
	if err != nil {
		log.Printf("could not open postgres database (%v) error (%v)", connMainURI, err)
		return
	}
	defer dbMain.Close()

	_, err = dbMain.Exec("CREATE DATABASE petinfoservice;")
	if err != nil {
		log.Printf("could not create petinfoservice database (%v)", err)
		return
	}

	// Connect to the prepared postgres database and assign it to a global variable
	connURI := fmt.Sprintf("postgres://%v:%v@%v/petinfoservice?sslmode=disable", postgresUser, postgresPass, postgresHost)
	db, err := sql.Open("postgres", connURI)
	if err != nil {
		log.Fatalf("could not open petinfoservice database (%v) error (%v)", connURI, err)
		return
	}

	schema, err := os.ReadFile("../database/schema.sql")
	if err != nil {
		log.Printf("could not read ../database/schema.sql error (%v)", err)
		return
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		log.Printf("could not create schema in database (%v)", err)
		return
	}

	// Startup gRPC service
	service := Service{
		Addr:     ":9999",
		Database: database.New(db),
	}
	serviceErr := make(chan error, 1)
	go func() {
		serviceErr <- service.Open()
	}()

	// Setup gRPC client
	conn, err := grpc.Dial("localhost:9999", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("could not dial test gRPC server error (%v)", err)
		return
	}
	defer conn.Close()

	client = petinfoproto.NewPetInfoServiceClient(conn)

	// Run tests.
	exitCode := m.Run()

	// Cleanup gRPC service and database.
	service.Close()
	err = <-serviceErr
	if err != nil {
		log.Printf("error while running grpc service (%v)", err)
		return
	}

	err = db.Close()
	if err != nil {
		log.Printf("could not close db connection. error: (%v)", err)
		return
	}

	_, err = dbMain.Exec("DROP DATABASE petinfoservice;")
	if err != nil {
		log.Printf("could not delete database during cleanup stage")
		return
	}

	err = dbMain.Close()
	if err != nil {
		log.Printf("could not close db connection. error: (%v)", err)
		return
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
			IdempotencyKey: "",
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
