package user

import (
	"context"
	"time"
)

type User struct {
	ID        string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type MagicLinkToken struct {
	ID        string
	Email     string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	Create(ctx context.Context, email, username string) (*User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

type MagicLinkRepository interface {
	Create(ctx context.Context, email, tokenHash string, expiresAt time.Time) error
	GetByHash(ctx context.Context, tokenHash string) (*MagicLinkToken, error)
	MarkUsed(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (*RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (*RefreshToken, error)
	Revoke(ctx context.Context, id string) error
}

type PATRepository interface {
	GetUserIDByHash(ctx context.Context, tokenHash string) (string, error)
}
