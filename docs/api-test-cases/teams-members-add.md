# POST /api/v1/teams/{prefix}/members

## 1. Успешное приглашение участника с правами (JWT)
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": ["edit.title", "team.manage_filters"] }
- Preconditions: пользователь — владелец команды, newmember@example.com существует и не участник
- Expected: 201 Created
- Side effects: участник добавлен с правами edit.title, team.manage_filters + автоматический view

## 2. Успешное приглашение участника с правами (PAT)
- Auth: PAT
- Body: { "email": "newmember@example.com", "permissions": ["task.create"] }
- Preconditions: пользователь — владелец команды
- Expected: 201 Created

## 3. Приглашение с пустым массивом permissions — выдаётся только view
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": [] }
- Preconditions: пользователь — владелец команды
- Expected: 201 Created
- Side effects: участник добавлен только с правом view

## 4. Приглашение без поля permissions — выдаётся только view
- Auth: JWT
- Body: { "email": "newmember@example.com" }
- Preconditions: пользователь — владелец команды
- Expected: 201 Created
- Side effects: участник добавлен только с правом view

## 5. Приглашение уже существующего участника
- Auth: JWT
- Body: { "email": "existing@example.com", "permissions": ["edit.title"] }
- Preconditions: existing@example.com уже участник команды
- Expected: 409
- Response: { error: { code: "CONFLICT", message: "..." } }

## 6. Приглашение несуществующего пользователя
- Auth: JWT
- Body: { "email": "nonexistent@example.com", "permissions": [] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 7. Невалидные значения permissions
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": ["invalid_permission"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Отсутствует поле email
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Пустая строка в email
- Auth: JWT
- Body: { "email": "", "permissions": [] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Невалидный формат email
- Auth: JWT
- Body: { "email": "not-an-email", "permissions": [] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Permissions с невалидным типом (строка вместо массива)
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": "edit.title" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Участник с team.manage_owners приглашает
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": ["edit.title"] }
- Preconditions: пользователь — участник с правом team.manage_owners
- Expected: 201 Created

## 13. Участник без team.manage_owners
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": [] }
- Preconditions: пользователь — участник без team.manage_owners
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 14. Не участник команды
- Auth: JWT
- Body: { "email": "newmember@example.com", "permissions": [] }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 15. Несуществующий prefix команды
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/members
- Body: { "email": "newmember@example.com", "permissions": [] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 16. Без авторизации
- Auth: нет
- Body: { "email": "newmember@example.com", "permissions": [] }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 17. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
