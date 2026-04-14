# Error Handling Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Заменить все `apierror.Internal` в auth-хэндлерах на конкретные коды ошибок с полезной мета-информацией.

**Architecture:** Три слоя изменений: (1) sentinel-ошибки в usecase, (2) новые хелперы в apierror, (3) маппинг в хэндлере. Каждый источник ошибки оборачивается через `errors.Join(sentinel, original)` — сохраняется стек, работает `errors.Is`.

**Tech Stack:** Go 1.23, `errors.Join`, `errors.Is`

---

## File Map

- Modify: `backend/pkg/apierror/apierror.go` — обновить `EmailUnavailable` (403→500), добавить хелперы
- Modify: `backend/internal/usecase/user/usecase.go` — новые sentinel-ошибки, обёртка на источниках
- Modify: `backend/internal/handler/v1/auth.go` — заменить `apierror.Internal` на конкретные хелперы

---

## Task 1: apierror — обновить EmailUnavailable и добавить хелперы

**Files:**
- Modify: `backend/pkg/apierror/apierror.go`

- [ ] **Step 1: Обновить `EmailUnavailable` с 403 на 500, добавить новые хелперы**

Итоговый файл `backend/pkg/apierror/apierror.go`:

```go
package apierror

import (
	"encoding/json"
	"net/http"
)

type detail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type response struct {
	Error detail `json:"error"`
}

func Write(writer http.ResponseWriter, status int, code, message string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(response{
		Error: detail{Code: code, Message: message},
	})
}

func BadRequest(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusBadRequest, "VALIDATION_ERROR", message)
}

func Unauthorized(writer http.ResponseWriter) {
	Write(writer, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
}

func Forbidden(writer http.ResponseWriter) {
	Write(writer, http.StatusForbidden, "PERMISSION_DENIED", "permission denied")
}

func NotFound(writer http.ResponseWriter) {
	Write(writer, http.StatusNotFound, "NOT_FOUND", "not found")
}

func UserNotFound(writer http.ResponseWriter) {
	Write(writer, http.StatusNotFound, "USER_NOT_FOUND", "user not found")
}

func Conflict(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusConflict, "CONFLICT", message)
}

func QuotaExceeded(writer http.ResponseWriter, message string) {
	Write(writer, http.StatusPaymentRequired, "QUOTA_EXCEEDED", message)
}

func EmailUnavailable(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "EMAIL_UNAVAILABLE", "email service is not available")
}

func DatabaseUnavailable(writer http.ResponseWriter) {
	Write(writer, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "database is unavailable, try again later")
}

func TokenGenerationFailed(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "failed to generate secure token")
}

func UsernameConflict(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "USERNAME_CONFLICT", "could not generate unique username")
}

func Internal(writer http.ResponseWriter) {
	Write(writer, http.StatusInternalServerError, "INTERNAL", "internal server error")
}
```

- [ ] **Step 2: Собрать пакет**

```bash
cd backend && go build ./pkg/apierror/...
```

Ожидание: без ошибок.

---

## Task 2: Usecase — sentinel-ошибки и обёртка на источниках

**Files:**
- Modify: `backend/internal/usecase/user/usecase.go`

- [ ] **Step 1: Добавить новые sentinel-ошибки**

Заменить блок `var (...)` в начале файла:

```go
var (
	ErrInvalidToken          = errors.New("invalid or expired token")
	ErrTokenUsed             = errors.New("token already used")
	ErrEmailUnavailable      = errors.New("email service unavailable")
	ErrDatabaseUnavailable   = errors.New("database unavailable")
	ErrTokenGenerationFailed = errors.New("token generation failed")
	ErrUsernameConflict      = errors.New("could not generate unique username")
	ErrUserNotFound          = errors.New("user not found")
)
```

- [ ] **Step 2: Обернуть ошибки в `SendMagicLink`**

```go
func (useCase *UseCase) SendMagicLink(ctx context.Context, emailAddr string) error {
	rawToken, err := generateToken()
	if err != nil {
		return fmt.Errorf("generate token: %w", errors.Join(ErrTokenGenerationFailed, err))
	}

	tokenHash := hashToken(rawToken)
	expiresAt := time.Now().Add(magicLinkTTL)

	if err := useCase.magicLinks.Create(ctx, emailAddr, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("create magic link: %w", errors.Join(ErrDatabaseUnavailable, err))
	}

	link := fmt.Sprintf("%s/auth/verify?token=%s", useCase.frontendURL, rawToken)
	if err := useCase.email.SendMagicLink(ctx, emailAddr, link); err != nil {
		return fmt.Errorf("send email: %w", errors.Join(ErrEmailUnavailable, err))
	}

	return nil
}
```

- [ ] **Step 3: Обернуть ошибки в `VerifyMagicLink`**

```go
func (useCase *UseCase) VerifyMagicLink(ctx context.Context, rawToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawToken)

	magicToken, err := useCase.magicLinks.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("get magic link: %w", errors.Join(ErrDatabaseUnavailable, err))
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

	existingUser, err := useCase.users.GetByEmail(ctx, magicToken.Email)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", errors.Join(ErrDatabaseUnavailable, err))
	}

	var currentUser *user.User
	if existingUser != nil {
		currentUser = existingUser
	} else {
		username, err := useCase.generateUsername(ctx, magicToken.Email)
		if err != nil {
			return nil, fmt.Errorf("generate username: %w", err)
		}
		currentUser, err = useCase.users.Create(ctx, magicToken.Email, username)
		if err != nil {
			return nil, fmt.Errorf("create user: %w", errors.Join(ErrDatabaseUnavailable, err))
		}
	}

	if err := useCase.magicLinks.MarkUsed(ctx, magicToken.ID); err != nil {
		return nil, fmt.Errorf("mark used: %w", errors.Join(ErrDatabaseUnavailable, err))
	}

	return useCase.issueTokenPair(ctx, currentUser)
}
```

- [ ] **Step 4: Обернуть ошибки в `Refresh`**

```go
func (useCase *UseCase) Refresh(ctx context.Context, rawRefreshToken string) (*VerifyResult, error) {
	tokenHash := hashToken(rawRefreshToken)

	storedToken, err := useCase.refreshTokens.GetByHash(ctx, tokenHash)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", errors.Join(ErrDatabaseUnavailable, err))
	}
	if storedToken == nil || storedToken.RevokedAt != nil {
		return nil, ErrInvalidToken
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return nil, ErrInvalidToken
	}

	if err := useCase.refreshTokens.Revoke(ctx, storedToken.ID); err != nil {
		return nil, fmt.Errorf("revoke old token: %w", errors.Join(ErrDatabaseUnavailable, err))
	}

	currentUser, err := useCase.users.GetByID(ctx, storedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", errors.Join(ErrDatabaseUnavailable, err))
	}
	if currentUser == nil {
		return nil, ErrUserNotFound
	}

	return useCase.issueTokenPair(ctx, currentUser)
}
```

- [ ] **Step 5: Обернуть ошибки в `issueTokenPair`**

```go
func (useCase *UseCase) issueTokenPair(ctx context.Context, currentUser *user.User) (*VerifyResult, error) {
	accessToken, err := useCase.jwt.GenerateAccessToken(currentUser.ID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", errors.Join(ErrTokenGenerationFailed, err))
	}

	rawRefreshToken, err := generateToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", errors.Join(ErrTokenGenerationFailed, err))
	}

	refreshTokenHash := hashToken(rawRefreshToken)
	expiresAt := time.Now().Add(refreshTokenTTL)

	if _, err := useCase.refreshTokens.Create(ctx, currentUser.ID, refreshTokenHash, expiresAt); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", errors.Join(ErrDatabaseUnavailable, err))
	}

	return &VerifyResult{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		User:         currentUser,
	}, nil
}
```

- [ ] **Step 6: Обернуть ошибки в `generateUsername`**

Заменить финальный `return "", fmt.Errorf("could not generate unique username")`:

```go
func (useCase *UseCase) generateUsername(ctx context.Context, emailAddr string) (string, error) {
	prefix := strings.Split(emailAddr, "@")[0]
	prefix = strings.ToLower(prefix)

	exists, err := useCase.users.ExistsByUsername(ctx, prefix)
	if err != nil {
		return "", fmt.Errorf("check username: %w", errors.Join(ErrDatabaseUnavailable, err))
	}
	if !exists {
		return prefix, nil
	}

	for range 10 {
		suffix, err := generateToken()
		if err != nil {
			return "", fmt.Errorf("generate suffix: %w", errors.Join(ErrTokenGenerationFailed, err))
		}
		candidate := prefix + "_" + suffix[:4]
		exists, err := useCase.users.ExistsByUsername(ctx, candidate)
		if err != nil {
			return "", fmt.Errorf("check username: %w", errors.Join(ErrDatabaseUnavailable, err))
		}
		if !exists {
			return candidate, nil
		}
	}

	return "", ErrUsernameConflict
}
```

- [ ] **Step 7: Собрать usecase**

```bash
cd backend && go build ./internal/usecase/...
```

Ожидание: без ошибок.

---

## Task 3: Handler — заменить apierror.Internal на конкретные хелперы

**Files:**
- Modify: `backend/internal/handler/v1/auth.go`

- [ ] **Step 1: Обновить `SendMagicLink` в хэндлере**

```go
if err := handler.usecase.SendMagicLink(request.Context(), body.Email); err != nil {
    switch {
    case errors.Is(err, userusecase.ErrEmailUnavailable):
        apierror.EmailUnavailable(writer)
    case errors.Is(err, userusecase.ErrDatabaseUnavailable):
        apierror.DatabaseUnavailable(writer)
    case errors.Is(err, userusecase.ErrTokenGenerationFailed):
        apierror.TokenGenerationFailed(writer)
    default:
        apierror.Internal(writer)
    }
    return
}
```

- [ ] **Step 2: Обновить `VerifyMagicLink` в хэндлере**

```go
result, err := handler.usecase.VerifyMagicLink(request.Context(), body.Token)
if err != nil {
    switch {
    case errors.Is(err, userusecase.ErrInvalidToken), errors.Is(err, userusecase.ErrTokenUsed):
        apierror.Unauthorized(writer)
    case errors.Is(err, userusecase.ErrDatabaseUnavailable):
        apierror.DatabaseUnavailable(writer)
    case errors.Is(err, userusecase.ErrTokenGenerationFailed):
        apierror.TokenGenerationFailed(writer)
    case errors.Is(err, userusecase.ErrUsernameConflict):
        apierror.UsernameConflict(writer)
    default:
        apierror.Internal(writer)
    }
    return
}
```

- [ ] **Step 3: Обновить `Refresh` в хэндлере**

```go
result, err := handler.usecase.Refresh(request.Context(), body.RefreshToken)
if err != nil {
    switch {
    case errors.Is(err, userusecase.ErrInvalidToken):
        apierror.Unauthorized(writer)
    case errors.Is(err, userusecase.ErrUserNotFound):
        apierror.UserNotFound(writer)
    case errors.Is(err, userusecase.ErrDatabaseUnavailable):
        apierror.DatabaseUnavailable(writer)
    case errors.Is(err, userusecase.ErrTokenGenerationFailed):
        apierror.TokenGenerationFailed(writer)
    default:
        apierror.Internal(writer)
    }
    return
}
```

- [ ] **Step 4: Обновить `Me` в хэндлере**

```go
currentUser, err := handler.usecase.GetUserByID(request.Context(), userID)
if err != nil {
    apierror.DatabaseUnavailable(writer)
    return
}
if currentUser == nil {
    apierror.UserNotFound(writer)
    return
}
```

- [ ] **Step 5: Собрать и прогнать тесты**

```bash
cd backend && go build -o /tmp/pulse-api ./cmd/api/... && go test ./tests/... -v
```

Ожидание: все тесты зелёные.

- [ ] **Step 6: Линт**

```bash
cd backend && golangci-lint run ./...
```

Ожидание: 0 issues.
