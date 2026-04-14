# POST /api/v1/teams/{prefix}/members — превышение лимита участников

## 1. Лимит участников исчерпан (Free/Pro: 30 участников)
- Auth: JWT
- Preconditions: в команде 30 участников (включая владельца); у команды план free
- Body: { "email": "newuser@example.com", "permissions": [] }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 2. Лимит участников исчерпан (Team: 60 участников)
- Auth: JWT
- Preconditions: в команде 60 участников; у команды план team
- Body: { "email": "newuser@example.com", "permissions": [] }
- Expected: 403
- Response: { error: { code: "QUOTA_EXCEEDED", message: "..." } }

## 3. Enterprise позволяет добавлять участников свыше 60
- Auth: JWT
- Preconditions: в команде 60+ участников; у команды план enterprise
- Body: { "email": "newuser@example.com", "permissions": [] }
- Expected: 201
- Note: при превышении 60 участников выставляется billing event (доплата), но добавление не блокируется

## 4. Team-план — 59 участников (лимит не исчерпан)
- Auth: JWT
- Preconditions: в команде 59 участников; у команды план team
- Body: { "email": "newuser@example.com", "permissions": [] }
- Expected: 201
- Response: 201 Created
