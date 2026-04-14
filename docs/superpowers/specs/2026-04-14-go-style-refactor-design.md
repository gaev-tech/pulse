# Go Style Refactor Design

**Date:** 2026-04-14
**Scope:** All 35 `.go` files in `backend/`
**Goal:** Apply Uber Go Style Guide rules across the entire Go backend via style-only refactoring (no feature changes, no architectural changes).

---

## Approach

Single pass over all backend Go files using the `go-style` skill as the authoritative reference. Changes are grouped by category below. Business logic, SQL queries, and layer interfaces may be touched if they contain style violations.

Tests (`backend/tests/`) are run before and after to verify behaviour is unchanged.

---

## Change Categories

### 1. Error Handling

- Replace `fmt.Errorf("...: %w", errors.Join(SentinelErr, err))` with idiomatic wrapping. When a sentinel and a cause both need to be surfaced, use `errors.Join(SentinelErr, fmt.Errorf("context: %w", err))` or restructure so the sentinel is returned directly and the cause is logged.
- Ensure all error returns include context via `fmt.Errorf("operationName: %w", err)`.
- Avoid returning `nil, nil` where a typed "not found" sentinel is cleaner.

### 2. Naming

- Receiver names: shorten to 1–2 letter abbreviations (`h` for handler, `uc` for usecase, `r` for repo, etc.) consistently across all files.
- No single-letter variable names except loop indices and well-known math vars.
- No abbreviations like `f` for filter — use full names.
- Exported constants grouped as `iota` blocks where applicable.

### 3. Nil Checks

- Consolidate double nil checks: `if err != nil { ... } if x == nil { ... }` → `if err != nil || x == nil { ... }` where appropriate.
- Remove redundant nil guards after early returns.

### 4. Structs and Types

- Export JSON-marshaling-only structs (`detail`, `response` in `apierror`) or inline them — unexported types used only for `json.Marshal` are confusing.
- Group related struct fields logically (required first, optional last).

### 5. JSON Encoding Errors

- Document or log ignored `json.Encoder.Encode()` errors; don't silently discard with bare `_ =`.

### 6. Godoc Comments

- Add godoc comments on all exported types, functions, and methods that lack them.
- Comments must start with the name of the thing being documented.

### 7. Package-Level State

- No unexported mutable package-level vars unless they are initialized once at startup and never written again.
- No `init()` functions (Uber style: explicit initialization in `main`).

### 8. Config Validation

- `config.Load()` should validate required fields (e.g., JWT secret must not be `"changeme"` in non-dev environments, Resend API key must be set if email sender is Resend).

### 9. Username Generation

- Log each failed attempt in `generateUsername()` at debug level.
- Consider replacing magic number `10` with a named constant `maxUsernameAttempts`.

### 10. Miscellaneous

- Ensure `defer tx.Rollback()` errors are captured: `defer func() { _ = tx.Rollback(ctx) }()`.
- Remove any dead code or unreachable branches.
- Consistent import grouping: stdlib → external → internal, separated by blank lines.

---

## Files in Scope

```
backend/cmd/api/main.go
backend/cmd/seed/main.go
backend/internal/domain/user/user.go
backend/internal/domain/team/team.go
backend/internal/domain/task/task.go
backend/internal/domain/filter/filter.go
backend/internal/domain/permission/permission.go
backend/internal/domain/event/event.go
backend/internal/handler/v1/auth.go
backend/internal/handler/v1/router.go
backend/internal/handler/middleware/auth.go
backend/internal/usecase/user/usecase.go
backend/internal/usecase/user/user.go
backend/internal/usecase/team/team.go
backend/internal/usecase/task/task.go
backend/internal/usecase/filter/filter.go
backend/internal/usecase/permission/permission.go
backend/internal/repository/postgres/user/repo.go
backend/internal/repository/postgres/magic_link/repo.go
backend/internal/repository/postgres/refresh_token/repo.go
backend/internal/repository/postgres/pat/repo.go
backend/internal/repository/migrations/migrations.go
backend/internal/infrastructure/config/config.go
backend/internal/infrastructure/jwt/jwt.go
backend/internal/infrastructure/logger/logger.go
backend/internal/infrastructure/postgres/db.go
backend/internal/infrastructure/postgres/migrate.go
backend/internal/infrastructure/email/resend.go
backend/internal/infrastructure/email/log.go
backend/pkg/apierror/apierror.go
backend/pkg/pagination/pagination.go
backend/pkg/validator/validator.go
backend/api/docs/docs.go
backend/tests/health_test.go
backend/tests/auth_test.go
backend/tests/db_test.go
```

---

## Success Criteria

- All existing tests pass after refactoring.
- No `golangci-lint` violations for style rules covered by the go-style skill.
- No behaviour changes — only style and readability improvements.
