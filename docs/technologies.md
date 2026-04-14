# Technologies

## Stack

| Layer | Technology |
|---|---|
| Backend | Go |
| Frontend | SvelteKit + TypeScript |
| Database | PostgreSQL 16 |
| Server | Nginx |
| Containerization | Docker + docker-compose |
| Build scripts | Make |
| E2E tests | Playwright |
| API docs | Swagger / OpenAPI (swaggo/http-swagger) |
| File storage | MinIO |

## External services

| Service | Purpose |
|---|---|
| Resend | Отправка email (magic link). Free: 3 000 писем/мес |
| ЮKassa | Рекуррентные платежи и подписки (российский рынок) |

## Backend dependencies

| Package | Purpose |
|---|---|
| `go-chi/chi` | HTTP router |
| `jackc/pgx/v5` | PostgreSQL driver |
| `golang-migrate/migrate` | Database migrations |
| `golang-jwt/jwt` | JWT tokens |
| `uber-go/zap` | Structured logging |
| `golang.org/x/crypto/bcrypt` | Password hashing |
| `swaggo/swag` | OpenAPI spec generation |
| `cosmtrek/air` | Hot reload in dev |
| `minio/minio-go` | MinIO / S3 client |

## Frontend dependencies

| Package | Purpose |
|---|---|
| SvelteKit | Framework + file-based routing |
| TypeScript | Type safety |
| Vite | Build tool + dev server |
| Playwright | E2E testing |
| `@milkdown/crepe` | Markdown editor |
