# Infrastructure Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Поднять полный dev/prod стек (Go backend, SvelteKit frontend, PostgreSQL, MinIO, Nginx) со всеми миграциями БД, make-командами и smoke-тестом.

**Architecture:** Монорепо с разделением на `backend/`, `frontend/`, `e2e/`, `nginx/`. Backend — Go с chi, postgres через pgx/v5, миграции через golang-migrate при старте. Frontend — SvelteKit. Всё поднимается через docker-compose.

**Tech Stack:** Go 1.23, SvelteKit + TypeScript, PostgreSQL 16, MinIO, Nginx, Docker Compose, golang-migrate, chi, pgx/v5, zap, air (hot reload), Playwright

---

## File Map

**Root:**
- Create: `docker-compose.yml`
- Create: `docker-compose.dev.yml`
- Create: `Makefile`
- Create: `.env.example`
- Create: `nginx/conf.d/default.conf`

**Backend:**
- Create: `backend/go.mod`
- Create: `backend/Dockerfile`
- Create: `backend/Dockerfile.dev`
- Create: `backend/.air.toml`
- Create: `backend/cmd/api/main.go`
- Create: `backend/cmd/seed/main.go`
- Create: `backend/internal/infrastructure/config/config.go`
- Create: `backend/internal/infrastructure/logger/logger.go`
- Create: `backend/internal/infrastructure/postgres/pool.go`
- Create: `backend/internal/repository/migrations/000001_users.up.sql`
- Create: `backend/internal/repository/migrations/000001_users.down.sql`
- Create: `backend/internal/repository/migrations/000002_auth.up.sql`
- Create: `backend/internal/repository/migrations/000002_auth.down.sql`
- Create: `backend/internal/repository/migrations/000003_teams.up.sql`
- Create: `backend/internal/repository/migrations/000003_teams.down.sql`
- Create: `backend/internal/repository/migrations/000004_labels.up.sql`
- Create: `backend/internal/repository/migrations/000004_labels.down.sql`
- Create: `backend/internal/repository/migrations/000005_tasks.up.sql`
- Create: `backend/internal/repository/migrations/000005_tasks.down.sql`
- Create: `backend/internal/repository/migrations/000006_permissions.up.sql`
- Create: `backend/internal/repository/migrations/000006_permissions.down.sql`
- Create: `backend/internal/repository/migrations/000007_filters.up.sql`
- Create: `backend/internal/repository/migrations/000007_filters.down.sql`
- Create: `backend/internal/repository/migrations/000008_events.up.sql`
- Create: `backend/internal/repository/migrations/000008_events.down.sql`
- Create: `backend/internal/repository/migrations/000009_automations.up.sql`
- Create: `backend/internal/repository/migrations/000009_automations.down.sql`
- Create: `backend/internal/repository/migrations/000010_settings.up.sql`
- Create: `backend/internal/repository/migrations/000010_settings.down.sql`
- Create: `backend/internal/repository/migrations/000011_subscriptions.up.sql`
- Create: `backend/internal/repository/migrations/000011_subscriptions.down.sql`
- Create: `backend/internal/repository/migrations/000012_imports.up.sql`
- Create: `backend/internal/repository/migrations/000012_imports.down.sql`
- Create: `backend/pkg/apierror/apierror.go`
- Create: `backend/tests/health_test.go`

**Frontend:**
- Create: `frontend/package.json`
- Create: `frontend/svelte.config.js`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/src/app.html`
- Create: `frontend/src/routes/+layout.svelte`
- Create: `frontend/src/routes/+page.svelte`
- Create: `frontend/Dockerfile`
- Create: `frontend/Dockerfile.dev`
- Create: `frontend/.eslintrc.cjs`

**E2E:**
- Create: `e2e/package.json`
- Create: `e2e/playwright.config.ts`
- Create: `e2e/tests/smoke.spec.ts`

---

## Task 1: Root project files

**Files:**
- Create: `docker-compose.yml`
- Create: `docker-compose.dev.yml`
- Create: `.env.example`
- Create: `nginx/conf.d/default.conf`

- [ ] **Step 1: Create `.env.example`**

```
POSTGRES_DB=pulse
POSTGRES_USER=pulse
POSTGRES_PASSWORD=pulse_secret
POSTGRES_HOST=postgres
POSTGRES_PORT=5432

MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_ENDPOINT=minio:9000
MINIO_BUCKET=pulse-attachments

JWT_SECRET=change_me_in_production
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h

RESEND_API_KEY=re_xxxx
APP_URL=http://localhost
PORT=8080

ENV=development
```

- [ ] **Step 2: Create `nginx/conf.d/default.conf`**

```nginx
upstream backend {
    server backend:8080;
}

upstream frontend_prod {
    server frontend:3000;
}

upstream frontend_dev {
    server frontend:5173;
}

server {
    listen 80;

    location /api/ {
        proxy_pass http://backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /api/v1/ws/ {
        proxy_pass http://backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    location / {
        proxy_pass http://frontend_dev;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }
}
```

- [ ] **Step 3: Create `docker-compose.yml` (production)**

```yaml
services:
  backend:
    build: ./backend
    env_file: .env
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_started
    restart: unless-stopped

  frontend:
    build: ./frontend
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    env_file: .env
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $${POSTGRES_USER} -d $${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 10
    restart: unless-stopped

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    env_file: .env
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    volumes:
      - minio_data:/data
    restart: unless-stopped

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
    depends_on:
      - backend
      - frontend
    restart: unless-stopped

volumes:
  postgres_data:
  minio_data:
```

- [ ] **Step 4: Create `docker-compose.dev.yml` (development)**

```yaml
services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    volumes:
      - ./backend:/app
      - go_cache:/root/go/pkg/mod
    env_file: .env
    depends_on:
      postgres:
        condition: service_healthy
      minio:
        condition: service_started
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    volumes:
      - ./frontend:/app
      - /app/node_modules
    env_file: .env
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    env_file: .env
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $${POSTGRES_USER} -d $${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 10

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    env_file: .env
    environment:
      MINIO_ROOT_USER: ${MINIO_ROOT_USER}
      MINIO_ROOT_PASSWORD: ${MINIO_ROOT_PASSWORD}
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx/conf.d:/etc/nginx/conf.d:ro
    depends_on:
      - backend
      - frontend

volumes:
  postgres_data:
  minio_data:
  go_cache:
```

- [ ] **Step 5: Commit**

```bash
git add .env.example nginx/ docker-compose.yml docker-compose.dev.yml
git commit -m "infra: add docker-compose, nginx config, env example"
```

---

## Task 2: Makefile

**Files:**
- Create: `Makefile`

- [ ] **Step 1: Create `Makefile`**

```makefile
.PHONY: dev prod lint test-api test-e2e swagger seed

dev:
	cp -n .env.example .env 2>/dev/null || true
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose up --build

lint:
	docker compose -f docker-compose.dev.yml run --rm --no-deps backend \
		sh -c "cd /app && golangci-lint run ./..."
	docker compose -f docker-compose.dev.yml run --rm --no-deps frontend \
		sh -c "npm run lint"

test-api:
	docker compose -f docker-compose.dev.yml run --rm backend \
		go test ./tests/... -v -count=1

test-e2e:
	docker compose -f docker-compose.dev.yml up -d --wait
	cd e2e && npx playwright test

swagger:
	docker compose -f docker-compose.dev.yml run --rm --no-deps backend \
		sh -c "swag init -g cmd/api/main.go -o api"

seed:
	docker compose -f docker-compose.dev.yml run --rm backend \
		go run cmd/seed/main.go
```

- [ ] **Step 2: Commit**

```bash
git add Makefile
git commit -m "infra: add Makefile"
```

---

## Task 3: Backend skeleton — go.mod, Dockerfiles, air

**Files:**
- Create: `backend/go.mod`
- Create: `backend/Dockerfile`
- Create: `backend/Dockerfile.dev`
- Create: `backend/.air.toml`

- [ ] **Step 1: Create `backend/go.mod`**

```
module github.com/gaev-tech/pulse

go 1.23

require (
    github.com/go-chi/chi/v5 v5.1.0
    github.com/jackc/pgx/v5 v5.7.1
    github.com/golang-migrate/migrate/v4 v4.17.1
    github.com/golang-jwt/jwt/v5 v5.2.1
    go.uber.org/zap v1.27.0
    golang.org/x/crypto v0.28.0
    github.com/swaggo/swag v1.16.3
    github.com/swaggo/http-swagger/v2 v2.0.2
    github.com/minio/minio-go/v7 v7.0.78
    github.com/google/uuid v1.6.0
    github.com/resend/resend-go/v2 v2.12.0
)
```

- [ ] **Step 2: Create `backend/Dockerfile` (production)**

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /bin/api /api
COPY --from=builder /app/internal/repository/migrations /migrations
EXPOSE 8080
CMD ["/api"]
```

- [ ] **Step 3: Create `backend/Dockerfile.dev`**

```dockerfile
FROM golang:1.23-alpine
WORKDIR /app
RUN apk add --no-cache git curl
RUN go install github.com/air-verse/air@latest && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum ./
RUN go mod download
EXPOSE 8080
CMD ["air", "-c", ".air.toml"]
```

- [ ] **Step 4: Create `backend/.air.toml`**

```toml
[build]
  bin = "./tmp/api"
  cmd = "go build -o ./tmp/api ./cmd/api"
  delay = 1000
  exclude_dir = ["tmp", "tests", "api"]
  include_ext = ["go"]
  kill_delay = "0s"
  rerun = false

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = true
```

- [ ] **Step 5: Commit**

```bash
git add backend/go.mod backend/Dockerfile backend/Dockerfile.dev backend/.air.toml
git commit -m "infra: add backend go.mod, Dockerfiles, air config"
```

---

## Task 4: Backend infrastructure packages

**Files:**
- Create: `backend/internal/infrastructure/config/config.go`
- Create: `backend/internal/infrastructure/logger/logger.go`
- Create: `backend/internal/infrastructure/postgres/pool.go`
- Create: `backend/pkg/apierror/apierror.go`

- [ ] **Step 1: Create `backend/internal/infrastructure/config/config.go`**

```go
package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port            string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	MinioEndpoint   string
	MinioUser       string
	MinioPassword   string
	MinioBucket     string
	ResendAPIKey    string
	AppURL          string
	Env             string
}

func Load() *Config {
	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "720h"))

	return &Config{
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     buildDSN(),
		JWTSecret:       mustEnv("JWT_SECRET"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		MinioEndpoint:   getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioUser:       mustEnv("MINIO_ROOT_USER"),
		MinioPassword:   mustEnv("MINIO_ROOT_PASSWORD"),
		MinioBucket:     getEnv("MINIO_BUCKET", "pulse-attachments"),
		ResendAPIKey:    getEnv("RESEND_API_KEY", ""),
		AppURL:          getEnv("APP_URL", "http://localhost"),
		Env:             getEnv("ENV", "development"),
	}
}

func buildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		mustEnv("POSTGRES_USER"),
		mustEnv("POSTGRES_PASSWORD"),
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "pulse"),
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		panic("missing required env variable: " + key)
	}
	return v
}
```

- [ ] **Step 2: Create `backend/internal/infrastructure/logger/logger.go`**

```go
package logger

import "go.uber.org/zap"

func New(prod bool) *zap.Logger {
	var l *zap.Logger
	var err error
	if prod {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic("failed to init logger: " + err.Error())
	}
	return l
}
```

- [ ] **Step 3: Create `backend/internal/infrastructure/postgres/pool.go`**

```go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pool.Ping: %w", err)
	}
	return pool, nil
}
```

- [ ] **Step 4: Create `backend/pkg/apierror/apierror.go`**

```go
package apierror

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Response struct {
	Error Error `json:"error"`
}

func Write(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response{Error: Error{Code: code, Message: message}})
}

func Unauthorized(w http.ResponseWriter) {
	Write(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required")
}

func Forbidden(w http.ResponseWriter) {
	Write(w, http.StatusForbidden, "PERMISSION_DENIED", "permission denied")
}

func NotFound(w http.ResponseWriter) {
	Write(w, http.StatusNotFound, "NOT_FOUND", "resource not found")
}

func BadRequest(w http.ResponseWriter, message string) {
	Write(w, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func QuotaExceeded(w http.ResponseWriter) {
	Write(w, http.StatusForbidden, "QUOTA_EXCEEDED", "subscription quota exceeded")
}

func Conflict(w http.ResponseWriter, message string) {
	Write(w, http.StatusConflict, "CONFLICT", message)
}

func InternalError(w http.ResponseWriter) {
	Write(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal server error")
}
```

- [ ] **Step 5: Commit**

```bash
git add backend/internal/infrastructure/ backend/pkg/
git commit -m "infra: add config, logger, postgres pool, apierror packages"
```

---

## Task 5: Write smoke test first, then main.go

**Files:**
- Create: `backend/tests/health_test.go`
- Create: `backend/cmd/api/main.go`

- [ ] **Step 1: Write failing smoke test**

Create `backend/tests/health_test.go`:

```go
package tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
)

func baseURL() string {
	if u := os.Getenv("API_URL"); u != "" {
		return u
	}
	return "http://localhost:8080"
}

func TestHealthEndpoint(t *testing.T) {
	resp, err := http.Get(baseURL() + "/api/v1/health")
	if err != nil {
		t.Fatalf("GET /api/v1/health failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status=ok, got %q", body["status"])
	}
}
```

- [ ] **Step 2: Create `backend/cmd/api/main.go`**

```go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"

	"github.com/gaev-tech/pulse/internal/infrastructure/config"
	"github.com/gaev-tech/pulse/internal/infrastructure/logger"
	"github.com/gaev-tech/pulse/internal/infrastructure/postgres"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Env == "production")
	defer log.Sync()

	ctx := context.Background()

	pool, err := postgres.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()
	log.Info("connected to database")

	if err := runMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal("migrations failed", zap.Error(err))
	}
	log.Info("migrations applied")

	r := chi.NewRouter()

	r.Get("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Info("starting server", zap.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("server shutdown error", zap.Error(err))
	}
}

func runMigrations(databaseURL string) error {
	m, err := migrate.New("file:///migrations", databaseURL)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer m.Close()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("m.Up: %w", err)
	}
	return nil
}
```

- [ ] **Step 3: Run test (expect fail — server not running yet)**

```bash
cd backend && go test ./tests/... -run TestHealthEndpoint -v
```

Expected: `FAIL` — `connection refused`

- [ ] **Step 4: Commit**

```bash
git add backend/cmd/api/main.go backend/tests/
git commit -m "infra: add main.go with health endpoint and smoke test"
```

---

## Task 6: Database migrations

**Files:** `backend/internal/repository/migrations/000001_users.up.sql` … `000012_imports.down.sql`

- [ ] **Step 1: Create `000001_users.up.sql`**

```sql
CREATE TABLE users (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL UNIQUE,
    username   TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- [ ] **Step 2: Create `000001_users.down.sql`**

```sql
DROP TABLE IF EXISTS users;
```

- [ ] **Step 3: Create `000002_auth.up.sql`**

```sql
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON refresh_tokens (user_id);

CREATE TABLE private_access_tokens (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON private_access_tokens (user_id);

CREATE TABLE magic_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT NOT NULL,
    token_hash TEXT NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON magic_links (email);
```

- [ ] **Step 4: Create `000002_auth.down.sql`**

```sql
DROP TABLE IF EXISTS magic_links;
DROP TABLE IF EXISTS private_access_tokens;
DROP TABLE IF EXISTS refresh_tokens;
```

- [ ] **Step 5: Create `000003_teams.up.sql`**

```sql
CREATE TABLE teams (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL UNIQUE,
    prefix     TEXT NOT NULL UNIQUE,
    owner_id   UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE team_members (
    team_id   UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (team_id, user_id)
);
CREATE INDEX ON team_members (user_id);
```

- [ ] **Step 6: Create `000003_teams.down.sql`**

```sql
DROP TABLE IF EXISTS team_members;
DROP TABLE IF EXISTS teams;
```

- [ ] **Step 7: Create `000004_labels.up.sql`**

```sql
CREATE TABLE labels (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type TEXT NOT NULL,
    owner_id   UUID NOT NULL,
    name       TEXT NOT NULL,
    color      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (owner_type, owner_id, name)
);
CREATE INDEX ON labels (owner_type, owner_id);
```

- [ ] **Step 8: Create `000004_labels.down.sql`**

```sql
DROP TABLE IF EXISTS labels;
```

- [ ] **Step 9: Create `000005_tasks.up.sql`**

```sql
CREATE TABLE task_sequences (
    owner_type  TEXT   NOT NULL,
    owner_id    UUID   NOT NULL,
    last_number BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (owner_type, owner_id)
);

CREATE TABLE tasks (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_number  BIGINT NOT NULL,
    owner_type  TEXT   NOT NULL,
    owner_id    UUID   NOT NULL,
    title       TEXT   NOT NULL,
    description TEXT,
    status      TEXT   NOT NULL DEFAULT 'opened',
    assignee_id UUID   REFERENCES users(id),
    parent_id   UUID   REFERENCES tasks(id),
    created_by  UUID   NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    origin_id   TEXT UNIQUE,
    UNIQUE (owner_type, owner_id, key_number)
);
CREATE INDEX ON tasks (owner_type, owner_id);
CREATE INDEX ON tasks (assignee_id);
CREATE INDEX ON tasks (parent_id);

CREATE TABLE task_labels (
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, label_id)
);

CREATE TABLE task_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    title      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON task_links (task_id);

CREATE TABLE task_relations (
    task_id         UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    related_task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, related_task_id),
    CHECK (task_id < related_task_id)
);
CREATE INDEX ON task_relations (related_task_id);

CREATE TABLE task_blocking (
    blocker_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);
CREATE INDEX ON task_blocking (blocked_id);

CREATE TABLE task_attachments (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    name       TEXT   NOT NULL,
    url        TEXT   NOT NULL,
    size       BIGINT NOT NULL,
    created_by UUID   NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON task_attachments (task_id);

CREATE TABLE task_opens (
    task_id   UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, user_id)
);
CREATE INDEX ON task_opens (user_id, opened_at DESC);
```

- [ ] **Step 10: Create `000005_tasks.down.sql`**

```sql
DROP TABLE IF EXISTS task_opens;
DROP TABLE IF EXISTS task_attachments;
DROP TABLE IF EXISTS task_blocking;
DROP TABLE IF EXISTS task_relations;
DROP TABLE IF EXISTS task_links;
DROP TABLE IF EXISTS task_labels;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS task_sequences;
```

- [ ] **Step 11: Create `000006_permissions.up.sql`**

```sql
CREATE TABLE permissions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_type  TEXT NOT NULL,
    subject_id    UUID NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   UUID,
    action        TEXT NOT NULL,
    granted_by    UUID REFERENCES users(id),
    level         TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON permissions (subject_type, subject_id, resource_type, resource_id, action);
```

- [ ] **Step 12: Create `000006_permissions.down.sql`**

```sql
DROP TABLE IF EXISTS permissions;
```

- [ ] **Step 13: Create `000007_filters.up.sql`**

```sql
CREATE TABLE filters (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type      TEXT NOT NULL,
    owner_id        UUID NOT NULL,
    name            TEXT NOT NULL,
    filter_mode     TEXT NOT NULL DEFAULT 'simple',
    search_contains TEXT,
    assignee_ids    UUID[],
    status          TEXT,
    label_ids       UUID[],
    rsql            TEXT,
    team_id         UUID REFERENCES teams(id),
    created_by      UUID NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON filters (owner_type, owner_id);

CREATE TABLE filter_settings (
    filter_id    UUID PRIMARY KEY REFERENCES filters(id) ON DELETE CASCADE,
    columns      TEXT[] NOT NULL DEFAULT '{}',
    sort1_column TEXT,
    sort1_dir    TEXT,
    sort2_column TEXT,
    sort2_dir    TEXT,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- [ ] **Step 14: Create `000007_filters.down.sql`**

```sql
DROP TABLE IF EXISTS filter_settings;
DROP TABLE IF EXISTS filters;
```

- [ ] **Step 15: Create `000008_events.up.sql`**

```sql
CREATE TABLE events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type    TEXT NOT NULL,
    actor_id      UUID NOT NULL REFERENCES users(id),
    resource_type TEXT NOT NULL,
    resource_id   UUID NOT NULL,
    team_id       UUID REFERENCES teams(id),
    payload       JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON events (actor_id, created_at DESC);
CREATE INDEX ON events (team_id, created_at DESC) WHERE team_id IS NOT NULL;
CREATE INDEX ON events (resource_id, created_at DESC);
CREATE INDEX ON events USING GIN (payload);
```

- [ ] **Step 16: Create `000008_events.down.sql`**

```sql
DROP TABLE IF EXISTS events;
```

- [ ] **Step 17: Create `000009_automations.up.sql`**

```sql
CREATE TABLE automations (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type TEXT    NOT NULL,
    owner_id   UUID    NOT NULL,
    name       TEXT    NOT NULL,
    trigger    TEXT    NOT NULL,
    criteria   JSONB   NOT NULL DEFAULT '{}',
    actions    JSONB   NOT NULL DEFAULT '[]',
    enabled    BOOLEAN NOT NULL DEFAULT true,
    created_by UUID    NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON automations (owner_type, owner_id);
```

- [ ] **Step 18: Create `000009_automations.down.sql`**

```sql
DROP TABLE IF EXISTS automations;
```

- [ ] **Step 19: Create `000010_settings.up.sql`**

```sql
CREATE TABLE user_settings (
    user_id               UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    language              TEXT    NOT NULL DEFAULT 'en',
    theme                 TEXT    NOT NULL DEFAULT 'system',
    sidebar_personal_open BOOLEAN NOT NULL DEFAULT true,
    sidebar_teams_open    BOOLEAN NOT NULL DEFAULT true,
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_sidebar_team_states (
    user_id UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id UUID    NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    open    BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, team_id)
);
```

- [ ] **Step 20: Create `000010_settings.down.sql`**

```sql
DROP TABLE IF EXISTS user_sidebar_team_states;
DROP TABLE IF EXISTS user_settings;
```

- [ ] **Step 21: Create `000011_subscriptions.up.sql`**

```sql
CREATE TABLE subscriptions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_type TEXT NOT NULL,
    subject_id   UUID NOT NULL,
    plan         TEXT NOT NULL DEFAULT 'free',
    status       TEXT NOT NULL DEFAULT 'active',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (subject_type, subject_id)
);
CREATE INDEX ON subscriptions (subject_type, subject_id);
```

- [ ] **Step 22: Create `000011_subscriptions.down.sql`**

```sql
DROP TABLE IF EXISTS subscriptions;
```

- [ ] **Step 23: Create `000012_imports.up.sql`**

```sql
CREATE TABLE imports (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id            UUID NOT NULL REFERENCES users(id),
    team_id            UUID REFERENCES teams(id),
    source             TEXT NOT NULL,
    status             TEXT NOT NULL DEFAULT 'in_progress',
    progress_total     INT  NOT NULL DEFAULT 0,
    progress_processed INT  NOT NULL DEFAULT 0,
    result             JSONB,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON imports (user_id);
```

- [ ] **Step 24: Create `000012_imports.down.sql`**

```sql
DROP TABLE IF EXISTS imports;
```

- [ ] **Step 25: Commit**

```bash
git add backend/internal/repository/migrations/
git commit -m "infra: add all database migrations (000001-000012)"
```

---

## Task 7: Run `make dev` — verify stack starts and migrations apply

- [ ] **Step 1: Copy .env.example to .env**

```bash
cp .env.example .env
```

Edit `.env` — all defaults work for local dev as-is.

- [ ] **Step 2: Start the stack**

```bash
make dev
```

Expected in logs:
- `postgres` — `database system is ready to accept connections`
- `backend` — `connected to database`
- `backend` — `migrations applied`
- `backend` — `starting server addr=:8080`
- `nginx` — starts without error

- [ ] **Step 3: Verify health endpoint through nginx**

```bash
curl http://localhost/api/v1/health
```

Expected: `{"status":"ok"}`

- [ ] **Step 4: Commit**

```bash
git add .env.example  # .env is gitignored
git commit -m "infra: verify stack starts successfully"
```

---

## Task 8: Run smoke test via `make test-api`

- [ ] **Step 1: Run test against running stack**

```bash
make test-api
```

Expected output:
```
--- PASS: TestHealthEndpoint (0.01s)
PASS
ok  github.com/gaev-tech/pulse/tests  0.XXXs
```

If `FAIL`: check that `make dev` is running and `API_URL` env is reachable from the test container.

- [ ] **Step 2: Commit** (only if something was adjusted)

```bash
git add .
git commit -m "infra: smoke test passes"
```

---

## Task 9: Seed script

**Files:**
- Create: `backend/cmd/seed/main.go`

- [ ] **Step 1: Create `backend/cmd/seed/main.go`**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gaev-tech/pulse/internal/infrastructure/config"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("ping: %v", err)
	}

	// Seed user
	var userID string
	err = pool.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ('admin@example.com', 'admin')
		ON CONFLICT (email) DO UPDATE SET username = EXCLUDED.username
		RETURNING id
	`).Scan(&userID)
	if err != nil {
		log.Fatalf("seed user: %v", err)
	}
	fmt.Printf("user: %s\n", userID)

	// Seed user settings
	_, err = pool.Exec(ctx, `
		INSERT INTO user_settings (user_id) VALUES ($1)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	if err != nil {
		log.Fatalf("seed user_settings: %v", err)
	}

	// Seed subscription for user
	_, err = pool.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('user', $1, 'free')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, userID)
	if err != nil {
		log.Fatalf("seed subscription: %v", err)
	}

	// Seed team
	var teamID string
	err = pool.QueryRow(ctx, `
		INSERT INTO teams (name, prefix, owner_id)
		VALUES ('Backend Team', 'BACK', $1)
		ON CONFLICT (prefix) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, userID).Scan(&teamID)
	if err != nil {
		log.Fatalf("seed team: %v", err)
	}
	fmt.Printf("team: %s\n", teamID)

	// Seed team member
	_, err = pool.Exec(ctx, `
		INSERT INTO team_members (team_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT DO NOTHING
	`, teamID, userID)
	if err != nil {
		log.Fatalf("seed team_member: %v", err)
	}

	// Seed subscription for team
	_, err = pool.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('team', $1, 'free')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, teamID)
	if err != nil {
		log.Fatalf("seed team subscription: %v", err)
	}

	// Seed personal task sequence
	_, err = pool.Exec(ctx, `
		INSERT INTO task_sequences (owner_type, owner_id, last_number)
		VALUES ('user', $1, 0)
		ON CONFLICT DO NOTHING
	`, userID)
	if err != nil {
		log.Fatalf("seed task_sequence: %v", err)
	}

	// Seed personal task
	var taskID string
	err = pool.QueryRow(ctx, `
		INSERT INTO tasks (key_number, owner_type, owner_id, title, description, status, created_by)
		VALUES (1, 'user', $1, 'First task', 'This is a seed task.', 'opened', $1)
		ON CONFLICT (owner_type, owner_id, key_number) DO UPDATE SET title = EXCLUDED.title
		RETURNING id
	`, userID).Scan(&taskID)
	if err != nil {
		log.Fatalf("seed task: %v", err)
	}
	fmt.Printf("task: %s\n", taskID)

	_ = time.Now() // suppress import warning
	fmt.Println("seed complete")
}
```

- [ ] **Step 2: Run seed**

```bash
make seed
```

Expected:
```
user: <uuid>
team: <uuid>
task: <uuid>
seed complete
```

- [ ] **Step 3: Commit**

```bash
git add backend/cmd/seed/
git commit -m "infra: add seed script with user, team, task"
```

---

## Task 10: Frontend SvelteKit skeleton

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/svelte.config.js`
- Create: `frontend/vite.config.ts`
- Create: `frontend/tsconfig.json`
- Create: `frontend/src/app.html`
- Create: `frontend/src/routes/+layout.svelte`
- Create: `frontend/src/routes/+page.svelte`
- Create: `frontend/Dockerfile`
- Create: `frontend/Dockerfile.dev`
- Create: `frontend/.eslintrc.cjs`

- [ ] **Step 1: Create `frontend/package.json`**

```json
{
  "name": "pulse-frontend",
  "version": "0.0.1",
  "private": true,
  "scripts": {
    "dev": "vite dev",
    "build": "vite build",
    "preview": "vite preview",
    "check": "svelte-kit sync && svelte-check --tsconfig ./tsconfig.json",
    "lint": "eslint src"
  },
  "devDependencies": {
    "@sveltejs/adapter-node": "^5.2.0",
    "@sveltejs/kit": "^2.7.0",
    "@sveltejs/vite-plugin-svelte": "^4.0.0",
    "@typescript-eslint/eslint-plugin": "^8.0.0",
    "@typescript-eslint/parser": "^8.0.0",
    "eslint": "^9.0.0",
    "eslint-plugin-svelte": "^2.44.0",
    "svelte": "^5.0.0",
    "svelte-check": "^4.0.0",
    "typescript": "^5.6.0",
    "vite": "^5.4.0"
  }
}
```

- [ ] **Step 2: Create `frontend/svelte.config.js`**

```js
import adapter from '@sveltejs/adapter-node';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter(),
    alias: {
      '@pulse': './src'
    }
  }
};

export default config;
```

- [ ] **Step 3: Create `frontend/vite.config.ts`**

```ts
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    host: '0.0.0.0',
    port: 5173
  }
});
```

- [ ] **Step 4: Create `frontend/tsconfig.json`**

```json
{
  "extends": "./.svelte-kit/tsconfig.json",
  "compilerOptions": {
    "allowJs": true,
    "checkJs": true,
    "esModuleInterop": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "skipLibCheck": true,
    "sourceMap": true,
    "strict": true
  }
}
```

- [ ] **Step 5: Create `frontend/src/app.html`**

```html
<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <link rel="icon" href="%sveltekit.assets%/favicon.png" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    %sveltekit.head%
  </head>
  <body data-sveltekit-preload-data="hover">
    <div style="display: contents">%sveltekit.body%</div>
  </body>
</html>
```

- [ ] **Step 6: Create `frontend/src/routes/+layout.svelte`**

```svelte
<script lang="ts">
  let { children } = $props();
</script>

{@render children()}
```

- [ ] **Step 7: Create `frontend/src/routes/+page.svelte`**

```svelte
<h1>Pulse</h1>
<p>Task tracker.</p>
```

- [ ] **Step 8: Create `frontend/.eslintrc.cjs`**

```js
module.exports = {
  root: true,
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:svelte/recommended'
  ],
  parser: '@typescript-eslint/parser',
  plugins: ['@typescript-eslint'],
  parserOptions: {
    sourceType: 'module',
    ecmaVersion: 2020,
    extraFileExtensions: ['.svelte']
  },
  overrides: [
    {
      files: ['*.svelte'],
      parser: 'svelte-eslint-parser',
      parserOptions: {
        parser: '@typescript-eslint/parser'
      }
    }
  ],
  env: {
    browser: true,
    es2017: true,
    node: true
  }
};
```

- [ ] **Step 9: Create `frontend/Dockerfile` (production)**

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm ci --omit=dev
EXPOSE 3000
ENV NODE_ENV=production
CMD ["node", "build"]
```

- [ ] **Step 10: Create `frontend/Dockerfile.dev`**

```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci
EXPOSE 5173
CMD ["npm", "run", "dev"]
```

- [ ] **Step 11: Commit**

```bash
git add frontend/
git commit -m "infra: add SvelteKit frontend skeleton"
```

---

## Task 11: E2E Playwright skeleton

**Files:**
- Create: `e2e/package.json`
- Create: `e2e/playwright.config.ts`
- Create: `e2e/tests/smoke.spec.ts`

- [ ] **Step 1: Create `e2e/package.json`**

```json
{
  "name": "pulse-e2e",
  "version": "0.0.1",
  "private": true,
  "scripts": {
    "test": "playwright test",
    "test:ui": "playwright test --ui"
  },
  "devDependencies": {
    "@playwright/test": "^1.48.0"
  }
}
```

- [ ] **Step 2: Create `e2e/playwright.config.ts`**

```ts
import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './tests',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL: process.env.BASE_URL || 'http://localhost',
    trace: 'on-first-retry'
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] }
    },
    {
      name: 'Mobile Chrome',
      use: { ...devices['Pixel 5'] }
    }
  ]
});
```

- [ ] **Step 3: Create `e2e/tests/smoke.spec.ts`**

```ts
import { test, expect } from '@playwright/test';

test('app loads', async ({ page }) => {
  await page.goto('/');
  await expect(page).toHaveTitle(/Pulse/i);
});

test('health endpoint returns ok', async ({ request }) => {
  const response = await request.get('/api/v1/health');
  expect(response.status()).toBe(200);
  const body = await response.json();
  expect(body.status).toBe('ok');
});
```

- [ ] **Step 4: Run e2e smoke test**

```bash
make dev   # if not already running
make test-e2e
```

Expected: 2 tests pass (app loads + health endpoint).

- [ ] **Step 5: Commit**

```bash
git add e2e/
git commit -m "infra: add Playwright e2e skeleton with smoke tests"
```

---

## Verification

После выполнения всех задач:

```bash
# 1. Полный стек поднимается
make dev
# Ожидание: все контейнеры healthy, backend логирует "starting server"

# 2. Health endpoint через nginx
curl http://localhost/api/v1/health
# Ожидание: {"status":"ok"}

# 3. API smoke тест
make test-api
# Ожидание: PASS TestHealthEndpoint

# 4. Seed работает
make seed
# Ожидание: "seed complete", данные в БД

# 5. E2E smoke тест
make test-e2e
# Ожидание: 2 passed (app loads + health)

# 6. Lint (после установки зависимостей)
make lint
# Ожидание: no issues
```
