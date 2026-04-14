# POST /api/v1/teams

## 1. Успешное создание команды (JWT)
- Auth: JWT
- Body: { "name": "Acme Corp", "prefix": "ACME" }
- Expected: 201
- Response: { id, name: "Acme Corp", prefix: "ACME", owner_id: <current_user_id>, created_at }
- Side effects: текущий пользователь становится владельцем команды со всеми правами

## 2. Успешное создание команды (PAT)
- Auth: PAT
- Body: { "name": "Beta Team", "prefix": "BETA" }
- Expected: 201
- Response: { id, name: "Beta Team", prefix: "BETA", owner_id: <current_user_id>, created_at }

## 3. Без авторизации
- Auth: нет
- Body: { "name": "Acme Corp", "prefix": "ACME" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 4. Невалидный токен
- Auth: JWT (invalid)
- Body: { "name": "Acme Corp", "prefix": "ACME" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 5. Дублирующийся prefix
- Auth: JWT
- Body: { "name": "Another Team", "prefix": "ACME" }
- Preconditions: команда с prefix "ACME" уже существует
- Expected: 409
- Response: { error: { code: "CONFLICT", message: "..." } }

## 6. Отсутствует поле name
- Auth: JWT
- Body: { "prefix": "ACME" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Отсутствует поле prefix
- Auth: JWT
- Body: { "name": "Acme Corp" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Пустая строка в name
- Auth: JWT
- Body: { "name": "", "prefix": "ACME" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Пустая строка в prefix
- Auth: JWT
- Body: { "name": "Acme Corp", "prefix": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Name с невалидным типом (число)
- Auth: JWT
- Body: { "name": 123, "prefix": "ACME" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Prefix с невалидным типом (число)
- Auth: JWT
- Body: { "name": "Acme Corp", "prefix": 123 }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
