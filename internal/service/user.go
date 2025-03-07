package service

import (
	"context"
	"errors"
	"time"

	"test123/internal/domain"
	"test123/internal/repository/postgres"
	"test123/internal/repository/redis"
	"test123/pkg/hash"
	"test123/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// UserService определяет методы работы с пользователями.
type UserService struct {
	postgresRepo *postgres.UserRepository
	redisRepo    *redis.UserRepository
	tokenManager *token.Manager
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(pgRepo *postgres.UserRepository, redisRepo *redis.UserRepository, tokenManager *token.Manager) *UserService {
	return &UserService{
		postgresRepo: pgRepo,
		redisRepo:    redisRepo,
		tokenManager: tokenManager,
	}
}

// Register регистрирует нового пользователя.
func (s *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Проверяем, существует ли пользователь
	existingUser, _ := s.postgresRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("пользователь с таким email уже существует")
	}

	// Хэшируем пароль
	hashedPassword, err := hash.Generate(password)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
	}

	// Сохраняем в БД
	createdUser, err := s.postgresRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Кэшируем пользователя в Redis
	err = s.redisRepo.SaveUser(ctx, createdUser, time.Hour)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// Login аутентифицирует пользователя и выдает токены.
func (s *UserService) Login(ctx context.Context, email, password string) (*token.TokenPair, *domain.User, error) {
	// Получаем пользователя из БД
	user, err := s.postgresRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, nil, errors.New("неверный email или пароль")
	}

	// Проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, nil, errors.New("неверный email или пароль")
	}

	// Генерируем токены
	tokens, err := s.tokenManager.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	// Сохраняем refresh-токен в Redis
	err = s.redisRepo.SaveToken(ctx, user.ID, tokens.RefreshToken, 24*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

// RefreshTokens обновляет access и refresh токены.
func (s *UserService) RefreshTokens(ctx context.Context, refreshToken string) (*token.TokenPair, *domain.User, error) {
	// Проверяем refresh-токен
	userID, err := s.tokenManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil, errors.New("невалидный refresh-токен")
	}

	// Получаем пользователя
	user, err := s.postgresRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, nil, errors.New("пользователь не найден")
	}

	// Проверяем, совпадает ли токен с тем, что в Redis
	storedToken, err := s.redisRepo.GetToken(ctx, userID)
	if err != nil || storedToken != refreshToken {
		return nil, nil, errors.New("refresh-токен не совпадает")
	}

	// Генерируем новый токен
	tokens, err := s.tokenManager.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, nil, err
	}

	// Обновляем refresh-токен в Redis
	err = s.redisRepo.SaveToken(ctx, user.ID, tokens.RefreshToken, 24*time.Hour)
	if err != nil {
		return nil, nil, err
	}

	return tokens, user, nil
}

// GetMe возвращает данные текущего пользователя.
func (s *UserService) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	// Пробуем получить пользователя из кеша Redis
	user, err := s.redisRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// Если нет в Redis, ищем в БД
	user, err = s.postgresRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Кешируем пользователя
	err = s.redisRepo.SaveUser(ctx, user, time.Hour)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LogOut удаляет refresh-токен из Redis.
func (s *UserService) LogOut(ctx context.Context, userID string) error {
	return s.redisRepo.DeleteToken(ctx, userID)
}
