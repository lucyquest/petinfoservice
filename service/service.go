package service

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"

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

	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("could not open tcp socket (%v) error (%w)", s.Addr, err)
	}

	s.grpcServer = grpc.NewServer(
		serverOptions...,
	)

	// Initalize petInfoService
	s.petInfoService = &petInfoService{
		db: nil,
	}

	petinfoproto.RegisterPetInfoServiceServer(s.grpcServer, s.petInfoService)

	return s.grpcServer.Serve(l)
}

func (s *Service) Close() {
	s.grpcServer.GracefulStop()
}
