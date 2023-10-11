package service

import (
	"crypto/tls"
	"log/slog"

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

	grpcServer *grpc.Server

	petInfoService *petInfoService
}

// TODO: readiness probe for k8s

func (s *Service) Open() error {
	var serverOptions []grpc.ServerOption

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

	s.grpcServer = grpc.NewServer(
		serverOptions...,
	)

	// Initalize petInfoService
	s.petInfoService = &petInfoService{
		db: nil,
	}

	petinfoproto.RegisterPetInfoServiceServer(s.grpcServer, s.petInfoService)

	return nil
}

func (s *Service) Close() {
	s.grpcServer.GracefulStop()
}
