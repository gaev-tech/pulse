# PATCH /api/v1/filters/{filterID}/settings

## 1. Успешное изменение columns (JWT)
- Auth: JWT
- Body: { "columns": ["labels", "assignee", "status", "created_at"] }
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { columns: ["labels", "assignee", "status", "created_at"], sort1_column, sort1_dir, sort2_column, sort2_dir }

## 2. Успешное изменение настроек (PAT)
- Auth: PAT
- Body: { "columns": ["status"] }
- Expected: 200

## 3. Изменение первичной сортировки
- Auth: JWT
- Body: { "sort1_column": "created_at", "sort1_dir": "asc" }
- Expected: 200
- Response: { columns, sort1_column: "created_at", sort1_dir: "asc", sort2_column, sort2_dir }

## 4. Изменение вторичной сортировки
- Auth: JWT
- Body: { "sort2_column": "title", "sort2_dir": "desc" }
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column: "title", sort2_dir: "desc" }

## 5. Сброс вторичной сортировки
- Auth: JWT
- Body: { "sort2_column": null, "sort2_dir": null }
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column: null, sort2_dir: null }

## 6. Изменение всех настроек одновременно
- Auth: JWT
- Body: { "columns": ["assignee", "status"], "sort1_column": "key", "sort1_dir": "desc", "sort2_column": "assignee", "sort2_dir": "asc" }
- Expected: 200

## 7. Командный фильтр — с правом team.manage_filters
- Auth: JWT
- Body: { "columns": ["status"] }
- Preconditions: пользователь имеет team.manage_filters
- Expected: 200

## 8. Командный фильтр — без права team.manage_filters
- Auth: JWT
- Body: { "columns": ["status"] }
- Preconditions: пользователь — участник без team.manage_filters
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 9. Командный фильтр — не участник команды
- Auth: JWT
- Body: { "columns": ["status"] }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 10. Чужой личный фильтр
- Auth: JWT
- Body: { "columns": ["status"] }
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Невалидное имя колонки
- Auth: JWT
- Body: { "columns": ["invalid_column"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Невалидное направление сортировки
- Auth: JWT
- Body: { "sort1_dir": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Сортировка по колонке, не входящей в columns
- Auth: JWT
- Body: { "sort1_column": "assignee" }
- Preconditions: "assignee" не входит в текущий набор columns
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid/settings
- Body: { "columns": ["status"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 15. Без авторизации
- Auth: нет
- Body: { "columns": ["status"] }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 16. Пустое тело запроса
- Auth: JWT
- Body: {}
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir } (без изменений)
