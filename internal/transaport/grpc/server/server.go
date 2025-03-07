package server

import (
	"fmt"
	"log"
	"net"

	"test123/internal/service"
	"test123/internal/transport/grpc/handler"
	"test123/pkg/pb"

	"google.golang.org/grpc"
)

// GRPCServer представляет gRPC сервер.
type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer создает новый экземпляр gRPC сервера.
func NewGRPCServer(userService *service.UserService, port int) *GRPCServer {
	grpcServer := grpc.NewServer()
	authHandler := handler.NewAuthHandler(userService)

	pb.RegisterAuthServer(grpcServer, authHandler)

	return &GRPCServer{
		server: grpcServer,
		port:   port,
	}
}

// Run запускает gRPC сервер.
func (s *GRPCServer) Run() error {
	address := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
		return err
	}

	log.Printf("gRPC сервер запущен на %s", address)
	if err := s.server.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
		return err
	}

	return nil
}

// Stop останавливает gRPC сервер.
func (s *GRPCServer) Stop() {
	log.Println("Остановка gRPC сервера...")
	s.server.GracefulStop()
}
