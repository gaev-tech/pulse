# GET /api/v1/tasks/{key}/share/search?email=

## 1. Успешный поиск пользователя для шаринга (JWT)
- Auth: JWT
- Query: ?email=user@example.com
- Preconditions: пользователь существует и не имеет доступа к задаче
- Expected: 200
- Response: [{ id, email, username }]

## 2. Успешный поиск пользователя для шаринга (PAT)
- Auth: PAT
- Query: ?email=user@example.com
- Preconditions: пользователь существует и не имеет доступа к задаче
- Expected: 200
- Response: [{ id, email, username }]

## 3. Поиск — у пользователя уже есть доступ к задаче
- Auth: JWT
- Query: ?email=shared@example.com
- Preconditions: пользователь уже имеет доступ к задаче
- Expected: 200
- Response: [] (пустой массив — уже имеющие доступ исключаются)

## 4. Поиск — пользователь не найден
- Auth: JWT
- Query: ?email=nonexistent@example.com
- Expected: 200
- Response: []

## 5. Частичный поиск по email
- Auth: JWT
- Query: ?email=user
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

## 9. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999/share/search?email=user@example.com
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Нет права share на задачу
- Auth: JWT
- Query: ?email=user@example.com
- Preconditions: текущий пользователь имеет доступ к задаче, но без права share
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Нет доступа к задаче
- Auth: JWT
- Query: ?email=user@example.com
- Preconditions: текущий пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }
