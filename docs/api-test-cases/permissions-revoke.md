# DELETE /api/v1/tasks/{key}/permissions/{permissionID}

## 1. Успешный отзыв доступа (JWT)
- Auth: JWT
- Preconditions: текущий пользователь имеет право share, permissionID существует
- Expected: 204 No Content
- Side effects: доступ целевого пользователя к задаче полностью отозван (включая view)

## 2. Успешный отзыв доступа (PAT)
- Auth: PAT
- Preconditions: текущий пользователь имеет право share
- Expected: 204 No Content

## 3. Каскадный отзыв цепочки прав (A→B→C)
- Auth: JWT
- Preconditions: A выдал доступ B, B выдал доступ C; отзываем доступ у B
- Expected: 204 No Content
- Side effects: доступ у B и C каскадно удалён

## 4. Нет права share
- Auth: JWT
- Preconditions: текущий пользователь имеет доступ к задаче, но без права share
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Нет доступа к задаче
- Auth: JWT
- Preconditions: текущий пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Несуществующий permissionID
- Auth: JWT
- URL: /api/v1/tasks/{key}/permissions/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 7. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999/permissions/{permissionID}
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 9. Владелец задачи отзывает доступ
- Auth: JWT
- Preconditions: текущий пользователь — владелец задачи (personal task)
- Expected: 204 No Content
