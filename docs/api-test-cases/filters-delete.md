# DELETE /api/v1/filters/{filterID}

## 1. Успешное удаление личного фильтра (JWT)
- Auth: JWT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 204 No Content

## 2. Успешное удаление фильтра (PAT)
- Auth: PAT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 204 No Content

## 3. Удаление командного фильтра — с правом team.manage_filters
- Auth: JWT
- Preconditions: пользователь имеет team.manage_filters
- Expected: 204 No Content

## 4. Удаление командного фильтра — без права team.manage_filters
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_filters
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Командный фильтр — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Чужой личный фильтр
- Auth: JWT
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
