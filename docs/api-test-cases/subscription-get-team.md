# GET /api/v1/teams/{prefix}/subscription

## 1. Успешное получение подписки команды участником (Free по умолчанию)
- Auth: JWT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "free", status: "active" }

## 2. Успешное получение подписки команды с планом Team
- Auth: JWT
- Preconditions: команда имеет plan: "team"
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "team", status: "active" }

## 3. Успешное получение через PAT
- Auth: PAT
- Preconditions: пользователь — участник команды
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "free", status: "active" }

## 4. Пользователь не участник команды
- Auth: JWT
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Несуществующая команда
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/subscription
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 6. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
