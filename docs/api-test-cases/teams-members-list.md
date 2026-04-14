# GET /api/v1/teams/{prefix}/members

## 1. Успешное получение списка участников (JWT)
- Auth: JWT
- Preconditions: пользователь — участник команды, в команде несколько участников
- Expected: 200
- Response: [{ user: { id, email, username }, permissions: number, joined_at }, ...]

## 2. Успешное получение списка участников (PAT)
- Auth: PAT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: [{ user: { id, email, username }, permissions: number, joined_at }, ...]

## 3. Проверка permissions — владелец команды в списке
- Auth: JWT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: владелец команды присутствует в списке с полными permissions

## 4. Проверка permissions — участник с ограниченными правами
- Auth: JWT
- Preconditions: участник имеет права [task.create, edit.title]
- Expected: 200
- Response: permissions участника содержит соответствующие биты

## 5. Команда с одним участником (владельцем)
- Auth: JWT
- Preconditions: в команде только владелец
- Expected: 200
- Response: [{ user: { id, email, username }, permissions: number, joined_at }]

## 6. Не участник команды
- Auth: JWT
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий prefix
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/members
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
