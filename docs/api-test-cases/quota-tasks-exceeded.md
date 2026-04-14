# POST /api/v1/tasks — превышение лимита задач

## 1. Личный лимит задач исчерпан (Free: 200 личных задач)
- Auth: JWT
- Preconditions: у пользователя создано 200 личных задач; план — free
- Body: { "title": "New task", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 2. Командный лимит задач исчерпан (Free/Pro: 200 командных задач)
- Auth: JWT
- Preconditions: в команде создано 200 задач; у команды план free; у владельца команды план free/pro
- Body: { "title": "New task", "owner_type": "team", "owner_id": "<team_id>" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 3. Team-подписка команды снимает лимит командных задач
- Auth: JWT
- Preconditions: в команде создано 200+ задач; у команды план team
- Body: { "title": "New task", "owner_type": "team", "owner_id": "<team_id>" }
- Expected: 201
- Response: { ...task, task_permissions: <number> }

## 4. Pro-подписка снимает лимит личных задач
- Auth: JWT
- Preconditions: у пользователя 200+ личных задач; план пользователя — pro
- Body: { "title": "New task", "owner_type": "user", "owner_id": "<current_user_id>" }
- Expected: 201
- Response: { ...task, task_permissions: <number> }
