# POST /api/v1/automations

## 1. Успешное создание личной автоматизации (JWT)
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Close notification", "trigger": "task.status_changed", "actions": [{ "type": "api_call", "url": "https://hooks.example.com/notify", "method": "POST" }] }
- Expected: 201
- Response: { id, name: "Close notification", trigger: "task.status_changed", criteria: null, actions: [...], enabled: true }

## 2. Успешное создание автоматизации (PAT)
- Auth: PAT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Auto", "trigger": "task.created", "actions": [{ "type": "api_call", "url": "https://hooks.example.com", "method": "POST" }] }
- Expected: 201

## 3. Создание с критериями
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Auto", "trigger": "task.status_changed", "criteria": { "new_status": "closed" }, "actions": [{ "type": "api_call", "url": "https://hooks.example.com", "method": "POST" }] }
- Expected: 201
- Response: { id, name, trigger, criteria: { new_status: "closed" }, actions: [...], enabled: true }

## 4. Командная автоматизация — с правом team.manage_automations
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Team auto", "trigger": "task.created", "actions": [...] }
- Preconditions: пользователь имеет team.manage_automations
- Expected: 201

## 5. Командная автоматизация — без права team.manage_automations
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Preconditions: пользователь — участник без team.manage_automations
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Командная автоматизация — не участник команды
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Личная автоматизация от имени другого пользователя
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<other_user_id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Отсутствует поле name
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "trigger": "task.created", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Отсутствует поле trigger
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Auto", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Отсутствует поле actions
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Auto", "trigger": "task.created" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Отсутствует поле owner_type
- Auth: JWT
- Body: { "owner_id": "<id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Отсутствует поле owner_id
- Auth: JWT
- Body: { "owner_type": "user", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Невалидный trigger
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Auto", "trigger": "invalid.trigger", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Пустой массив actions
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Auto", "trigger": "task.created", "actions": [] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Пустая строка в name
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "", "trigger": "task.created", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Невалидный owner_type
- Auth: JWT
- Body: { "owner_type": "invalid", "owner_id": "<id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Без авторизации
- Auth: нет
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Auto", "trigger": "task.created", "actions": [...] }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 18. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
