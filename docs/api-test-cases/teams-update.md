# PATCH /api/v1/teams/{prefix}

## 1. Успешное изменение названия команды (JWT)
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: { id, name: "New Name", prefix, owner_id, created_at }

## 2. Успешное изменение названия команды (PAT)
- Auth: PAT
- Body: { "name": "New Name" }
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: { id, name: "New Name", prefix, owner_id, created_at }

## 3. Успешное изменение prefix команды
- Auth: JWT
- Body: { "prefix": "NEWPFX" }
- Preconditions: пользователь — владелец команды, prefix "NEWPFX" свободен
- Expected: 200
- Response: { id, name, prefix: "NEWPFX", owner_id, created_at }

## 4. Изменение name и prefix одновременно
- Auth: JWT
- Body: { "name": "New Name", "prefix": "NEWPFX" }
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: { id, name: "New Name", prefix: "NEWPFX", owner_id, created_at }

## 5. Изменение prefix на уже занятый
- Auth: JWT
- Body: { "prefix": "TAKEN" }
- Preconditions: команда с prefix "TAKEN" уже существует
- Expected: 409
- Response: { error: { code: "CONFLICT", message: "..." } }

## 6. Пустая строка в name
- Auth: JWT
- Body: { "name": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Пустая строка в prefix
- Auth: JWT
- Body: { "prefix": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Name с невалидным типом
- Auth: JWT
- Body: { "name": 123 }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Prefix с невалидным типом
- Auth: JWT
- Body: { "prefix": 123 }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Не владелец команды — участник с team.manage_owners
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь — участник с правом team.manage_owners
- Expected: 200
- Response: { id, name: "New Name", prefix, owner_id, created_at }

## 11. Участник без нужных прав
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь — участник без team.manage_owners
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. Не участник команды
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 13. Несуществующий prefix команды
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT
- Body: { "name": "New Name" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 14. Без авторизации
- Auth: нет
- Body: { "name": "New Name" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 15. Пустое тело запроса
- Auth: JWT
- Body: {}
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: { id, name, prefix, owner_id, created_at } (без изменений)
