# GET /api/v1/tasks/{key}

## 1. Успешное получение задачи по ключу (JWT)
- Auth: JWT
- Preconditions: задача существует, пользователь имеет право view
- Expected: 200
- Response: { id, key, title, description, status, owner_type, owner_id, assignee_id, parent_id, label_ids, links, relations, blocking, task_permissions: number, created_at, updated_at }
- Side effects: фиксируется факт открытия в task_opens

## 2. Успешное получение задачи по ключу (PAT)
- Auth: PAT
- Preconditions: задача существует, пользователь имеет право view
- Expected: 200
- Response: { ...task, task_permissions: number }
- Side effects: фиксируется факт открытия в task_opens

## 3. Личная задача по ключу #123
- Auth: JWT
- URL: /api/v1/tasks/%23123
- Preconditions: личная задача текущего пользователя с key "#123"
- Expected: 200
- Response: { key: "#123", owner_type: "user", ... }

## 4. Командная задача по ключу PREFIX-123
- Auth: JWT
- URL: /api/v1/tasks/ACME-123
- Preconditions: командная задача, пользователь — участник команды
- Expected: 200
- Response: { key: "ACME-123", owner_type: "team", ... }

## 5. Проверка task_permissions — владелец задачи (все права)
- Auth: JWT
- Preconditions: пользователь — владелец личной задачи
- Expected: 200
- Response: task_permissions содержит все биты (view + все edit + share)

## 6. Проверка task_permissions — участник команды с ограниченными правами
- Auth: JWT
- Preconditions: пользователь — участник команды с правами [edit.title, edit.status]
- Expected: 200
- Response: task_permissions содержит биты view, edit.title, edit.status

## 7. Проверка task_permissions — пользователь с direct share
- Auth: JWT
- Preconditions: задача расшарена с правами [edit.title]
- Expected: 200
- Response: task_permissions содержит биты view, edit.title

## 8. Нет доступа к задаче
- Auth: JWT
- Preconditions: пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 9. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
