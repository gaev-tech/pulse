# GET /api/v1/subscription

## 1. Успешное получение подписки нового пользователя (Free по умолчанию)
- Auth: JWT
- Expected: 200
- Response: { subject_type: "user", subject_id: "<current_user_id>", plan: "free", status: "active" }

## 2. Успешное получение подписки после апгрейда до Pro
- Auth: JWT
- Preconditions: подписка пользователя изменена на plan: "pro"
- Expected: 200
- Response: { subject_type: "user", subject_id: "<current_user_id>", plan: "pro", status: "active" }

## 3. Успешное получение через PAT
- Auth: PAT
- Expected: 200
- Response: { subject_type: "user", subject_id: "<current_user_id>", plan: "free", status: "active" }

## 4. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
