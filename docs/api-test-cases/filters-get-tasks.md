# GET /api/v1/filters/{filterID}/tasks?cursor=

## 1. Успешное получение задач личного фильтра (JWT)
- Auth: JWT
- Query: (без cursor)
- Preconditions: фильтр принадлежит текущему пользователю, есть задачи, соответствующие критериям
- Expected: 200
- Response: { items: [{ ...task, task_permissions: number }, ...], next_cursor }

## 2. Успешное получение задач фильтра (PAT)
- Auth: PAT
- Query: (без cursor)
- Expected: 200
- Response: { items: [...], next_cursor }

## 3. Командный фильтр — участник команды
- Auth: JWT
- Preconditions: фильтр принадлежит команде, пользователь — участник
- Expected: 200
- Response: { items: [...], next_cursor }

## 4. Пагинация — первая страница
- Auth: JWT
- Query: (без cursor)
- Preconditions: фильтр возвращает больше задач, чем limit
- Expected: 200
- Response: { items: [...], next_cursor: "<non-null>" }

## 5. Пагинация — вторая страница по cursor
- Auth: JWT
- Query: ?cursor=<next_cursor_from_previous_response>
- Expected: 200
- Response: { items: [...], next_cursor }

## 6. Пагинация — последняя страница
- Auth: JWT
- Query: ?cursor=<cursor>
- Preconditions: это последняя страница
- Expected: 200
- Response: { items: [...], next_cursor: null }

## 7. Пагинация — с параметром limit
- Auth: JWT
- Query: ?limit=10
- Expected: 200
- Response: { items: (максимум 10 элементов), next_cursor }

## 8. Пагинация — limit больше максимума (100)
- Auth: JWT
- Query: ?limit=200
- Expected: 200
- Response: items содержит максимум 100 элементов

## 9. Пагинация — limit=0
- Auth: JWT
- Query: ?limit=0
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Пагинация — невалидный cursor
- Auth: JWT
- Query: ?cursor=invalid-cursor
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Фильтр без результатов
- Auth: JWT
- Preconditions: нет задач, соответствующих критериям
- Expected: 200
- Response: { items: [], next_cursor: null }

## 12. Проверка task_permissions в items
- Auth: JWT
- Preconditions: пользователь имеет разные права на разные задачи в результатах
- Expected: 200
- Response: каждая задача содержит корректный task_permissions

## 13. Командный фильтр — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 14. Чужой личный фильтр
- Auth: JWT
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 15. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid/tasks
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 16. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
