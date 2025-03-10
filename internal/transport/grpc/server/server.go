package server

import (
	"fmt"
	"log"
	"net"

	"github.com/Olegnemlii/test123/internal/config"
	"github.com/Olegnemlii/test123/internal/transport/grpc/handler"
	"github.com/Olegnemlii/test123/pkg/pb"

	"google.golang.org/grpc"
)

// StartGRPCServer starts the gRPC server
func StartGRPCServer(cfg *config.Config, authHandler *handler.AuthHandler) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Port))
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return fmt.Errorf("failed to listen: %w", err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServer(s, authHandler)

	log.Printf("gRPC server listening on: %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Printf("failed to serve: %v", err)
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
