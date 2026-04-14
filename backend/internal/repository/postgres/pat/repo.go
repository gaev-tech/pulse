package pat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/gaevivan/pulse/internal/domain/user"
)

var _ domain.PATRepository = (*Repo)(nil)

// Repo is the PostgreSQL implementation of domain.PATRepository.
type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new PAT Repo.
func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (repo *Repo) GetUserIDByHash(ctx context.Context, tokenHash string) (string, error) {
	var userID string
	err := repo.pool.QueryRow(ctx,
		`SELECT user_id FROM private_access_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("query: %w", err)
	}
	return userID, nil
}
