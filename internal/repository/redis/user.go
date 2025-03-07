package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"test123/internal/domain"

	"github.com/redis/go-redis/v9"
)

// UserRepository хранит пользователей в Redis.
type UserRepository struct {
	client *redis.Client
}

// NewUserRepository создает новый экземпляр UserRepository для Redis.
func NewUserRepository(client *redis.Client) *UserRepository {
	return &UserRepository{client: client}
}

// SaveToken сохраняет refresh-токен в Redis с заданным временем жизни.
func (r *UserRepository) SaveToken(ctx context.Context, userID string, token string, expiration time.Duration) error {
	return r.client.Set(ctx, userID, token, expiration).Err()
}

// GetToken получает refresh-токен пользователя по его ID.
func (r *UserRepository) GetToken(ctx context.Context, userID string) (string, error) {
	token, err := r.client.Get(ctx, userID).Result()
	if err == redis.Nil {
		return "", nil // Токен не найден
	}
	if err != nil {
		return "", err
	}
	return token, nil
}

// DeleteToken удаляет refresh-токен пользователя.
func (r *UserRepository) DeleteToken(ctx context.Context, userID string) error {
	return r.client.Del(ctx, userID).Err()
}

// SaveUser кэширует данные пользователя в Redis.
func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User, expiration time.Duration) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "user:"+user.ID, data, expiration).Err()
}

// GetUser получает пользователя из кэша Redis.
func (r *UserRepository) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	data, err := r.client.Get(ctx, "user:"+userID).Result()
	if err == redis.Nil {
		return nil, nil // Пользователь не найден
	}
	if err != nil {
		return nil, err
	}

	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, errors.New("ошибка декодирования данных пользователя")
	}

	return &user, nil
}

// DeleteUser удаляет пользователя из Redis.
func (r *UserRepository) DeleteUser(ctx context.Context, userID string) error {
	return r.client.Del(ctx, "user:"+userID).Err()
}
