# Code Conventions

## Git Commits

Каждый коммит, завершающий реализацию задачи из GitHub Issues, должен содержать closing keyword с номером задачи:

```
feat: add magic link auth

closes #3
```

Поддерживаемые ключевые слова: `closes`, `fixes`, `resolves`. При мерже в main GitHub автоматически закроет указанную задачу.

Если коммит частично реализует задачу, использовать `ref #N` вместо `closes #N` — задача останется открытой.

Если один коммит закрывает несколько задач, перечислить каждую на отдельной строке:

```
feat: add task CRUD and labels API

closes #11
closes #10
```

## Именование переменных и аргументов (Go)

Не использовать однобуквенные или сокращённые имена переменных, аргументов и receiver-ов. Имена должны быть полными и осмысленными.

```go
// Плохо
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) { ... }
func (r *Repo) GetByEmail(ctx context.Context, email string) { ... }
func (m *Manager) GenerateAccessToken(userID string) { ... }

// Хорошо
func (handler *AuthHandler) Me(writer http.ResponseWriter, request *http.Request) { ... }
func (repo *Repo) GetByEmail(ctx context.Context, email string) { ... }
func (manager *Manager) GenerateAccessToken(userID string) { ... }
```

Это касается любого Go-кода: handler-ов, middleware, репозиториев, usecase-ов, инфраструктурных компонентов.

## API-типы (frontend)

Zod-схемы и типизированный API-клиент генерируются из OpenAPI-схемы бэкенда (`openapi-zod-client`).

Рабочий процесс:
1. `make swagger` — обновить схему из аннотаций Go и сконвертировать в OpenAPI 3.0
2. `make generate-types` — сгенерировать `frontend/src/api/generated.ts`

Файл `generated.ts` коммитить не нужно — он генерируется при сборке.
При добавлении нового эндпоинта: обновить аннотации в Go-обработчике, запустить оба шага.
Клиент доступен через `api` из `src/api/client.ts`, методы соответствуют алиасам из схемы (`api.postV1authmagicLink(...)`).
