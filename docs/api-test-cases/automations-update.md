# PATCH /api/v1/automations/{automationID}

## 1. Успешное изменение name (JWT)
- Auth: JWT
- Body: { "name": "Updated Automation" }
- Preconditions: автоматизация принадлежит текущему пользователю
- Expected: 200
- Response: { id, name: "Updated Automation", trigger, criteria, actions, enabled }

## 2. Успешное изменение name (PAT)
- Auth: PAT
- Body: { "name": "Updated" }
- Expected: 200

## 3. Изменение trigger
- Auth: JWT
- Body: { "trigger": "task.title_changed" }
- Expected: 200
- Response: { trigger: "task.title_changed", ... }

## 4. Изменение criteria
- Auth: JWT
- Body: { "criteria": { "new_status": "opened" } }
- Expected: 200
- Response: { criteria: { new_status: "opened" }, ... }

## 5. Изменение actions
- Auth: JWT
- Body: { "actions": [{ "type": "api_call", "url": "https://new-hook.example.com", "method": "POST" }] }
- Expected: 200

## 6. Включение автоматизации
- Auth: JWT
- Body: { "enabled": true }
- Preconditions: автоматизация выключена
- Expected: 200
- Response: { enabled: true, ... }

## 7. Выключение автоматизации
- Auth: JWT
- Body: { "enabled": false }
- Preconditions: автоматизация включена
- Expected: 200
- Response: { enabled: false, ... }

## 8. Командная автоматизация — с правом team.manage_automations
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь имеет team.manage_automations
- Expected: 200

## 9. Командная автоматизация — без права team.manage_automations
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь — участник без team.manage_automations
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 10. Командная автоматизация — не участник команды
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Чужая личная автоматизация
- Auth: JWT
- Body: { "name": "Updated" }
- Preconditions: автоматизация принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. Невалидный trigger
- Auth: JWT
- Body: { "trigger": "invalid.trigger" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Пустой массив actions
- Auth: JWT
- Body: { "actions": [] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Пустая строка в name
- Auth: JWT
- Body: { "name": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Несуществующий automationID
- Auth: JWT
- URL: /api/v1/automations/nonexistent-uuid
- Body: { "name": "Updated" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 16. Без авторизации
- Auth: нет
- Body: { "name": "Updated" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 17. Пустое тело запроса
- Auth: JWT
- Body: {}
- Preconditions: автоматизация принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, trigger, criteria, actions, enabled } (без изменений)
