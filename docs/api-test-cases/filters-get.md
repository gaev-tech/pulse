# GET /api/v1/filters/{filterID}

## 1. Успешное получение личного фильтра (JWT)
- Auth: JWT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, owner_type: "user", owner_id, filter_mode, ...criteria fields }

## 2. Успешное получение личного фильтра (PAT)
- Auth: PAT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, owner_type, owner_id, filter_mode, ... }

## 3. Получение командного фильтра — участник команды
- Auth: JWT
- Preconditions: фильтр принадлежит команде, пользователь — участник
- Expected: 200
- Response: { id, name, owner_type: "team", owner_id, filter_mode, ... }

## 4. Simple mode фильтр — все критерии присутствуют
- Auth: JWT
- Preconditions: фильтр в simple mode со всеми критериями
- Expected: 200
- Response: { filter_mode: "simple", search_contains, assignee_ids, status, label_ids, team_id }

## 5. RSQL mode фильтр
- Auth: JWT
- Preconditions: фильтр в rsql mode
- Expected: 200
- Response: { filter_mode: "rsql", rsql: "..." }

## 6. Командный фильтр — не участник команды
- Auth: JWT
- Preconditions: фильтр принадлежит команде, пользователь не участник
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Чужой личный фильтр
- Auth: JWT
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 9. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
