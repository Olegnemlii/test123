package main

import (
	"context"
	"log"
	"time"

	pb "test123/pkg/pb"

	"google.golang.org/grpc"
)

const serverAddr = "localhost:50051"

func main() {
	// Устанавливаем соединение с gRPC-сервером
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Ошибка подключения к серверу: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Пример запроса регистрации пользователя
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := client.Register(ctx, req)
	if err != nil {
		log.Fatalf("Ошибка регистрации: %v", err)
	}

	log.Printf("Успешная регистрация! Подпись: %s", resp.Signature)
}
