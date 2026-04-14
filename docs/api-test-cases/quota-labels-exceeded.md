# POST /api/v1/labels — превышение лимита меток

## 1. Личный лимит меток исчерпан (Free: 20 личных меток)
- Auth: JWT
- Preconditions: у пользователя создано 20 личных меток; план — free
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "New", "color": "#ff0000" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 2. Командный лимит меток исчерпан (Free/Pro: 20 командных меток)
- Auth: JWT
- Preconditions: в команде создано 20 меток; у команды план free
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "New", "color": "#ff0000" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 3. Team-подписка снимает лимит командных меток
- Auth: JWT
- Preconditions: в команде 20+ меток; у команды план team
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "New", "color": "#ff0000" }
- Expected: 201
- Response: { id, name: "New", color: "#ff0000" }

## 4. Pro-подписка снимает лимит личных меток
- Auth: JWT
- Preconditions: у пользователя 20+ личных меток; план — pro
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "New", "color": "#ff0000" }
- Expected: 201
- Response: { id, name: "New", color: "#ff0000" }
