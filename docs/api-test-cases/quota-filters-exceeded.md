# POST /api/v1/filters — превышение лимита фильтров

## 1. Личный лимит фильтров исчерпан (Free: 10 личных фильтров)
- Auth: JWT
- Preconditions: у пользователя создано 10 личных фильтров; план — free
- Body: { "name": "New filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 2. Командный лимит фильтров исчерпан (Free/Pro: 10 командных фильтров)
- Auth: JWT
- Preconditions: в команде создано 10 фильтров; у команды план free
- Body: { "name": "New filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple" }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 3. Team-подписка снимает лимит командных фильтров
- Auth: JWT
- Preconditions: в команде 10+ фильтров; у команды план team
- Body: { "name": "New filter", "owner_type": "team", "owner_id": "<team_id>", "filter_mode": "simple" }
- Expected: 201
- Response: { id, name: "New filter", owner_type: "team", ... }

## 4. Pro-подписка снимает лимит личных фильтров
- Auth: JWT
- Preconditions: у пользователя 10+ личных фильтров; план — pro
- Body: { "name": "New filter", "owner_type": "user", "owner_id": "<current_user_id>", "filter_mode": "simple" }
- Expected: 201
- Response: { id, name: "New filter", owner_type: "user", ... }
