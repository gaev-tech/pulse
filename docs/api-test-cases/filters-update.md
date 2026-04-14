# PATCH /api/v1/filters/{filterID}

## 1. Успешное изменение name личного фильтра (JWT)
- Auth: JWT
- Body: { "name": "Updated Filter" }
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { id, name: "Updated Filter", owner_type, owner_id, filter_mode, ... }

## 2. Успешное изменение name фильтра (PAT)
- Auth: PAT
- Body: { "name": "Updated Filter" }
- Expected: 200

## 3. Изменение критериев simple mode
- Auth: JWT
- Body: { "search_contains": "new search", "assignee_ids": ["<user_id>"], "status": "closed", "label_ids": ["<label_id>"] }
- Preconditions: фильтр в simple mode
- Expected: 200
- Response: { filter_mode: "simple", search_contains: "new search", ... }

## 4. Изменение rsql выражения
- Auth: JWT
- Body: { "rsql": "status==closed;title==*bug*" }
- Preconditions: фильтр в rsql mode
- Expected: 200
- Response: { filter_mode: "rsql", rsql: "status==closed;title==*bug*" }

## 5. Переключение filter_mode с simple на rsql
- Auth: JWT
- Body: { "filter_mode": "rsql", "rsql": "status==opened" }
- Preconditions: фильтр в simple mode
- Expected: 200
- Response: { filter_mode: "rsql", rsql: "status==opened" }
- Side effects: критерии simple mode сбрасываются

## 6. Переключение filter_mode с rsql на simple
- Auth: JWT
- Body: { "filter_mode": "simple", "status": "opened" }
- Preconditions: фильтр в rsql mode
- Expected: 200
- Response: { filter_mode: "simple", status: "opened" }
- Side effects: rsql критерий сбрасывается

## 7. Личный фильтр — изменение team_id в simple mode
- Auth: JWT
- Body: { "team_id": "<team_id>" }
- Preconditions: личный фильтр в simple mode
- Expected: 200

## 8. Командный фильтр — попытка установить team_id
- Auth: JWT
- Body: { "team_id": "<other_team_id>" }
- Preconditions: командный фильтр в simple mode
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Командный фильтр — с правом team.manage_filters
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь имеет team.manage_filters
- Expected: 200

## 10. Командный фильтр — без права team.manage_filters
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь — участник без team.manage_filters
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Командный фильтр — не участник команды
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. Чужой личный фильтр
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 13. Невалидное RSQL выражение
- Auth: JWT
- Body: { "rsql": "invalid %%%" }
- Preconditions: фильтр в rsql mode
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Пустая строка в name
- Auth: JWT
- Body: { "name": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Невалидный filter_mode
- Auth: JWT
- Body: { "filter_mode": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid
- Body: { "name": "Updated" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 17. Без авторизации
- Auth: нет
- Body: { "name": "Updated" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 18. Пустое тело запроса
- Auth: JWT
- Body: {}
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, ... } (без изменений)
