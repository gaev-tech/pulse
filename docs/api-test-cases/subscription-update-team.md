# PATCH /api/v1/teams/{prefix}/subscription

## 1. Успешный апгрейд команды Free → Team (владелец)
- Auth: JWT
- Preconditions: пользователь — владелец команды с правом team.manage_owners; текущий план — free
- Body: { "plan": "team" }
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "team", status: "active" }

## 2. Успешный апгрейд Free → Enterprise
- Auth: JWT
- Preconditions: пользователь — владелец команды с правом team.manage_owners
- Body: { "plan": "enterprise" }
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "enterprise", status: "active" }

## 3. Успешный downgrade Team → Free
- Auth: JWT
- Preconditions: пользователь — владелец команды с правом team.manage_owners; текущий план — team
- Body: { "plan": "free" }
- Expected: 200
- Response: { subject_type: "team", subject_id: "<team_id>", plan: "free", status: "active" }

## 4. Нет права team.manage_owners
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_owners
- Body: { "plan": "team" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Невалидное значение plan
- Auth: JWT
- Preconditions: пользователь — владелец команды с правом team.manage_owners
- Body: { "plan": "ultimate" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 6. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Несуществующая команда
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/subscription
- Body: { "plan": "team" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Body: { "plan": "team" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
