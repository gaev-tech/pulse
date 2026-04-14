# DELETE /api/v1/tasks/{key}/attachments/{attachmentID}

## 1. Успешное удаление вложения (JWT)
- Auth: JWT
- Preconditions: пользователь имеет право edit.attachments, вложение существует
- Expected: 204 No Content
- Side effects: файл удалён из хранилища

## 2. Успешное удаление вложения (PAT)
- Auth: PAT
- Preconditions: пользователь имеет право edit.attachments
- Expected: 204 No Content

## 3. Владелец задачи удаляет вложение
- Auth: JWT
- Preconditions: пользователь — владелец личной задачи
- Expected: 204 No Content

## 4. Участник команды с edit.attachments
- Auth: JWT
- Preconditions: пользователь — участник команды с правом edit.attachments
- Expected: 204 No Content

## 5. Пользователь с direct share и edit.attachments
- Auth: JWT
- Preconditions: задача расшарена с правом edit.attachments
- Expected: 204 No Content

## 6. Нет права edit.attachments
- Auth: JWT
- Preconditions: пользователь имеет view, но не edit.attachments
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Нет доступа к задаче
- Auth: JWT
- Preconditions: пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Несуществующий attachmentID
- Auth: JWT
- URL: /api/v1/tasks/{key}/attachments/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 9. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999/attachments/{attachmentID}
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
