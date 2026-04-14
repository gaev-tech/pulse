# POST /api/v1/tasks

## 1. Успешное создание личной задачи (JWT)
- Auth: JWT
- Body: { "title": "My Task", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 201
- Response: { id, key: "#<number>", title: "My Task", description: null, status: "opened", owner_type: "user", owner_id, assignee_id: null, parent_id: null, label_ids: [], task_permissions: number, created_at, updated_at }

## 2. Успешное создание личной задачи (PAT)
- Auth: PAT
- Body: { "title": "My Task", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 201
- Response: { ...task, task_permissions: number }

## 3. Успешное создание командной задачи
- Auth: JWT
- Body: { "title": "Team Task", "owner_type": "team", "owner_id": "<team_id>" }
- Preconditions: пользователь имеет право task.create в команде
- Expected: 201
- Response: { id, key: "PREFIX-<number>", title: "Team Task", owner_type: "team", owner_id: "<team_id>", task_permissions: number, ... }

## 4. Создание задачи со всеми опциональными полями
- Auth: JWT
- Body: { "title": "Full Task", "description": "## Description", "owner_type": "user", "owner_id": "<current_user_id>", "assignee_id": "<user_id>", "parent_id": "<parent_task_id>", "label_ids": ["<label_id_1>", "<label_id_2>"] }
- Expected: 201
- Response: { ...task, description: "## Description", assignee_id, parent_id, label_ids: [...] }

## 5. Создание командной задачи без права task.create
- Auth: JWT
- Body: { "title": "Team Task", "owner_type": "team", "owner_id": "<team_id>" }
- Preconditions: пользователь — участник команды без права task.create
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Создание командной задачи — не участник команды
- Auth: JWT
- Body: { "title": "Team Task", "owner_type": "team", "owner_id": "<team_id>" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Создание личной задачи от имени другого пользователя
- Auth: JWT
- Body: { "title": "Task", "owner_type": "user", "owner_id": "<other_user_id>" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Отсутствует поле title
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Пустая строка в title
- Auth: JWT
- Body: { "title": "", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Отсутствует поле owner_type
- Auth: JWT
- Body: { "title": "Task", "owner_id": "<current_user_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Отсутствует поле owner_id
- Auth: JWT
- Body: { "title": "Task", "owner_type": "user" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Невалидный owner_type
- Auth: JWT
- Body: { "title": "Task", "owner_type": "invalid", "owner_id": "<id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Несуществующий owner_id (team)
- Auth: JWT
- Body: { "title": "Task", "owner_type": "team", "owner_id": "nonexistent-uuid" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 14. Несуществующий parent_id
- Auth: JWT
- Body: { "title": "Task", "owner_type": "user", "owner_id": "<current_user_id>", "parent_id": "nonexistent-uuid" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 15. label_ids с чужими метками (личная метка для командной задачи)
- Auth: JWT
- Body: { "title": "Task", "owner_type": "team", "owner_id": "<team_id>", "label_ids": ["<personal_label_id>"] }
- Preconditions: label принадлежит другому пользователю, а не команде
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. label_ids с несуществующим ID
- Auth: JWT
- Body: { "title": "Task", "owner_type": "user", "owner_id": "<current_user_id>", "label_ids": ["nonexistent-uuid"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Title с невалидным типом (число)
- Auth: JWT
- Body: { "title": 123, "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 18. Без авторизации
- Auth: нет
- Body: { "title": "Task", "owner_type": "user", "owner_id": "<id>" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 19. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 20. Атомарная генерация key_number
- Auth: JWT
- Preconditions: создать несколько задач параллельно для одного owner
- Expected: 201 для каждой, все key_number уникальны и последовательны
