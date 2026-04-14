# POST /api/v1/filters

## 1. Успешное создание личного фильтра в simple mode (JWT)
- Auth: JWT
- Body: { "name": "My Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple", "status": "opened", "assignee_ids": ["<user_id>"] }
- Expected: 201
- Response: { id, name: "My Filter", owner_type: "user", owner_id, filter_mode: "simple", status: "opened", assignee_ids: [...], ... }

## 2. Успешное создание личного фильтра (PAT)
- Auth: PAT
- Body: { "name": "My Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple" }
- Expected: 201
- Response: { id, name, owner_type, owner_id, filter_mode, ... }

## 3. Личный фильтр в simple mode с team_id
- Auth: JWT
- Body: { "name": "Team Tasks Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple", "team_id": "<team_id>" }
- Expected: 201
- Response: { id, name, filter_mode: "simple", team_id: "<team_id>", ... }

## 4. Личный фильтр в simple mode со всеми критериями
- Auth: JWT
- Body: { "name": "Full Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple", "search_contains": "bug", "assignee_ids": ["<user_id>"], "status": "opened", "label_ids": ["<label_id>"], "team_id": "<team_id>" }
- Expected: 201

## 5. Личный фильтр в rsql mode
- Auth: JWT
- Body: { "name": "RSQL Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "rsql", "rsql": "status==opened;assignee_id==<user_id>" }
- Expected: 201
- Response: { id, name, filter_mode: "rsql", rsql: "status==opened;assignee_id==<user_id>", ... }

## 6. Командный фильтр в simple mode
- Auth: JWT
- Body: { "name": "Team Filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple", "status": "opened" }
- Preconditions: пользователь имеет право team.manage_filters
- Expected: 201

## 7. Командный фильтр в rsql mode
- Auth: JWT
- Body: { "name": "Team RSQL", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "rsql", "rsql": "status==closed" }
- Preconditions: пользователь имеет право team.manage_filters
- Expected: 201

## 8. Командный фильтр с team_id — ошибка (team_id только для личных)
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple", "team_id": "<other_team_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Командный фильтр — нет права team.manage_filters
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple" }
- Preconditions: пользователь — участник без team.manage_filters
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 10. Командный фильтр — не участник команды
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Личный фильтр от имени другого пользователя
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "user", "owner_id": "<other_user_id>", "filter_mode": "simple" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. RSQL mode с невалидным выражением
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "rsql", "rsql": "invalid expression %%%" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Отсутствует поле name
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Отсутствует поле owner_type
- Auth: JWT
- Body: { "name": "Filter", "owner_id": "<id>", "filter_mode": "simple" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Отсутствует поле owner_id
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "user", "filter_mode": "simple" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Отсутствует поле filter_mode
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Невалидный filter_mode
- Auth: JWT
- Body: { "name": "Filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 18. Пустая строка в name
- Auth: JWT
- Body: { "name": "", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 19. Без авторизации
- Auth: нет
- Body: { "name": "Filter", "owner_type": "user", "owner_id": "<id>", "filter_mode": "simple" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 20. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
