# Error Handling Design — Backend Auth

## Контекст

Все неожиданные ошибки в auth-хэндлерах сейчас возвращают 500 без полезной информации. Цель — заменить на конкретные коды, чтобы фронтенд мог реагировать дифференцированно.

## Карта ошибок

| Место | Причина | Код | HTTP |
|---|---|---|---|
| `SendMagicLink` | email-сервис недоступен | `EMAIL_UNAVAILABLE` | 500 |
| `SendMagicLink` | `rand.Read` упал | `TOKEN_GENERATION_FAILED` | 500 |
| `SendMagicLink` | `magicLinks.Create` упал | `DATABASE_UNAVAILABLE` | 503 |
| `VerifyMagicLink` | любой `db.*` упал | `DATABASE_UNAVAILABLE` | 503 |
| `VerifyMagicLink` | `jwt.Generate` / `rand.Read` | `TOKEN_GENERATION_FAILED` | 500 |
| `VerifyMagicLink` | username не сгенерировался за 10 попыток | `USERNAME_CONFLICT` | 500 |
| `Refresh` | любой `db.*` упал | `DATABASE_UNAVAILABLE` | 503 |
| `Refresh` | `GetByID` → nil (юзер удалён) | `USER_NOT_FOUND` | 404 |
| `Refresh` | `jwt.Generate` / `rand.Read` | `TOKEN_GENERATION_FAILED` | 500 |
| `Me` | `GetByID` упал | `DATABASE_UNAVAILABLE` | 503 |
| `Me` | `GetByID` → nil | `USER_NOT_FOUND` | 404 |

## Архитектура

### 1. Usecase — новые sentinel-ошибки

`backend/internal/usecase/user/usecase.go`:

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

Все ошибки оборачиваются через `errors.Join(sentinel, original)` чтобы сохранить оригинальный стек и позволить `errors.Is` работать.

### 2. apierror — новые хелперы

`backend/pkg/apierror/apierror.go`:

```go
// Существует, но возвращает 403 — меняем на 500
func EmailUnavailable(w http.ResponseWriter) {
    Write(w, 500, "EMAIL_UNAVAILABLE", "email service is not available")
}
func DatabaseUnavailable(w http.ResponseWriter) {
    Write(w, 503, "DATABASE_UNAVAILABLE", "database is unavailable, try again later")
}
func TokenGenerationFailed(w http.ResponseWriter) {
    Write(w, 500, "TOKEN_GENERATION_FAILED", "failed to generate secure token")
}
func UsernameConflict(w http.ResponseWriter) {
    Write(w, 500, "USERNAME_CONFLICT", "could not generate unique username")
}
// UserNotFound уже есть как NotFound, но добавить отдельный для явности
func UserNotFound(w http.ResponseWriter) {
    Write(w, 404, "USER_NOT_FOUND", "user not found")
}
```

### 3. Handler — замена apierror.Internal

Каждый `apierror.Internal` заменяется на цепочку `errors.Is`:

```go
switch {
case errors.Is(err, userusecase.ErrDatabaseUnavailable):
    apierror.DatabaseUnavailable(writer)
case errors.Is(err, userusecase.ErrTokenGenerationFailed):
    apierror.TokenGenerationFailed(writer)
// ...
default:
    apierror.Internal(writer) // fallback для непредвиденного
}
```

## Изменяемые файлы

- `backend/internal/usecase/user/usecase.go` — новые ошибки + обёртка на источниках
- `backend/pkg/apierror/apierror.go` — новые хелперы
- `backend/internal/handler/v1/auth.go` — замена `apierror.Internal` на конкретные

## Проверка

- `make test-api` — все тесты зелёные
- `go build ./...` — компилируется
- `golangci-lint run ./...` — без замечаний
