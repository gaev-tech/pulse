# POST /api/v1/automations — превышение лимита автоматизаций

## 1. Личный лимит автоматизаций исчерпан (Free: 1 личная автоматизация)
- Auth: JWT
- Preconditions: у пользователя создана 1 личная автоматизация; план — free
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "New", "trigger": "task.status_changed", "actions": [] }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 2. Командный лимит автоматизаций исчерпан (Free/Pro: 1 командная)
- Auth: JWT
- Preconditions: в команде создана 1 командная автоматизация; у команды план free
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "New", "trigger": "task.status_changed", "actions": [] }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 3. Pro-подписка расширяет лимит личных автоматизаций до 20
- Auth: JWT
- Preconditions: у пользователя создано 20 личных автоматизаций; план — pro
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "New", "trigger": "task.status_changed", "actions": [] }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 4. Pro допускает до 19 личных автоматизаций (лимит не исчерпан)
- Auth: JWT
- Preconditions: у пользователя 19 личных автоматизаций; план — pro
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "New", "trigger": "task.status_changed", "actions": [] }
- Expected: 201
- Response: { id, name: "New", trigger: "task.status_changed", enabled: true, ... }

## 5. Team-подписка снимает лимит командных автоматизаций
- Auth: JWT
- Preconditions: в команде 10+ автоматизаций; у команды план team
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "New", "trigger": "task.status_changed", "actions": [] }
- Expected: 201
- Response: { id, name: "New", trigger: "task.status_changed", enabled: true, ... }
