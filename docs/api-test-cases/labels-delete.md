# DELETE /api/v1/labels/{labelID}

## 1. Успешное удаление личной метки (JWT)
- Auth: JWT
- Preconditions: метка принадлежит текущему пользователю
- Expected: 204 No Content
- Side effects: метка убирается из всех задач, в которых она использовалась

## 2. Успешное удаление метки (PAT)
- Auth: PAT
- Preconditions: метка принадлежит текущему пользователю
- Expected: 204 No Content

## 3. Удаление командной метки — с правом team.manage_labels
- Auth: JWT
- Preconditions: метка принадлежит команде, пользователь имеет team.manage_labels
- Expected: 204 No Content
- Side effects: метка убирается из всех командных задач

## 4. Удаление командной метки — без права team.manage_labels
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_labels
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Удаление чужой личной метки
- Auth: JWT
- Preconditions: метка принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Командная метка — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий labelID
- Auth: JWT
- URL: /api/v1/labels/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
