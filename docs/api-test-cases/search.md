# GET /api/v1/search?q=&cursor=

## 1. Успешный поиск по названию задачи (JWT)
- Auth: JWT
- Query: ?q=bug
- Preconditions: есть задачи с "bug" в названии
- Expected: 200
- Response: { items: [{ id, key, title, status, owner_type, owner_id }, ...], next_cursor }

## 2. Успешный поиск (PAT)
- Auth: PAT
- Query: ?q=task
- Expected: 200
- Response: { items: [...], next_cursor }

## 3. Поиск по ключу задачи (полное совпадение)
- Auth: JWT
- Query: ?q=ACME-42
- Preconditions: задача с ключом ACME-42 существует, пользователь имеет доступ
- Expected: 200
- Response: { items: [...] } — задача ACME-42 первая в результатах

## 4. Поиск по личному ключу задачи
- Auth: JWT
- Query: ?q=#123
- Preconditions: личная задача #123 существует
- Expected: 200
- Response: { items: [...] } — задача #123 первая в результатах

## 5. Поиск по описанию (contains)
- Auth: JWT
- Query: ?q=urgent
- Preconditions: есть задачи с "urgent" в описании
- Expected: 200
- Response: { items: [...] }

## 6. Сортировка — точное совпадение ключа первым
- Auth: JWT
- Query: ?q=ACME-42
- Preconditions: есть ACME-42 и ACME-421
- Expected: 200
- Response: ACME-42 идёт перед ACME-421

## 7. Сортировка — точное совпадение названия первым
- Auth: JWT
- Query: ?q=Deploy script
- Preconditions: есть задача с точным названием "Deploy script" и задача "Deploy script v2"
- Expected: 200
- Response: задача с точным совпадением названия первая

## 8. Сортировка — по дате последнего открытия
- Auth: JWT
- Query: ?q=task
- Preconditions: несколько задач содержат "task", у некоторых есть записи в task_opens
- Expected: 200
- Response: после точных совпадений — задачи отсортированы по дате последнего открытия (от новых к старым)

## 9. Поиск — нет результатов
- Auth: JWT
- Query: ?q=xyznonexistent
- Expected: 200
- Response: { items: [], next_cursor: null }

## 10. Поиск — видны только доступные задачи
- Auth: JWT
- Query: ?q=secret
- Preconditions: есть задачи с "secret", но к некоторым пользователь не имеет доступа
- Expected: 200
- Response: items содержит только задачи, к которым пользователь имеет доступ

## 11. Пагинация — с cursor
- Auth: JWT
- Query: ?q=task&cursor=<next_cursor>
- Expected: 200
- Response: { items: [...], next_cursor }

## 12. Пагинация — последняя страница
- Auth: JWT
- Query: ?q=task&cursor=<cursor>
- Expected: 200
- Response: { items: [...], next_cursor: null }

## 13. Отсутствует параметр q
- Auth: JWT
- Query: (без q)
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Пустая строка в q
- Auth: JWT
- Query: ?q=
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Без авторизации
- Auth: нет
- Query: ?q=task
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
