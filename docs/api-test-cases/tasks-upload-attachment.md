# POST /api/v1/tasks/{key}/attachments

## 1. Успешная загрузка вложения (JWT)
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь имеет право edit.attachments
- Expected: 200
- Response: { id, name: "<filename>", url, size, created_at }

## 2. Успешная загрузка вложения (PAT)
- Auth: PAT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь имеет право edit.attachments
- Expected: 200
- Response: { id, name, url, size, created_at }

## 3. Владелец задачи загружает вложение
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь — владелец личной задачи
- Expected: 200
- Response: { id, name, url, size, created_at }

## 4. Участник команды с edit.attachments
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь — участник команды с правом edit.attachments
- Expected: 200

## 5. Пользователь с direct share и edit.attachments
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: задача расшарена с правом edit.attachments
- Expected: 200

## 6. Нет права edit.attachments
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь имеет view, но не edit.attachments
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Нет доступа к задаче
- Auth: JWT
- Body: multipart/form-data { file: <binary> }
- Preconditions: пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Отсутствует файл в запросе
- Auth: JWT
- Body: multipart/form-data {} (без file)
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999/attachments
- Body: multipart/form-data { file: <binary> }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Без авторизации
- Auth: нет
- Body: multipart/form-data { file: <binary> }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
