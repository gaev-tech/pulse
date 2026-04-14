# DELETE /api/v1/automations/{automationID}

## 1. Успешное удаление личной автоматизации (JWT)
- Auth: JWT
- Preconditions: автоматизация принадлежит текущему пользователю
- Expected: 204 No Content

## 2. Успешное удаление автоматизации (PAT)
- Auth: PAT
- Preconditions: автоматизация принадлежит текущему пользователю
- Expected: 204 No Content

## 3. Командная автоматизация — с правом team.manage_automations
- Auth: JWT
- Preconditions: пользователь имеет team.manage_automations
- Expected: 204 No Content

## 4. Командная автоматизация — без права team.manage_automations
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_automations
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Командная автоматизация — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Чужая личная автоматизация
- Auth: JWT
- Preconditions: автоматизация принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий automationID
- Auth: JWT
- URL: /api/v1/automations/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
