# DELETE /api/v1/tasks/{key}

## 1. Успешное удаление личной задачи (JWT)
- Auth: JWT
- Preconditions: пользователь — владелец личной задачи
- Expected: 204 No Content
- Side effects: задача удалена, связанные permissions, attachments удалены

## 2. Успешное удаление задачи (PAT)
- Auth: PAT
- Preconditions: пользователь — владелец личной задачи
- Expected: 204 No Content

## 3. Удаление командной задачи — владелец команды
- Auth: JWT
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content

## 4. Удаление командной задачи — участник с полными правами
- Auth: JWT
- Preconditions: пользователь — участник команды со всеми правами
- Expected: 204 No Content

## 5. Нет доступа к задаче
- Auth: JWT
- Preconditions: пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Пользователь с view-only доступом
- Auth: JWT
- Preconditions: пользователь имеет только view на задачу (через share)
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 9. Удаление задачи с дочерними задачами
- Auth: JWT
- Preconditions: у задачи есть дочерние задачи (parent_id ссылается на неё)
- Expected: 204 No Content
- Side effects: у дочерних задач parent_id сбрасывается в null

## 10. Удаление задачи, которая блокирует другие
- Auth: JWT
- Preconditions: задача находится в blocking других задач
- Expected: 204 No Content
- Side effects: задача убирается из blocking списков других задач

## 11. Удаление задачи с relations
- Auth: JWT
- Preconditions: задача имеет связи с другими задачами
- Expected: 204 No Content
- Side effects: связи удалены у обеих сторон
