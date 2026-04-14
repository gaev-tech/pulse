package magiclink

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/gaevivan/pulse/internal/domain/user"
)

var _ domain.MagicLinkRepository = (*Repo)(nil)

// Repo is the PostgreSQL implementation of domain.MagicLinkRepository.
type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new magic link Repo.
func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (repo *Repo) Create(ctx context.Context, email, tokenHash string, expiresAt time.Time) error {
	_, err := repo.pool.Exec(ctx,
		`INSERT INTO magic_link_tokens (email, token_hash, expires_at) VALUES ($1, $2, $3)`,
		email, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func (repo *Repo) GetByHash(ctx context.Context, tokenHash string) (*domain.MagicLinkToken, error) {
	token := &domain.MagicLinkToken{}
	err := repo.pool.QueryRow(ctx,
		`SELECT id, email, token_hash, expires_at, used_at, created_at
		 FROM magic_link_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&token.ID, &token.Email, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return token, nil
}

func (repo *Repo) MarkUsed(ctx context.Context, id string) error {
	_, err := repo.pool.Exec(ctx,
		`UPDATE magic_link_tokens SET used_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
