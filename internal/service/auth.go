package service

import (
	"context"
	"log"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"insta/auth/database"
	authpb "insta/auth/pkg/pb"
)

type AuthService struct {
	DB *database.Database
}

func NewAuthService(db *database.Database) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	// Здесь будет запрос к БД для сохранения пользователя (пока не добавлен)
	return &authpb.RegisterResponse{Success: true}, nil
}
