# GET /api/v1/filters/{filterID}/feed?cursor=

## 1. Успешное получение ленты по фильтру (JWT)
- Auth: JWT
- Query: (без cursor)
- Preconditions: фильтр принадлежит текущему пользователю, есть события
- Expected: 200
- Response: { items: [{ id, event_type, actor: { id, username }, resource_type, resource_id, payload, created_at }, ...], next_cursor }

## 2. Успешное получение ленты по фильтру (PAT)
- Auth: PAT
- Query: (без cursor)
- Expected: 200
- Response: { items: [...], next_cursor }

## 3. Командный фильтр — участник команды
- Auth: JWT
- Preconditions: фильтр принадлежит команде, пользователь — участник
- Expected: 200
- Response: { items: [...], next_cursor }

## 4. Лента содержит события filter.task_entered и filter.task_left
- Auth: JWT
- Preconditions: задачи входили и выходили из результатов фильтра
- Expected: 200
- Response: items содержит события с event_type filter.task_entered и filter.task_left

## 5. Лента содержит события filter.created, filter.updated, filter.deleted
- Auth: JWT
- Preconditions: фильтр был изменён
- Expected: 200
- Response: items содержит события с event_type filter.updated

## 6. Лента содержит события по задачам в период их вхождения в результаты
- Auth: JWT
- Preconditions: задача была в результатах фильтра и имела изменения
- Expected: 200
- Response: items содержит события задачи за период вхождения

## 7. Пустая лента
- Auth: JWT
- Preconditions: нет событий по фильтру
- Expected: 200
- Response: { items: [], next_cursor: null }

## 8. Пагинация — с cursor
- Auth: JWT
- Query: ?cursor=<next_cursor>
- Expected: 200
- Response: { items: [...], next_cursor }

## 9. Пагинация — последняя страница
- Auth: JWT
- Query: ?cursor=<cursor>
- Expected: 200
- Response: { items: [...], next_cursor: null }

## 10. Командный фильтр — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Чужой личный фильтр
- Auth: JWT
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 12. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid/feed
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 13. Невалидный cursor
- Auth: JWT
- Query: ?cursor=invalid
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
