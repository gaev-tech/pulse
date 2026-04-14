# GET /api/v1/teams/{prefix}/members/search?email=

## 1. Успешный поиск пользователя для приглашения (JWT)
- Auth: JWT
- Query: ?email=user@example.com
- Preconditions: пользователь с таким email существует и не является участником команды
- Expected: 200
- Response: [{ id, email, username }]

## 2. Успешный поиск пользователя для приглашения (PAT)
- Auth: PAT
- Query: ?email=user@example.com
- Preconditions: пользователь существует и не участник команды
- Expected: 200
- Response: [{ id, email, username }]

## 3. Поиск — пользователь уже является участником команды
- Auth: JWT
- Query: ?email=member@example.com
- Preconditions: пользователь с этим email уже участник команды
- Expected: 200
- Response: [] (пустой массив — участники исключаются)

## 4. Поиск — пользователь не найден
- Auth: JWT
- Query: ?email=nonexistent@example.com
- Preconditions: пользователя с таким email нет
- Expected: 200
- Response: []

## 5. Частичный поиск по email
- Auth: JWT
- Query: ?email=user
- Preconditions: есть пользователи с email, содержащим "user"
- Expected: 200
- Response: [{ id, email, username }, ...]

## 6. Без параметра email
- Auth: JWT
- Query: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Пустая строка в email
- Auth: JWT
- Query: ?email=
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Query: ?email=user@example.com
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 9. Несуществующий prefix команды
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/members/search?email=user@example.com
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Не участник команды
- Auth: JWT
- Query: ?email=user@example.com
- Preconditions: текущий пользователь не участник команды {prefix}
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }
