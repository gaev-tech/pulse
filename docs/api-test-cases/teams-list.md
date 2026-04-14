# GET /api/v1/teams

## 1. Успешное получение списка команд (JWT)
- Auth: JWT
- Preconditions: пользователь состоит в нескольких командах
- Expected: 200
- Response: [{ id, name, prefix, owner_id, team_permissions: number, created_at }, ...]

## 2. Успешное получение списка команд (PAT)
- Auth: PAT
- Preconditions: пользователь состоит в нескольких командах
- Expected: 200
- Response: [{ id, name, prefix, owner_id, team_permissions: number, created_at }, ...]

## 3. Пользователь не состоит ни в одной команде
- Auth: JWT
- Preconditions: пользователь не участник ни одной команды
- Expected: 200
- Response: []

## 4. Проверка team_permissions — владелец команды
- Auth: JWT
- Preconditions: пользователь является владельцем команды
- Expected: 200
- Response: team_permissions содержит все биты (все права)

## 5. Проверка team_permissions — участник с ограниченными правами
- Auth: JWT
- Preconditions: пользователь — участник с правами [task.create]
- Expected: 200
- Response: team_permissions содержит только бит task.create (бит 0)

## 6. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 7. Невалидный токен
- Auth: JWT (invalid)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
