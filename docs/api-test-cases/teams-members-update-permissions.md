# PATCH /api/v1/teams/{prefix}/members/{userID}/permissions

## 1. Успешное изменение прав участника (JWT)
- Auth: JWT
- Body: { "permissions": ["edit.title", "edit.description", "task.create"] }
- Preconditions: пользователь — владелец команды, целевой участник существует
- Expected: 204 No Content
- Side effects: права участника заменены на указанные + автоматический view

## 2. Успешное изменение прав участника (PAT)
- Auth: PAT
- Body: { "permissions": ["edit.title"] }
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content

## 3. Пустой массив permissions — остаётся только view
- Auth: JWT
- Body: { "permissions": [] }
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content
- Side effects: у участника остаётся только право view

## 4. Полный набор прав
- Auth: JWT
- Body: { "permissions": ["task.create", "team.manage_owners", "team.edit_labels", "team.manage_labels", "team.manage_filters", "team.manage_automations", "team.import", "edit.title", "edit.description", "edit.status", "edit.assignee", "edit.labels", "edit.links", "edit.relations", "edit.blocking", "edit.parent", "edit.attachments", "share"] }
- Expected: 204 No Content

## 5. Невалидные значения permissions
- Auth: JWT
- Body: { "permissions": ["invalid_permission"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 6. Permissions с невалидным типом (строка вместо массива)
- Auth: JWT
- Body: { "permissions": "edit.title" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Отсутствует поле permissions
- Auth: JWT
- Body: {}
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Изменение прав владельца команды
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Preconditions: userID — владелец команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 9. Участник с team.manage_owners меняет права
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Preconditions: пользователь — участник с team.manage_owners
- Expected: 204 No Content

## 10. Участник без team.manage_owners
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Preconditions: пользователь — участник без team.manage_owners
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Не участник команды
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. Несуществующий userID
- Auth: JWT
- URL: /api/v1/teams/ACME/members/nonexistent-uuid/permissions
- Body: { "permissions": ["edit.title"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 13. Несуществующий prefix команды
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/members/{userID}/permissions
- Body: { "permissions": ["edit.title"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 14. Без авторизации
- Auth: нет
- Body: { "permissions": ["edit.title"] }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 15. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
