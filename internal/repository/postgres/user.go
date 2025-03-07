package postgres

import (
	"context"
	"database/sql"
	"errors"

	"test123/internal/domain"

	"github.com/google/uuid"
)

// UserRepository реализует взаимодействие с PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создает новый экземпляр UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create добавляет нового пользователя в базу данных.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, email, password, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetByEmail ищет пользователя по email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, err
	}
	return &user, nil
}

// Delete удаляет пользователя по ID.
func (r *UserRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
