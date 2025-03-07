package handler

import (
	"context"

	"test123/internal/service"
	"test123/pkg/pb"
)

// AuthHandler реализует gRPC сервис аутентификации.
type AuthHandler struct {
	pb.UnimplementedAuthServer
	userService *service.UserService
}

// NewAuthHandler создает новый экземпляр AuthHandler.
func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
	}
}

// Register обрабатывает регистрацию нового пользователя.
func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	user, err := h.userService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Signature: user.ID,
	}, nil
}

// Login обрабатывает вход пользователя.
func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	tokens, user, err := h.userService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		AccessToken:  &pb.Token{Data: tokens.AccessToken, ExpiresAt: tokens.AccessTokenExpiresAt},
		RefreshToken: &pb.Token{Data: tokens.RefreshToken, ExpiresAt: tokens.RefreshTokenExpiresAt},
		User:         &pb.User{Id: user.ID, Email: user.Email},
	}, nil
}

// VerifyCode обрабатывает проверку кода подтверждения.
func (h *AuthHandler) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.VerifyCodeResponse, error) {
	// Здесь можно добавить логику проверки кода, если она потребуется.
	return nil, nil
}

// RefreshTokens обновляет access и refresh токены.
func (h *AuthHandler) RefreshTokens(ctx context.Context, req *pb.RefreshTokensRequest) (*pb.RefreshTokensResponse, error) {
	tokens, user, err := h.userService.RefreshTokens(ctx, req.RefreshToken.Data)
	if err != nil {
		return nil, err
	}

	return &pb.RefreshTokensResponse{
		AccessToken:  &pb.Token{Data: tokens.AccessToken, ExpiresAt: tokens.AccessTokenExpiresAt},
		RefreshToken: &pb.Token{Data: tokens.RefreshToken, ExpiresAt: tokens.RefreshTokenExpiresAt},
		User:         &pb.User{Id: user.ID, Email: user.Email},
	}, nil
}

// GetMe возвращает данные текущего пользователя.
func (h *AuthHandler) GetMe(ctx context.Context, req *pb.GetMeRequest) (*pb.GetMeResponse, error) {
	user, err := h.userService.GetMe(ctx, req.AccessToken.Data)
	if err != nil {
		return nil, err
	}

	return &pb.GetMeResponse{
		User: &pb.User{Id: user.ID, Email: user.Email},
	}, nil
}

// LogOut обрабатывает выход пользователя из системы.
func (h *AuthHandler) LogOut(ctx context.Context, req *pb.LogOutRequest) (*pb.LogOutResponse, error) {
	err := h.userService.LogOut(ctx, req.AccessToken.Data)
	if err != nil {
		return nil, err
	}

	return &pb.LogOutResponse{Success: true}, nil
}
