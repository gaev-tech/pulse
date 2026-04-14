# POST /api/v1/labels

## 1. Успешное создание личной метки (JWT)
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Bug", "color": "#ff0000" }
- Expected: 201
- Response: { id, name: "Bug", color: "#ff0000" }

## 2. Успешное создание личной метки (PAT)
- Auth: PAT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Feature", "color": "#00ff00" }
- Expected: 201
- Response: { id, name: "Feature", color: "#00ff00" }

## 3. Успешное создание командной метки
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Priority", "color": "#0000ff" }
- Preconditions: пользователь имеет право team.manage_labels в команде
- Expected: 201
- Response: { id, name: "Priority", color: "#0000ff" }

## 4. Командная метка — нет права team.manage_labels
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Label", "color": "#000000" }
- Preconditions: пользователь — участник без team.manage_labels
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Командная метка — не участник команды
- Auth: JWT
- Body: { "owner_type": "team", "owner_id": "<team_id>", "name": "Label", "color": "#000000" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Личная метка от имени другого пользователя
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<other_user_id>", "name": "Label", "color": "#000000" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Отсутствует поле name
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "color": "#ff0000" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Отсутствует поле color
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Bug" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Отсутствует поле owner_type
- Auth: JWT
- Body: { "owner_id": "<current_user_id>", "name": "Bug", "color": "#ff0000" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Отсутствует поле owner_id
- Auth: JWT
- Body: { "owner_type": "user", "name": "Bug", "color": "#ff0000" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Пустая строка в name
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "", "color": "#ff0000" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Пустая строка в color
- Auth: JWT
- Body: { "owner_type": "user", "owner_id": "<current_user_id>", "name": "Bug", "color": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Невалидный owner_type
- Auth: JWT
- Body: { "owner_type": "invalid", "owner_id": "<id>", "name": "Bug", "color": "#ff0000" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Без авторизации
- Auth: нет
- Body: { "owner_type": "user", "owner_id": "<id>", "name": "Bug", "color": "#ff0000" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 15. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
