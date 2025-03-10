package handler

import (
	"context"
	"log"

	"github.com/Olegnemlii/test123/internal/config"
	"github.com/Olegnemlii/test123/internal/domain"
	"github.com/Olegnemlii/test123/internal/service"
	"github.com/Olegnemlii/test123/pkg/pb"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authService service.UserService
	cfg         config.Config
	pb.UnimplementedAuthServer
}

func NewAuthHandler(authService service.UserService, cfg config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

// Регистрация пользователя
func (s *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	user := &domain.User{
		Email:       email,
		Password:    password,
		IsConfirmed: false,
	}

	createdUser, err := s.authService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("error creating user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	// Generate verification code and send it (omitted for brevity)
	code, err := s.authService.GenerateVerificationCode(ctx, email)
	if err != nil {
		log.Printf("error generating verification code: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate verification code")
	}

	log.Printf("Verification code for %s: %s", email, code)

	//TODO: create MailopostService and call it here for sending email
	//mailopostService := service.NewMailopostService(s.cfg.MailopostApiKey, s.cfg.MailopostURL)
	//err = mailopostService.SendVerificationEmail(email, code)
	//if err != nil {
	//	log.Printf("error sending verification email: %v", err)
	//	return nil, status.Errorf(codes.Internal, "failed to send verification email")
	//}

	signature, err := s.authService.CreateSignature(ctx)

	if err != nil {
		log.Printf("error creating signature: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create signature")
	}

	err = s.authService.UpdateUserSignature(ctx, createdUser.ID, signature)
	if err != nil {
		log.Printf("error updating user signature: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update user signature")
	}

	return &pb.RegisterResponse{Signature: signature.String()}, nil
}

// Подтверждение почты
func (s *AuthHandler) ConfirmEmail(ctx context.Context, req *pb.ConfirmEmailRequest) (*pb.ConfirmEmailResponse, error) {
	email := req.GetEmail()
	code := req.GetCode()

	if email == "" || code == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and code are required")
	}

	isValid, err := s.authService.VerifyCode(ctx, email, code)
	if err != nil {
		log.Printf("error verifying code: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to verify code")
	}

	if !isValid {
		return nil, status.Errorf(codes.InvalidArgument, "invalid code")
	}

	return &pb.ConfirmEmailResponse{Success: true}, nil
}

// Авторизация пользователя
func (s *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if email == "" || password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and password are required")
	}

	user, err := s.authService.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("error getting user: %v", err)
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	//TODO:: реализовать сравнение паролей
	if user.Password != password {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	// Store refresh token
	refreshToken := uuid.New().String()
	err = s.authService.StoreRefreshToken(ctx, email, refreshToken)
	if err != nil {
		log.Printf("error storing refresh token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to store refresh token")
	}

	return &pb.LoginResponse{
		AccessToken:  "accessToken",
		RefreshToken: refreshToken,
	}, nil
}

// Обновление токена
func (s *AuthHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	email := req.GetEmail()
	refreshToken := req.GetRefreshToken()

	if email == "" || refreshToken == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email and refresh token are required")
	}

	storedRefreshToken, err := s.authService.GetRefreshToken(ctx, email)

	if err != nil {
		log.Printf("error while get refresh token %v", err)
		return nil, status.Errorf(codes.Internal, "error while get refresh token")
	}

	if storedRefreshToken != refreshToken {
		return nil, status.Errorf(codes.Unauthenticated, "invalid refresh token")
	}

	newAccessToken := uuid.New().String()

	return &pb.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}
