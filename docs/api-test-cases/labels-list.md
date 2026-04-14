# GET /api/v1/labels?owner_type=&owner_id=

## 1. Успешное получение личных меток (JWT)
- Auth: JWT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Preconditions: у пользователя есть несколько меток
- Expected: 200
- Response: [{ id, name, color }, ...]

## 2. Успешное получение личных меток (PAT)
- Auth: PAT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Expected: 200
- Response: [{ id, name, color }, ...]

## 3. Успешное получение командных меток
- Auth: JWT
- Query: ?owner_type=team&owner_id=<team_id>
- Preconditions: пользователь — участник команды, у команды есть метки
- Expected: 200
- Response: [{ id, name, color }, ...]

## 4. Нет меток — пустой массив
- Auth: JWT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Preconditions: у пользователя нет меток
- Expected: 200
- Response: []

## 5. Метки чужого пользователя
- Auth: JWT
- Query: ?owner_type=user&owner_id=<other_user_id>
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Командные метки — не участник команды
- Auth: JWT
- Query: ?owner_type=team&owner_id=<team_id>
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Отсутствует параметр owner_type
- Auth: JWT
- Query: ?owner_id=<id>
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Отсутствует параметр owner_id
- Auth: JWT
- Query: ?owner_type=user
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Невалидный owner_type
- Auth: JWT
- Query: ?owner_type=invalid&owner_id=<id>
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Несуществующий owner_id
- Auth: JWT
- Query: ?owner_type=team&owner_id=nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 11. Без авторизации
- Auth: нет
- Query: ?owner_type=user&owner_id=<id>
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
