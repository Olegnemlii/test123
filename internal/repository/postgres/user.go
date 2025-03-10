package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Olegnemlii/test123/internal/domain"
	"github.com/Olegnemlii/test123/internal/repository"

	"github.com/google/uuid"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) repository.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	// SQL для вставки нового пользователя
	insertUserSQL := `
		INSERT INTO users (id, email, password, created_at, updated_at, is_confirmed)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	id := uuid.New()

	_, err := r.db.ExecContext(ctx, insertUserSQL, id, user.Email, user.Password, user.CreatedAt, user.UpdatedAt, user.IsConfirmed)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	user.ID = id

	return user, nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	// SQL для получения пользователя по ID
	getUserSQL := `
		SELECT id, email, password, created_at, updated_at, deleted_at, is_confirmed
		FROM users
		WHERE id = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, getUserSQL, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.IsConfirmed)
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	// SQL для получения пользователя по email
	getUserSQL := `
		SELECT id, email, password, created_at, updated_at, deleted_at, is_confirmed
		FROM users
		WHERE email = $1
	`
	var user domain.User
	err := r.db.QueryRowContext(ctx, getUserSQL, email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &user.IsConfirmed)
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	// SQL для обновления пользователя
	updateUserSQL := `
		UPDATE users
		SET email = $2, password = $3, updated_at = $4, is_confirmed = $5
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, updateUserSQL, user.ID, user.Email, user.Password, user.UpdatedAt, user.IsConfirmed)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// SQL для удаления пользователя
	deleteUserSQL := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, deleteUserSQL, id)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetEmailBySignature(ctx context.Context, signature uuid.UUID) (string, error) {
	// SQL для получения email по подписи
	getEmailSQL := `
		SELECT u.email
		FROM users u
		JOIN codes_signatures cs ON u.id = cs.user_id
		WHERE cs.signature = $1
	`

	var email string

	err := r.db.QueryRowContext(ctx, getEmailSQL, signature).Scan(&email)

	if err != nil {
		log.Printf("Failed to get email by signature: %v", err)
		return "", fmt.Errorf("failed to get email by signature: %w", err)
	}

	return email, nil
}

func (r *PostgresUserRepository) UpdateUserSignature(ctx context.Context, userID uuid.UUID, signature uuid.UUID) error {
	// SQL для обновления подписи пользователя
	updateSignatureSQL := `
		UPDATE codes_signatures
		SET signature = $2
		WHERE user_id = $1
	`
	_, err := r.db.ExecContext(ctx, updateSignatureSQL, userID, signature)
	if err != nil {
		log.Printf("Failed to update user signature: %v", err)
		return fmt.Errorf("failed to update user signature: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) StoreVerificationCode(ctx context.Context, email string, code string) error {
	userID, err := r.getUserIDByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user ID by email: %w", err)
	}

	codeUUID, err := uuid.Parse(code)
	if err != nil {
		return fmt.Errorf("failed to parse verification code as UUID: %w", err)
	}

	signatureUUID := uuid.New()

	expiresAt := time.Now().UTC().Add(time.Hour * 24)

	// SQL для сохранения кода подтверждения в codes_signatures
	storeCodeSQL := `
		INSERT INTO codes_signatures (code, signature, user_id, is_used, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.ExecContext(ctx, storeCodeSQL, codeUUID, signatureUUID, userID, false, expiresAt)
	if err != nil {
		log.Printf("Failed to store verification code: %v", err)
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetVerificationCode(ctx context.Context, email string) (string, error) {
	userID, err := r.getUserIDByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user ID by email: %w", err)
	}

	// SQL для получения кода подтверждения из codes_signatures
	getCodeSQL := `
		SELECT code
		FROM codes_signatures
		WHERE user_id = $1 AND is_used = false AND expires_at > NOW()
	`

	var code uuid.UUID
	err = r.db.QueryRowContext(ctx, getCodeSQL, userID).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("verification code not found")
		}
		log.Printf("Failed to get verification code: %v", err)
		return "", fmt.Errorf("failed to get verification code: %w", err)
	}

	return code.String(), nil
}

func (r *PostgresUserRepository) DeleteVerificationCode(ctx context.Context, email string) error {
	userID, err := r.getUserIDByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to get user ID by email: %w", err)
	}

	// SQL для удаления кода подтверждения из codes_signatures
	deleteCodeSQL := `
		UPDATE codes_signatures SET is_used = true WHERE user_id = $1
	`
	_, err = r.db.ExecContext(ctx, deleteCodeSQL, userID)
	if err != nil {
		log.Printf("Failed to delete verification code: %v", err)
		return fmt.Errorf("failed to delete verification code: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) StoreRefreshToken(ctx context.Context, email string, refreshToken string) error {
	// SQL для сохранения refresh токена в таблице tokens
	storeTokenSQL := `
		INSERT INTO tokens (access_token, refresh_token, user_id)
		VALUES ($1, $2, (SELECT id FROM users WHERE email = $3))
	`
	accessToken := uuid.New().String()

	_, err := r.db.ExecContext(ctx, storeTokenSQL, accessToken, refreshToken, email)
	if err != nil {
		log.Printf("Failed to store refresh token: %v", err)
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetRefreshToken(ctx context.Context, email string) (string, error) {
	// SQL для получения refresh токена из таблицы tokens
	getTokenSQL := `
		SELECT refresh_token
		FROM tokens
		WHERE user_id = (SELECT id FROM users WHERE email = $1)
	`

	var refreshToken string
	err := r.db.QueryRowContext(ctx, getTokenSQL, email).Scan(&refreshToken)
	if err != nil {
		// Если refresh токен не найден, возвращаем nil без ошибки
		if err == sql.ErrNoRows {
			return "", nil
		}

		log.Printf("Failed to get refresh token: %v", err)
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return refreshToken, nil
}

func (r *PostgresUserRepository) DeleteRefreshToken(ctx context.Context, email string) error {
	// SQL для удаления refresh токена из таблицы tokens
	deleteTokenSQL := `
		DELETE FROM tokens
		WHERE user_id = (SELECT id FROM users WHERE email = $1)
	`
	_, err := r.db.ExecContext(ctx, deleteTokenSQL, email)
	if err != nil {
		log.Printf("Failed to delete refresh token: %v", err)
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) getUserIDByEmail(ctx context.Context, email string) (uuid.UUID, error) {
	// SQL для получения ID пользователя по email
	getUserIDSQL := `
		SELECT id
		FROM users
		WHERE email = $1
	`

	var userID uuid.UUID
	err := r.db.QueryRowContext(ctx, getUserIDSQL, email).Scan(&userID)
	if err != nil {
		log.Printf("Failed to get user ID by email: %v", err)
		return uuid.Nil, fmt.Errorf("failed to get user ID by email: %w", err)
	}

	return userID, nil
}
