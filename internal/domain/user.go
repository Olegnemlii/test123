package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the database
type User struct {
	ID          uuid.UUID
	Email       string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   sql.NullTime
	IsConfirmed bool
}

// Token represents a token in the database
type Token struct {
	ID           int
	AccessToken  string
	RefreshToken string
	UserID       uuid.UUID
}

// CodeSignature represents a code and signature in the database
type CodeSignature struct {
	Code      uuid.UUID
	Signature uuid.UUID
	UserID    uuid.UUID
	IsUsed    bool
	ExpiresAt time.Time
}
