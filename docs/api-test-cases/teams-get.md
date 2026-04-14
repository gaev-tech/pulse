# GET /api/v1/teams/{prefix}

## 1. Успешное получение команды по prefix (JWT)
- Auth: JWT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: { id, name, prefix, owner_id, team_permissions: number, created_at }

## 2. Успешное получение команды по prefix (PAT)
- Auth: PAT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: { id, name, prefix, owner_id, team_permissions: number, created_at }

## 3. Проверка team_permissions — владелец
- Auth: JWT
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: team_permissions содержит все биты

## 4. Проверка team_permissions — участник с ограниченными правами
- Auth: JWT
- Preconditions: пользователь — участник с правами [task.create, team.manage_filters]
- Expected: 200
- Response: team_permissions содержит биты 0 и 4

## 5. Несуществующий prefix
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 6. Не участник команды
- Auth: JWT
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
