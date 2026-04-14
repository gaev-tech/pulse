package tests

import (
	"context"
	"os"
	"testing"

	"github.com/gaevivan/pulse/internal/infrastructure/config"
	"github.com/gaevivan/pulse/internal/infrastructure/postgres"
	"github.com/gaevivan/pulse/internal/repository/migrations"
)

func TestDBSmoke(t *testing.T) {
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("RUN_DB_TESTS=1 not set, skipping DB smoke test")
	}

	cfg := config.Load()

	if err := postgres.Migrate(cfg.Database, migrations.FS); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	ctx := context.Background()
	pool, err := postgres.New(ctx, cfg.Database)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	tables := []string{
		"users", "tasks", "teams", "labels", "filters",
		"subscriptions", "events", "permissions",
	}

	for _, table := range tables {
		var exists bool
		err := pool.QueryRow(ctx,
			`SELECT EXISTS (
				SELECT FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)`, table,
		).Scan(&exists)
		if err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}
		if !exists {
			t.Errorf("table %q not found", table)
		}
	}
}
