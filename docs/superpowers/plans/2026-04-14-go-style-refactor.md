# Go Style Refactor Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Apply Uber Go Style Guide rules to all 35 `.go` files in `backend/` without changing behaviour.

**Architecture:** Style-only refactoring across all layers (cmd, domain, repository, usecase, handler, infrastructure, pkg). Existing tests act as the safety net — they must pass before and after every task.

**Tech Stack:** Go, chi, pgx, zap, golang-jwt

---

## Violations summary (by rule)

| Rule | Location |
|------|----------|
| Package names with `_` | `repository/postgres/magic_link`, `repository/postgres/refresh_token` |
| Receiver names too long | `auth.go`, `usecase.go`, `jwt.go`, `resend.go`, `log.go`, all repos |
| Unexported package-level var without `_` | `pkg/validator/validator.go` |
| Magic number without named constant | `usecase/user/usecase.go` |
| Variable named `bytes` shadows stdlib | `usecase/user/usecase.go` |
| `fmt.Errorf("%w", errors.Join(...))` — wrong nesting | `usecase/user/usecase.go` |
| `var err` at top + `err =` reuse | `cmd/seed/main.go` |
| Missing interface compliance assertions | all repos |
| Struct field alignment | `infrastructure/email/resend.go` |

---

## Task 1: Rename packages `magic_link` → `magiclink` and `refresh_token` → `refreshtoken`

**Files:**
- Rename dir: `backend/internal/repository/postgres/magic_link/` → `magiclink/`
- Rename dir: `backend/internal/repository/postgres/refresh_token/` → `refreshtoken/`
- Modify: `backend/internal/repository/postgres/magic_link/repo.go`
- Modify: `backend/internal/repository/postgres/refresh_token/repo.go`
- Modify: `backend/cmd/api/main.go`

Uber rule: "Package names: Short, lowercase, singular — no `_`"

- [ ] **Step 1: Rename directories**

```bash
cd /Users/gaevivan/projects/pulse/backend
mv internal/repository/postgres/magic_link internal/repository/postgres/magiclink
mv internal/repository/postgres/refresh_token internal/repository/postgres/refreshtoken
```

- [ ] **Step 2: Update package declaration in magiclink/repo.go**

Change line 1:
```go
package magiclink
```

- [ ] **Step 3: Update package declaration in refreshtoken/repo.go**

Change line 1:
```go
package refreshtoken
```

- [ ] **Step 4: Update imports in cmd/api/main.go**

Replace:
```go
repomagiclink "github.com/gaevivan/pulse/internal/repository/postgres/magic_link"
repopat "github.com/gaevivan/pulse/internal/repository/postgres/pat"
reporefreshtoken "github.com/gaevivan/pulse/internal/repository/postgres/refresh_token"
repouser "github.com/gaevivan/pulse/internal/repository/postgres/user"
```
With:
```go
repomagiclink "github.com/gaevivan/pulse/internal/repository/postgres/magiclink"
repopat "github.com/gaevivan/pulse/internal/repository/postgres/pat"
reporefreshtoken "github.com/gaevivan/pulse/internal/repository/postgres/refreshtoken"
repouser "github.com/gaevivan/pulse/internal/repository/postgres/user"
```

- [ ] **Step 5: Verify it compiles**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./...
```
Expected: no errors.

- [ ] **Step 6: Run tests**

```bash
cd /Users/gaevivan/projects/pulse/backend && go test ./tests/...
```
Expected: all pass.

---

## Task 2: Fix receiver names — handler layer

**Files:**
- Modify: `backend/internal/handler/v1/auth.go`
- Modify: `backend/internal/handler/middleware/auth.go`

Uber rule: "Receivers: Consistent pointer `(s *Service)` — use 1–2 letter abbreviation"

- [ ] **Step 1: Fix AuthHandler receiver in auth.go**

Replace all occurrences of `(handler *AuthHandler)` with `(h *AuthHandler)`:
- Line 30: `func (handler *AuthHandler) SendMagicLink(` → `func (h *AuthHandler) SendMagicLink(`
- Line 70: `func (handler *AuthHandler) VerifyMagicLink(` → `func (h *AuthHandler) VerifyMagicLink(`
- Line 123: `func (handler *AuthHandler) Refresh(` → `func (h *AuthHandler) Refresh(`
- Line 169: `func (handler *AuthHandler) Logout(` → `func (h *AuthHandler) Logout(`
- Line 195: `func (handler *AuthHandler) Me(` → `func (h *AuthHandler) Me(`

Replace all `handler.usecase` → `h.usecase` in method bodies (5 occurrences).

Full updated `auth.go`:
```go
package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gaevivan/pulse/internal/handler/middleware"
	userusecase "github.com/gaevivan/pulse/internal/usecase/user"
	"github.com/gaevivan/pulse/pkg/apierror"
	"github.com/gaevivan/pulse/pkg/validator"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	usecase *userusecase.UseCase
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(usecase *userusecase.UseCase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

// SendMagicLink godoc
// @Summary     Отправить magic-link
// @Tags        auth
// @Accept      json
// @Param       body body object{email=string} true "Email"
// @Success     200
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/magic-link [post]
func (h *AuthHandler) SendMagicLink(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := validator.Decode(r, &body); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}
	if err := validator.Email(body.Email); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}

	if err := h.usecase.SendMagicLink(r.Context(), body.Email); err != nil {
		switch {
		case errors.Is(err, userusecase.ErrEmailUnavailable):
			apierror.EmailUnavailable(w)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(w)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(w)
		default:
			apierror.Internal(w)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

// VerifyMagicLink godoc
// @Summary     Верифицировать magic-link токен
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object{token=string} true "Токен"
// @Success     200 {object} object{access_token=string,refresh_token=string,user=object{id=string,email=string,username=string}}
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/magic-link/verify [post]
func (h *AuthHandler) VerifyMagicLink(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token string `json:"token"`
	}
	if err := validator.Decode(r, &body); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}
	if err := validator.Required(body.Token, "token"); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}

	result, err := h.usecase.VerifyMagicLink(r.Context(), body.Token)
	if err != nil {
		switch {
		case errors.Is(err, userusecase.ErrInvalidToken), errors.Is(err, userusecase.ErrTokenUsed):
			apierror.Unauthorized(w)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(w)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(w)
		case errors.Is(err, userusecase.ErrUsernameConflict):
			apierror.UsernameConflict(w)
		default:
			apierror.Internal(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"user": map[string]string{
			"id":       result.User.ID,
			"email":    result.User.Email,
			"username": result.User.Username,
		},
	})
}

// Refresh godoc
// @Summary     Обновить токены
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       body body object{refresh_token=string} true "Refresh token"
// @Success     200 {object} object{access_token=string,refresh_token=string}
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := validator.Decode(r, &body); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}
	if err := validator.Required(body.RefreshToken, "refresh_token"); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}

	result, err := h.usecase.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, userusecase.ErrInvalidToken):
			apierror.Unauthorized(w)
		case errors.Is(err, userusecase.ErrUserNotFound):
			apierror.UserNotFound(w)
		case errors.Is(err, userusecase.ErrDatabaseUnavailable):
			apierror.DatabaseUnavailable(w)
		case errors.Is(err, userusecase.ErrTokenGenerationFailed):
			apierror.TokenGenerationFailed(w)
		default:
			apierror.Internal(w)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
	})
}

// Logout godoc
// @Summary     Выйти из системы
// @Tags        auth
// @Accept      json
// @Param       body body object{refresh_token=string} true "Refresh token"
// @Success     204
// @Failure     400 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := validator.Decode(r, &body); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}
	if err := validator.Required(body.RefreshToken, "refresh_token"); err != nil {
		apierror.BadRequest(w, err.Error())
		return
	}

	// Silently succeed even for invalid/revoked tokens.
	_ = h.usecase.Logout(r.Context(), body.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

// Me godoc
// @Summary     Текущий пользователь
// @Tags        auth
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} object{id=string,email=string,username=string}
// @Failure     401 {object} object{error=object{code=string,message=string}}
// @Router      /v1/auth/me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		apierror.Unauthorized(w)
		return
	}

	currentUser, err := h.usecase.GetUserByID(r.Context(), userID)
	if err != nil {
		apierror.DatabaseUnavailable(w)
		return
	}
	if currentUser == nil {
		apierror.UserNotFound(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"id":       currentUser.ID,
		"email":    currentUser.Email,
		"username": currentUser.Username,
	})
}
```

- [ ] **Step 2: Fix Auth middleware receiver in middleware/auth.go**

Replace `(auth *Auth)` → `(a *Auth)` in both methods, and `auth.jwt` → `a.jwt`, `auth.pats` → `a.pats`.

Full updated `middleware/auth.go`:
```go
package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gaevivan/pulse/internal/domain/user"
	"github.com/gaevivan/pulse/internal/infrastructure/jwt"
	"github.com/gaevivan/pulse/pkg/apierror"
)

type contextKey string

// UserIDKey is the context key for the authenticated user ID.
const UserIDKey contextKey = "user_id"

// UserIDFromContext returns the authenticated user ID from ctx.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok && userID != ""
}

// Auth is HTTP middleware for JWT and PAT authentication.
type Auth struct {
	jwt  *jwt.Manager
	pats user.PATRepository
}

// NewAuth creates a new Auth middleware.
func NewAuth(jwtManager *jwt.Manager, pats user.PATRepository) *Auth {
	return &Auth{jwt: jwtManager, pats: pats}
}

// Required enforces authentication, returning 401 if no valid credential is present.
func (a *Auth) Required(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := a.resolveUserID(r)
		if err != nil || userID == "" {
			apierror.Unauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Auth) resolveUserID(r *http.Request) (string, error) {
	if patHeader := r.Header.Get("X-API-Token"); patHeader != "" {
		hash := sha256Hash(patHeader)
		return a.pats.GetUserIDByHash(r.Context(), hash)
	}

	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	return a.jwt.ParseAccessToken(token)
}

func sha256Hash(value string) string {
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}
```

- [ ] **Step 3: Run tests**

```bash
cd /Users/gaevivan/projects/pulse/backend && go test ./tests/...
```
Expected: all pass.

---

## Task 3: Fix receiver names — repository layer

**Files:**
- Modify: `backend/internal/repository/postgres/user/repo.go`
- Modify: `backend/internal/repository/postgres/magiclink/repo.go`
- Modify: `backend/internal/repository/postgres/refreshtoken/repo.go`
- Modify: `backend/internal/repository/postgres/pat/repo.go`

Uber rule: receivers should be 1–2 letter abbreviation.

- [ ] **Step 1: Fix user repo — replace `(repo *Repo)` → `(r *Repo)` and `repo.pool` → `r.pool`**

Full updated `repository/postgres/user/repo.go`:
```go
package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/gaevivan/pulse/internal/domain/user"
)

var _ domain.Repository = (*Repo)(nil)

// Repo is the PostgreSQL implementation of domain.Repository.
type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new user Repo.
func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	u := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, username, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return u, nil
}

func (r *Repo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	u := &domain.User{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, username, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return u, nil
}

func (r *Repo) Create(ctx context.Context, email, username string) (*domain.User, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	u := &domain.User{}
	err = tx.QueryRow(ctx,
		`INSERT INTO users (email, username) VALUES ($1, $2)
		 RETURNING id, email, username, created_at, updated_at`,
		email, username,
	).Scan(&u.ID, &u.Email, &u.Username, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert user: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO subscriptions (subject_type, subject_id, plan) VALUES ('user', $1, 'free')`,
		u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert subscription: %w", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO user_settings (user_id) VALUES ($1)`,
		u.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert user_settings: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return u, nil
}

func (r *Repo) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
		username,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("query: %w", err)
	}
	return exists, nil
}
```

Note: local variable renamed `user` → `u` to avoid shadowing the package name.

- [ ] **Step 2: Fix magiclink repo — replace `(repo *Repo)` → `(r *Repo)` and `repo.pool` → `r.pool`**

Full updated `repository/postgres/magiclink/repo.go`:
```go
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

func (r *Repo) Create(ctx context.Context, email, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO magic_link_tokens (email, token_hash, expires_at) VALUES ($1, $2, $3)`,
		email, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}

func (r *Repo) GetByHash(ctx context.Context, tokenHash string) (*domain.MagicLinkToken, error) {
	t := &domain.MagicLinkToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, token_hash, expires_at, used_at, created_at
		 FROM magic_link_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&t.ID, &t.Email, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return t, nil
}

func (r *Repo) MarkUsed(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE magic_link_tokens SET used_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
```

- [ ] **Step 3: Fix refreshtoken repo — replace `(repo *Repo)` → `(r *Repo)` and `repo.pool` → `r.pool`**

Full updated `repository/postgres/refreshtoken/repo.go`:
```go
package refreshtoken

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	domain "github.com/gaevivan/pulse/internal/domain/user"
)

var _ domain.RefreshTokenRepository = (*Repo)(nil)

// Repo is the PostgreSQL implementation of domain.RefreshTokenRepository.
type Repo struct {
	pool *pgxpool.Pool
}

// New creates a new refresh token Repo.
func New(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, userID, tokenHash string, expiresAt time.Time) (*domain.RefreshToken, error) {
	t := &domain.RefreshToken{}
	err := r.pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, token_hash, expires_at, revoked_at, created_at`,
		userID, tokenHash, expiresAt,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert: %w", err)
	}
	return t, nil
}

func (r *Repo) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshToken, error) {
	t := &domain.RefreshToken{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
		 FROM refresh_tokens WHERE token_hash = $1`,
		tokenHash,
	).Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.RevokedAt, &t.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return t, nil
}

func (r *Repo) Revoke(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE refresh_tokens SET revoked_at = now() WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
```

- [ ] **Step 4: Fix pat repo — replace `(repo *Repo)` → `(r *Repo)` and `repo.pool` → `r.pool`, add interface assertion**

Full updated `repository/postgres/pat/repo.go`:
```go
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

func (r *Repo) GetUserIDByHash(ctx context.Context, tokenHash string) (string, error) {
	var userID string
	err := r.pool.QueryRow(ctx,
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
```

- [ ] **Step 5: Build and test**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./... && go test ./tests/...
```
Expected: no errors, all tests pass.

---

## Task 4: Fix receiver names — infrastructure layer

**Files:**
- Modify: `backend/internal/infrastructure/jwt/jwt.go`
- Modify: `backend/internal/infrastructure/email/resend.go`
- Modify: `backend/internal/infrastructure/email/log.go`

- [ ] **Step 1: Fix jwt.Manager receiver — `manager` → `m`**

Full updated `infrastructure/jwt/jwt.go`:
```go
package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessTokenDuration is how long an access token is valid.
const AccessTokenDuration = 15 * time.Minute

// Manager creates and validates JWT access tokens.
type Manager struct {
	secret []byte
}

// New creates a new Manager with the given signing secret.
func New(secret string) *Manager {
	return &Manager{secret: []byte(secret)}
}

// Claims are the JWT claims for an access token.
type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed access token for the given user ID.
func (m *Manager) GenerateAccessToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseAccessToken validates the token and returns the user ID.
func (m *Manager) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}
```

- [ ] **Step 2: Fix ResendClient receiver and struct alignment**

Uber rule: struct field tags must be aligned.

Full updated `infrastructure/email/resend.go`:
```go
package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Sender sends transactional email.
type Sender interface {
	SendMagicLink(ctx context.Context, toEmail, link string) error
}

// ResendClient sends email via the Resend API.
type ResendClient struct {
	apiKey    string
	fromEmail string
	client    *http.Client
}

// NewResend creates a new ResendClient.
func NewResend(apiKey, fromEmail string) *ResendClient {
	return &ResendClient{
		apiKey:    apiKey,
		fromEmail: fromEmail,
		client:    &http.Client{},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

// SendMagicLink sends a magic link email via Resend.
func (c *ResendClient) SendMagicLink(ctx context.Context, toEmail, link string) error {
	body := resendRequest{
		From:    c.fromEmail,
		To:      []string{toEmail},
		Subject: "Your Pulse login link",
		HTML:    fmt.Sprintf(`<p>Click <a href="%s">here</a> to log in to Pulse. Link expires in 15 minutes.</p>`, link),
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resend API error: status %d", resp.StatusCode)
	}

	return nil
}
```

- [ ] **Step 3: Fix LogSender receiver — `sender` → `s`**

Full updated `infrastructure/email/log.go`:
```go
package email

import (
	"context"

	"go.uber.org/zap"
)

// LogSender prints magic links to stdout instead of sending email (non-production use).
type LogSender struct {
	log *zap.Logger
}

// NewLogSender creates a new LogSender.
func NewLogSender(log *zap.Logger) *LogSender {
	return &LogSender{log: log}
}

// SendMagicLink logs the magic link instead of sending an email.
func (s *LogSender) SendMagicLink(_ context.Context, toEmail, link string) error {
	s.log.Info("magic link (dev mode — email not sent)",
		zap.String("to", toEmail),
		zap.String("link", link),
	)
	return nil
}
```

- [ ] **Step 4: Build and test**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./... && go test ./tests/...
```
Expected: no errors, all tests pass.

---

## Task 5: Fix usecase — error wrapping, named constant, variable naming

**Files:**
- Modify: `backend/internal/usecase/user/usecase.go`

Uber rules:
- Error wrapping: wrap sentinel and cause correctly — `errors.Join(ErrSentinel, fmt.Errorf("context: %w", err))`
- Named constants, not magic numbers
- No variable named `bytes` (shadows stdlib)

- [ ] **Step 1: Rewrite usecase.go**

Full updated `usecase/user/usecase.go`:
```go
package user

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gaevivan/pulse/internal/domain/user"
	"github.com/gaevivan/pulse/internal/infrastructure/email"
	"github.com/gaevivan/pulse/internal/infrastructure/jwt"
)

const (
	magicLinkTTL        = 15 * time.Minute
	refreshTokenTTL     = 30 * 24 * time.Hour
	maxUsernameAttempts = 10
)

var (
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrTokenUsed             = errors.New("token already used")
	ErrEmailUnavailable      = errors.New("email service unavailable")
	ErrDatabaseUnavailable   = errors.New("database unavailable")
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrUsernameConflict      = errors.New("could not generate unique username")
	ErrUserNotFound          = errors.New("user not found")
)

// UseCase implements auth business logic.
type UseCase struct {
	users         user.Repository
	magicLinks    user.MagicLinkRepository
	refreshTokens user.RefreshTokenRepository
	jwt           *jwt.Manager
	email         email.Sender
	frontendURL   string
}

// New creates a new UseCase.
func New(
	users user.Repository,
	magicLinks user.MagicLinkRepository,
	refreshTokens user.RefreshTokenRepository,
	jwtManager *jwt.Manager,
	emailSender email.Sender,
	frontendURL string,
) *UseCase {
	return &UseCase{
		users:         users,
		magicLinks:    magicLinks,
		refreshTokens: refreshTokens,
		jwt:           jwtManager,
		email:         emailSender,
		frontendURL:   frontendURL,
	}
}

// SendMagicLink generates a magic link and emails it to the given address.
func (uc *UseCase) SendMagicLink(ctx context.Context, emailAddr string) error {
	rawToken, err := generateToken()
	if err != nil {
		return errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate token: %w", err))
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(magicLinkTTL)

	if err := uc.magicLinks.Create(ctx, emailAddr, tokenHash, expiresAt); err != nil {
		return errors.Join(ErrDatabaseUnavailable, fmt.Errorf("create magic link: %w", err))
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", uc.frontendURL, rawToken)
	if err := uc.email.SendMagicLink(ctx, emailAddr, link); err != nil {
		return errors.Join(ErrEmailUnavailable, fmt.Errorf("send email: %w", err))
	}

	return nil
}

// VerifyResult holds the tokens and user returned after successful magic-link verification.
type VerifyResult struct {
	AccessToken  string
	RefreshToken string
	User         *user.User
}

// VerifyMagicLink validates the raw token and returns a token pair and the user.
func (uc *UseCase) VerifyMagicLink(ctx context.Context, rawToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawToken)

	magicToken, err := uc.magicLinks.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get magic link: %w", err))
	}
	if magicToken == nil {
		return nil, ErrInvalidToken
	}
	if magicToken.UsedAt != nil {
		return nil, ErrTokenUsed
	}
	if time.Now().After(magicToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	existingUser, err := uc.users.GetByEmail(ctx, magicToken.Email)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get user: %w", err))
	}

	currentUser := existingUser
	if currentUser == nil {
		username, err := uc.generateUsername(ctx, magicToken.Email)
		if err != nil {
			return nil, fmt.Errorf("generate username: %w", err)
		}
		currentUser, err = uc.users.Create(ctx, magicToken.Email, username)
		if err != nil {
			return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("create user: %w", err))
		}
	}

	if err := uc.magicLinks.MarkUsed(ctx, magicToken.ID); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("mark used: %w", err))
	}

	return uc.issueTokenPair(ctx, currentUser)
}

// Refresh rotates the refresh token and returns a new token pair.
func (uc *UseCase) Refresh(ctx context.Context, rawRefreshToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := uc.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get refresh token: %w", err))
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	if err := uc.refreshTokens.Revoke(ctx, storedToken.ID); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("revoke old token: %w", err))
	}

	currentUser, err := uc.users.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("get user: %w", err))
	}
	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	return uc.issueTokenPair(ctx, currentUser)
}

// Logout revokes the given refresh token. Silently succeeds for unknown tokens.
func (uc *UseCase) Logout(ctx context.Context, rawRefreshToken string) error {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := uc.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("get refresh token: %w", err)
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return ErrInvalidToken
	}

	return uc.refreshTokens.Revoke(ctx, storedToken.ID)
}

// GetUserByID returns the user with the given ID, or nil if not found.
func (uc *UseCase) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return uc.users.GetByID(ctx, id)
}

func (uc *UseCase) issueTokenPair(ctx context.Context, u *user.User) (*VerifyResult, error) {
	accessToken, err := uc.jwt.GenerateAccessToken(u.ID)
	if err != nil {
		return nil, errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate access token: %w", err))
	}

	rawRefreshToken, err := generateToken()
	if err != nil {
		return nil, errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate refresh token: %w", err))
	}

	refreshTokenHash := hashToken(rawRefreshToken)
	expiresAt := time.Now().Add(refreshTokenTTL)

	if _, err := uc.refreshTokens.Create(ctx, u.ID, refreshTokenHash, expiresAt); err != nil {
		return nil, errors.Join(ErrDatabaseUnavailable, fmt.Errorf("store refresh token: %w", err))
	}

	return &VerifyResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		User:         u,
	}, nil
}

func (uc *UseCase) generateUsername(ctx context.Context, emailAddr string) (string, error) {
	prefix := strings.ToLower(strings.Split(emailAddr, "@")[0])

	exists, err := uc.users.ExistsByUsername(ctx, prefix)
	if err != nil {
		return "", errors.Join(ErrDatabaseUnavailable, fmt.Errorf("check username: %w", err))
	}
	if !exists {
		return prefix, nil
	}

	for range maxUsernameAttempts {
		suffix, err := generateToken()
		if err != nil {
			return "", errors.Join(ErrTokenGenerationFailed, fmt.Errorf("generate suffix: %w", err))
		}
		candidate := prefix + "_" + suffix[:4]
		exists, err := uc.users.ExistsByUsername(ctx, candidate)
		if err != nil {
			return "", errors.Join(ErrDatabaseUnavailable, fmt.Errorf("check username: %w", err))
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", ErrUsernameConflict
}

func generateToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buf), nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
```

- [ ] **Step 2: Build and test**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./... && go test ./tests/...
```
Expected: no errors, all tests pass.

---

## Task 6: Fix validator — unexported package-level var naming

**Files:**
- Modify: `backend/pkg/validator/validator.go`

Uber rule: "Unexported globals: `var _cache = ...`" — underscore prefix for unexported package-level vars.

- [ ] **Step 1: Rename `emailRegex` → `_emailRegex`**

Full updated `pkg/validator/validator.go`:
```go
package validator

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

var _emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Decode decodes the JSON body of r into destination, disallowing unknown fields.
func Decode(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return errors.New("invalid request body")
	}
	return nil
}

// Email returns an error if the email address is empty or malformed.
func Email(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}
	if !_emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// Required returns an error if value is empty after trimming whitespace.
func Required(value, field string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New(field + " is required")
	}
	return nil
}
```

- [ ] **Step 2: Build and test**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./... && go test ./tests/...
```
Expected: no errors, all tests pass.

---

## Task 7: Fix seed.go — `var err` at top with `=` reuse

**Files:**
- Modify: `backend/cmd/seed/main.go`

Uber rule: "Handle errors once, immediately with `:=`" and "Local var: `:=` for most". Declaring `var err error` at the top of a function and reusing with `=` is the "Don't" pattern.

- [ ] **Step 1: Rewrite seed function to use `:=` per operation**

Full updated `cmd/seed/main.go`:
```go
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
	var aliceID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ('alice@example.com', 'alice')
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`).Scan(&aliceID); err != nil {
		return fmt.Errorf("insert alice: %w", err)
	}

	var bobID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ('bob@example.com', 'bob')
		ON CONFLICT (email) DO UPDATE SET email = EXCLUDED.email
		RETURNING id
	`).Scan(&bobID); err != nil {
		return fmt.Errorf("insert bob: %w", err)
	}

	// Subscriptions
	if _, err := tx.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('user', $1, 'pro'), ('user', $2, 'free')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, aliceID, bobID); err != nil {
		return fmt.Errorf("insert subscriptions: %w", err)
	}

	// Team
	var teamID string
	if err := tx.QueryRow(ctx, `
		INSERT INTO teams (name, prefix, owner_id)
		VALUES ('Pulse Dev', 'PLS', $1)
		ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
		RETURNING id
	`, aliceID).Scan(&teamID); err != nil {
		return fmt.Errorf("insert team: %w", err)
	}

	// Team subscription
	if _, err := tx.Exec(ctx, `
		INSERT INTO subscriptions (subject_type, subject_id, plan)
		VALUES ('team', $1, 'team')
		ON CONFLICT (subject_type, subject_id) DO NOTHING
	`, teamID); err != nil {
		return fmt.Errorf("insert team subscription: %w", err)
	}

	// Team members
	if _, err := tx.Exec(ctx, `
		INSERT INTO team_members (team_id, user_id)
		VALUES ($1, $2), ($1, $3)
		ON CONFLICT DO NOTHING
	`, teamID, aliceID, bobID); err != nil {
		return fmt.Errorf("insert team members: %w", err)
	}

	// Labels
	if _, err := tx.Exec(ctx, `
		INSERT INTO labels (owner_type, owner_id, name, color)
		VALUES
			('team', $1, 'bug', '#d73a4a'),
			('team', $1, 'feature', '#0075ca'),
			('team', $1, 'chore', '#e4e669')
		ON CONFLICT (owner_type, owner_id, name) DO NOTHING
	`, teamID); err != nil {
		return fmt.Errorf("insert labels: %w", err)
	}

	// Task sequence
	if _, err := tx.Exec(ctx, `
		INSERT INTO task_sequences (owner_type, owner_id, last_number)
		VALUES ('team', $1, 0)
		ON CONFLICT DO NOTHING
	`, teamID); err != nil {
		return fmt.Errorf("insert task sequence: %w", err)
	}

	// Tasks
	var taskNumber int64
	if err := tx.QueryRow(ctx, `
		UPDATE task_sequences SET last_number = last_number + 1
		WHERE owner_type = 'team' AND owner_id = $1
		RETURNING last_number
	`, teamID).Scan(&taskNumber); err != nil {
		return fmt.Errorf("increment task sequence: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO tasks (key_number, owner_type, owner_id, title, status, created_by)
		VALUES ($1, 'team', $2, 'Настроить CI/CD', 'opened', $3)
		ON CONFLICT (owner_type, owner_id, key_number) DO NOTHING
	`, taskNumber, teamID, aliceID); err != nil {
		return fmt.Errorf("insert task: %w", err)
	}

	// User settings
	if _, err := tx.Exec(ctx, `
		INSERT INTO user_settings (user_id, language, theme)
		VALUES ($1, 'ru', 'dark'), ($2, 'en', 'system')
		ON CONFLICT (user_id) DO NOTHING
	`, aliceID, bobID); err != nil {
		return fmt.Errorf("insert user settings: %w", err)
	}

	return tx.Commit(ctx)
}
```

- [ ] **Step 2: Build and test**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./... && go test ./tests/...
```
Expected: no errors, all tests pass.

---

## Task 8: Final build verification

- [ ] **Step 1: Full build**

```bash
cd /Users/gaevivan/projects/pulse/backend && go build ./...
```
Expected: no errors.

- [ ] **Step 2: Full test run**

```bash
cd /Users/gaevivan/projects/pulse/backend && go test ./...
```
Expected: all tests pass (db tests may be skipped without a running database — that is acceptable).

- [ ] **Step 3: Vet**

```bash
cd /Users/gaevivan/projects/pulse/backend && go vet ./...
```
Expected: no warnings.
