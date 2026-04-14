package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/gaevivan/pulse/internal/domain/user"
)

type Repo struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (repo *Repo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	err := repo.pool.QueryRow(ctx,
		`SELECT id, email, username, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return user, nil
}

func (repo *Repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	err := repo.pool.QueryRow(ctx,
		`SELECT id, email, username, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return user, nil
}

func (repo *Repo) Create(ctx context.Context, email, username string) (*domain.User, error) {
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	user := &domain.User{}
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, username) VALUES ($1, $2)
		 RETURNING id, email, username, created_at, updated_at`,
		email, username,
	).Scan(&user.ID, &user.Email, &user.Username, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO subscriptions (subject_type, subject_id, plan) VALUES ('user', $1, 'free')`,
		user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert subscription: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_settings (user_id) VALUES ($1)`,
		user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user_settings: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return user, nil
}

func (repo *Repo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := repo.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query: %w", err)
	}
	return exists, nil
}
