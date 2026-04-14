# DELETE /api/v1/teams/{prefix}

## 1. Успешное удаление команды (JWT)
- Auth: JWT
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content
- Side effects: каскадное удаление задач, фильтров, меток, автоматизаций, прав участников

## 2. Успешное удаление команды (PAT)
- Auth: PAT
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content
- Side effects: каскадное удаление всех ресурсов команды

## 3. Участник с team.manage_owners
- Auth: JWT
- Preconditions: пользователь — участник с правом team.manage_owners
- Expected: 204 No Content

## 4. Участник без нужных прав
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_owners
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Не участник команды
- Auth: JWT
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Несуществующий prefix
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 7. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
