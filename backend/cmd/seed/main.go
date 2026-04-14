package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gaevivan/pulse/internal/infrastructure/config"
	"github.com/gaevivan/pulse/internal/infrastructure/postgres"
	"github.com/gaevivan/pulse/internal/repository/migrations"
)

func main() {
	cfg := config.Load()

	if err := postgres.Migrate(cfg.Database, migrations.FS); err != nil {
		fmt.Fprintf(os.Stderr, "migrate: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := seed(ctx, pool); err != nil {
		fmt.Fprintf(os.Stderr, "seed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("seed completed")
}

func seed(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// Users
	var aliceID, bobID string
	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ('alice@example.com', 'alice')
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`).Scan(&aliceID)
	if err != nil {
		return fmt.Errorf("insert alice: %w", err)
	}

	err = tx.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ('bob@example.com', 'bob')
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`).Scan(&bobID)
	if err != nil {
		return fmt.Errorf("insert bob: %w", err)
	}

	// Subscriptions
	_, err = tx.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('user', $1, 'pro'), ('user', $2, 'free')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, aliceID, bobID)
	if err != nil {
		return fmt.Errorf("insert subscriptions: %w", err)
	}

	// Team
	var teamID string
	err = tx.QueryRow(ctx, `
		INSERT INTO teams (name, prefix, owner_id)
		VALUES ('Pulse Dev', 'PLS', $1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, aliceID).Scan(&teamID)
	if err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	// Team subscription
	_, err = tx.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('team', $1, 'team')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, teamID)
	if err != nil {
		return fmt.Errorf("insert team subscription: %w", err)
	}

	// Team members
	_, err = tx.Exec(ctx, `
		INSERT INTO team_members (team_id, user_id)
		VALUES ($1, $2), ($1, $3)
		ON CONFLICT DO NOTHING
	`, teamID, aliceID, bobID)
	if err != nil {
		return fmt.Errorf("insert team members: %w", err)
	}

	// Labels
	_, err = tx.Exec(ctx, `
		INSERT INTO labels (owner_type, owner_id, name, color)
		VALUES
			('team', $1, 'bug', '#d73a4a'),
			('team', $1, 'feature', '#0075ca'),
			('team', $1, 'chore', '#e4e669')
		ON CONFLICT (owner_type, owner_id, name) DO NOTHING
	`, teamID)
	if err != nil {
		return fmt.Errorf("insert labels: %w", err)
	}

	// Task sequence
	_, err = tx.Exec(ctx, `
		INSERT INTO task_sequences (owner_type, owner_id, last_number)
		VALUES ('team', $1, 0)
		ON CONFLICT DO NOTHING
	`, teamID)
	if err != nil {
		return fmt.Errorf("insert task sequence: %w", err)
	}

	// Tasks
	var taskNumber int64
	err = tx.QueryRow(ctx, `
		UPDATE task_sequences SET last_number = last_number + 1
		WHERE owner_type = 'team' AND owner_id = $1
		RETURNING last_number
	`, teamID).Scan(&taskNumber)
	if err != nil {
		return fmt.Errorf("increment task sequence: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO tasks (key_number, owner_type, owner_id, title, status, created_by)
		VALUES ($1, 'team', $2, 'Настроить CI/CD', 'opened', $3)
		ON CONFLICT (owner_type, owner_id, key_number) DO NOTHING
	`, taskNumber, teamID, aliceID)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	// User settings
	_, err = tx.Exec(ctx, `
		INSERT INTO user_settings (user_id, language, theme)
		VALUES ($1, 'ru', 'dark'), ($2, 'en', 'system')
		ON CONFLICT (user_id) DO NOTHING
	`, aliceID, bobID)
	if err != nil {
		return fmt.Errorf("insert user settings: %w", err)
	}

	return tx.Commit(ctx)
}
