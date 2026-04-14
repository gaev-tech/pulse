# PATCH /api/v1/tasks/{key}

## 1. Успешное изменение title (JWT)
- Auth: JWT
- Body: { "title": "Updated Title" }
- Preconditions: пользователь имеет право edit.title
- Expected: 200
- Response: { ...task, title: "Updated Title", task_permissions: number }

## 2. Успешное изменение title (PAT)
- Auth: PAT
- Body: { "title": "Updated Title" }
- Preconditions: пользователь имеет право edit.title
- Expected: 200
- Response: { ...task, title: "Updated Title", task_permissions: number }

## 3. Успешное изменение description
- Auth: JWT
- Body: { "description": "## New Description" }
- Preconditions: пользователь имеет право edit.description
- Expected: 200
- Response: { ...task, description: "## New Description" }

## 4. Успешное изменение status на closed
- Auth: JWT
- Body: { "status": "closed" }
- Preconditions: пользователь имеет право edit.status, задача в статусе opened
- Expected: 200
- Response: { ...task, status: "closed" }

## 5. Успешное изменение status на opened
- Auth: JWT
- Body: { "status": "opened" }
- Preconditions: пользователь имеет право edit.status, задача в статусе closed
- Expected: 200
- Response: { ...task, status: "opened" }

## 6. Успешное изменение assignee_id
- Auth: JWT
- Body: { "assignee_id": "<user_id>" }
- Preconditions: пользователь имеет право edit.assignee
- Expected: 200
- Response: { ...task, assignee_id: "<user_id>" }

## 7. Снятие ответственного (assignee_id: null)
- Auth: JWT
- Body: { "assignee_id": null }
- Preconditions: пользователь имеет право edit.assignee, задача имеет assignee
- Expected: 200
- Response: { ...task, assignee_id: null }

## 8. Успешное изменение label_ids
- Auth: JWT
- Body: { "label_ids": ["<label_id_1>", "<label_id_2>"] }
- Preconditions: пользователь имеет право edit.labels, метки принадлежат тому же owner
- Expected: 200
- Response: { ...task, label_ids: ["<label_id_1>", "<label_id_2>"] }

## 9. Удаление всех меток (label_ids: [])
- Auth: JWT
- Body: { "label_ids": [] }
- Preconditions: пользователь имеет право edit.labels
- Expected: 200
- Response: { ...task, label_ids: [] }

## 10. Успешное изменение links
- Auth: JWT
- Body: { "links": [{ "url": "https://example.com", "title": "Example" }] }
- Preconditions: пользователь имеет право edit.links
- Expected: 200
- Response: { ...task, links: [{ url, title }] }

## 11. Успешное изменение relations
- Auth: JWT
- Body: { "relations": ["ACME-42", "#5"] }
- Preconditions: пользователь имеет право edit.relations, указанные задачи существуют
- Expected: 200
- Response: { ...task, relations: ["ACME-42", "#5"] }

## 12. Успешное изменение blocking
- Auth: JWT
- Body: { "blocking": ["ACME-43"] }
- Preconditions: пользователь имеет право edit.blocking, нет цикла
- Expected: 200
- Response: { ...task, blocking: ["ACME-43"] }

## 13. Blocking с циклом — ошибка
- Auth: JWT
- Body: { "blocking": ["ACME-44"] }
- Preconditions: ACME-44 уже блокирует текущую задачу (прямо или транзитивно)
- Expected: 409
- Response: { error: { code: "CONFLICT", message: "..." } }

## 14. Успешное изменение parent_id
- Auth: JWT
- Body: { "parent_id": "<parent_task_id>" }
- Preconditions: пользователь имеет право edit.parent, родительская задача существует
- Expected: 200
- Response: { ...task, parent_id: "<parent_task_id>" }

## 15. Снятие родительской задачи (parent_id: null)
- Auth: JWT
- Body: { "parent_id": null }
- Preconditions: пользователь имеет право edit.parent
- Expected: 200
- Response: { ...task, parent_id: null }

## 16. Изменение нескольких полей одновременно
- Auth: JWT
- Body: { "title": "New Title", "status": "closed", "assignee_id": null }
- Preconditions: пользователь имеет права edit.title, edit.status, edit.assignee
- Expected: 200
- Response: { ...task, title: "New Title", status: "closed", assignee_id: null }

## 17. Нет права edit.title — изменение title
- Auth: JWT
- Body: { "title": "New Title" }
- Preconditions: пользователь имеет view, но не edit.title
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 18. Нет права edit.description — изменение description
- Auth: JWT
- Body: { "description": "New" }
- Preconditions: пользователь имеет view, но не edit.description
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 19. Нет права edit.status — изменение status
- Auth: JWT
- Body: { "status": "closed" }
- Preconditions: пользователь имеет view, но не edit.status
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 20. Нет права edit.assignee — изменение assignee_id
- Auth: JWT
- Body: { "assignee_id": "<user_id>" }
- Preconditions: пользователь имеет view, но не edit.assignee
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 21. Нет права edit.labels — изменение label_ids
- Auth: JWT
- Body: { "label_ids": ["<label_id>"] }
- Preconditions: пользователь имеет view, но не edit.labels
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 22. Нет права edit.links — изменение links
- Auth: JWT
- Body: { "links": [{ "url": "https://example.com" }] }
- Preconditions: пользователь имеет view, но не edit.links
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 23. Нет права edit.relations — изменение relations
- Auth: JWT
- Body: { "relations": ["ACME-42"] }
- Preconditions: пользователь имеет view, но не edit.relations
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 24. Нет права edit.blocking — изменение blocking
- Auth: JWT
- Body: { "blocking": ["ACME-43"] }
- Preconditions: пользователь имеет view, но не edit.blocking
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 25. Нет права edit.parent — изменение parent_id
- Auth: JWT
- Body: { "parent_id": "<task_id>" }
- Preconditions: пользователь имеет view, но не edit.parent
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 26. Нет доступа к задаче
- Auth: JWT
- Body: { "title": "New Title" }
- Preconditions: пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 27. Невалидный status
- Auth: JWT
- Body: { "status": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 28. Пустая строка в title
- Auth: JWT
- Body: { "title": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 29. label_ids с чужими метками
- Auth: JWT
- Body: { "label_ids": ["<foreign_label_id>"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 30. Relations с самой собой
- Auth: JWT
- Body: { "relations": ["ACME-42"] }
- Preconditions: текущая задача — ACME-42
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 31. Relations с несуществующей задачей
- Auth: JWT
- Body: { "relations": ["NONEXISTENT-999"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 32. Blocking с несуществующей задачей
- Auth: JWT
- Body: { "blocking": ["NONEXISTENT-999"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 33. Parent_id с несуществующей задачей
- Auth: JWT
- Body: { "parent_id": "nonexistent-uuid" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 34. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999
- Body: { "title": "New Title" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 35. Без авторизации
- Auth: нет
- Body: { "title": "New Title" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 36. Пустое тело запроса (нет изменений)
- Auth: JWT
- Body: {}
- Preconditions: пользователь имеет доступ к задаче
- Expected: 200
- Response: { ...task } (без изменений)

## 37. Участник команды с edit.title обновляет title командной задачи
- Auth: JWT
- Body: { "title": "Updated" }
- Preconditions: пользователь — участник команды с правом edit.title
- Expected: 200

## 38. Пользователь с direct share и edit.title обновляет title
- Auth: JWT
- Body: { "title": "Updated" }
- Preconditions: задача расшарена с правом edit.title
- Expected: 200
