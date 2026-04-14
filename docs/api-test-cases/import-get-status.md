# GET /api/v1/import/{importID}

## 1. Успешное получение статуса — импорт в процессе (JWT)
- Auth: JWT
- Preconditions: импорт запущен, ещё не завершён
- Expected: 200
- Response: { id, source: "csv", status: "in_progress", progress: { total: 100, processed: 42 }, result: null }

## 2. Успешное получение статуса — импорт завершён
- Auth: JWT
- Preconditions: импорт завершён успешно
- Expected: 200
- Response: { id, source: "jira", status: "completed", progress: { total: 50, processed: 50 }, result: { imported: 45, updated: 5, errors: [] } }

## 3. Успешное получение статуса — импорт завершён с ошибками
- Auth: JWT
- Preconditions: импорт завершён, часть задач не импортирована
- Expected: 200
- Response: { id, source: "github", status: "completed", progress: { total: 30, processed: 30 }, result: { imported: 28, updated: 0, errors: [{ origin_id: "github:acme/app:5", error: "..." }] } }

## 4. Успешное получение статуса (PAT)
- Auth: PAT
- Preconditions: импорт создан текущим пользователем
- Expected: 200

## 5. Импорт принадлежит другому пользователю
- Auth: JWT
- Preconditions: importID принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Несуществующий importID
- Auth: JWT
- URL: /api/v1/import/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 7. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
