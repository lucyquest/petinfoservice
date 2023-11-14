package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log/slog"
	"net"

	"github.com/google/uuid"
	"github.com/lucyquest/petinfoservice/database"
	"github.com/lucyquest/petinfoservice/petinfoproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Service struct {
	// Configuration Options

	// Addr to listen on
	Addr string
	// Certificate to use if using TLS
	Certificate *tls.Certificate
	// Queries to use
	Queries *database.Queries
	// SQL DB
	DB *sql.DB

	grpcServer *grpc.Server

	petInfoService *petInfoService
}

type UserID struct{}

func (s *Service) authenticationUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	return handler(context.WithValue(ctx, UserID{}, uuid.New()), req)
}

// TODO: readiness probe for k8s

func (s *Service) Open() error {
	var serverOptions []grpc.ServerOption

	// TODO: add client certificate support

	// Decide if we should use TLS
	if s.Certificate == nil {
		slog.Info("Using default server with no TLS")
	} else {
		slog.Info("Using certificate")

		serverOptions = append(serverOptions,
			grpc.Creds(credentials.NewTLS(&tls.Config{
				Certificates: []tls.Certificate{*s.Certificate},
				ClientAuth:   tls.NoClientCert,
				MinVersion:   tls.VersionTLS13,
			})),
		)
	}

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("could not open tcp socket (%v) error (%w)", s.Addr, err)
	}
	s.Addr = l.Addr().String()

	// Set max grpc message we can receive to 1 MiB.
	serverOptions = append(serverOptions, grpc.MaxRecvMsgSize(1024*1024))
	serverOptions = append(serverOptions, grpc.UnaryInterceptor(s.authenticationUnaryInterceptor))

	s.grpcServer = grpc.NewServer(
		serverOptions...,
	)

	// Initalize petInfoService
	s.petInfoService = &petInfoService{queries: *s.Queries, db: s.DB}

	petinfoproto.RegisterPetInfoServiceServer(s.grpcServer, s.petInfoService)

	return s.grpcServer.Serve(l)
}

func (s *Service) Close() {
	// Protect against a nil pointer if Close() is called and server was never setup.
	// This might happen if the Close() is in a defer and we returned before the grpc.NewServer was assigned.
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
}
