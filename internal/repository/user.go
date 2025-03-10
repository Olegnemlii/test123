package repository

import (
	"context"

	"github.com/Olegnemlii/test123/internal/domain" // Замените 'insta' на имя вашего модуля

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetEmailBySignature(ctx context.Context, signature uuid.UUID) (string, error)
	UpdateUserSignature(ctx context.Context, userID uuid.UUID, signature uuid.UUID) error
	StoreVerificationCode(ctx context.Context, email string, code string) error
	GetVerificationCode(ctx context.Context, email string) (string, error)
	DeleteVerificationCode(ctx context.Context, email string) error
	StoreRefreshToken(ctx context.Context, email string, refreshToken string) error
	GetRefreshToken(ctx context.Context, email string) (string, error)
	DeleteRefreshToken(ctx context.Context, email string) error
	// Добавьте другие методы, которые вам нужны для работы с User
}
