package service

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Olegnemlii/test123/internal/domain"
	"github.com/Olegnemlii/test123/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Генерация подписи
func (s *UserService) CreateSignature(ctx context.Context) (uuid.UUID, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Printf("error generating uuid: %v", err)
		return uuid.Nil, err
	}
	return newUUID, nil
}

// Верификация кода
func (s *UserService) VerifyCode(ctx context.Context, email string, code string) (bool, error) {
	storedCode, err := s.userRepo.GetVerificationCode(ctx, email)
	if err != nil {
		log.Printf("error getting verification code: %v", err)
		return false, err
	}

	if storedCode != code {
		log.Printf("verification code does not match")
		return false, fmt.Errorf("invalid verification code")
	}

	err = s.userRepo.DeleteVerificationCode(ctx, email)
	if err != nil {
		log.Printf("error deleting verification code: %v", err)
		return false, err
	}
	return true, nil
}

// Создание пользователя
func (s *UserService) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()

	return s.userRepo.CreateUser(ctx, user)
}

// Получение пользователя по ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}

// Получение пользователя по Email
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.GetUserByEmail(ctx, email)
}

// Обновление пользователя
func (s *UserService) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.userRepo.UpdateUser(ctx, user)
}

// Удаление пользователя
func (s *UserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.DeleteUser(ctx, id)
}

// Получение почты по подписи
func (s *UserService) GetEmailBySignature(ctx context.Context, signature uuid.UUID) (string, error) {
	return s.userRepo.GetEmailBySignature(ctx, signature)
}

// Обновление подписи пользователя
func (s *UserService) UpdateUserSignature(ctx context.Context, userID uuid.UUID, signature uuid.UUID) error {
	return s.userRepo.UpdateUserSignature(ctx, userID, signature)
}

// Хранение кода подтверждения
func (s *UserService) StoreRefreshToken(ctx context.Context, email, refreshToken string) error {
	return s.userRepo.StoreRefreshToken(ctx, email, refreshToken)
}

// Получение кода подтверждения
func (s *UserService) GetRefreshToken(ctx context.Context, email string) (string, error) {
	return s.userRepo.GetRefreshToken(ctx, email)
}

// Удаление кода подтверждения
func (s *UserService) DeleteRefreshToken(ctx context.Context, email string) error {
	return s.userRepo.DeleteRefreshToken(ctx, email)
}

// StoreVerificationCode - сохраняет код верификации
func (s *UserService) StoreVerificationCode(ctx context.Context, email string, code string) error {
	return s.userRepo.StoreVerificationCode(ctx, email, code)
}

// GenerateVerificationCode - генерирует код верификации
func (s *UserService) GenerateVerificationCode(ctx context.Context, email string) (string, error) {
	code := generateRandomCode(6) // Generate a 6-digit code
	err := s.userRepo.StoreVerificationCode(ctx, email, code)
	if err != nil {
		log.Printf("error storing verification code: %v", err)
		return "", err
	}
	return code, nil
}

// Функция для генерации случайной подписи
func generateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	const digits = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = digits[rand.Intn(len(digits))]
	}
	return string(code)
}

// GetVerificationCode
func (s *UserService) GetVerificationCode(ctx context.Context, email string) (string, error) {
	return s.userRepo.GetVerificationCode(ctx, email)
}

// DeleteVerificationCode
func (s *UserService) DeleteVerificationCode(ctx context.Context, email string) error {
	return s.userRepo.DeleteVerificationCode(ctx, email)
}
