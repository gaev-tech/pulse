# GET /api/v1/automations?owner_type=&owner_id=

## 1. Успешное получение личных автоматизаций (JWT)
- Auth: JWT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Preconditions: у пользователя есть автоматизации
- Expected: 200
- Response: [{ id, name, trigger, enabled }, ...]

## 2. Успешное получение автоматизаций (PAT)
- Auth: PAT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Expected: 200
- Response: [{ id, name, trigger, enabled }, ...]

## 3. Командные автоматизации — участник команды
- Auth: JWT
- Query: ?owner_type=team&owner_id=<team_id>
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: [{ id, name, trigger, enabled }, ...]

## 4. Нет автоматизаций
- Auth: JWT
- Query: ?owner_type=user&owner_id=<current_user_id>
- Preconditions: у пользователя нет автоматизаций
- Expected: 200
- Response: []

## 5. Автоматизации чужого пользователя
- Auth: JWT
- Query: ?owner_type=user&owner_id=<other_user_id>
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Командные автоматизации — не участник команды
- Auth: JWT
- Query: ?owner_type=team&owner_id=<team_id>
- Preconditions: пользователь не участник команды
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

## 10. Без авторизации
- Auth: нет
- Query: ?owner_type=user&owner_id=<id>
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
