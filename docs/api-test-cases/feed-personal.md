# GET /api/v1/feed?mode=personal&cursor=

## 1. Успешное получение личной ленты (JWT)
- Auth: JWT
- Query: ?mode=personal
- Preconditions: есть события по доступным задачам
- Expected: 200
- Response: { items: [{ id, event_type, actor: { id, username }, resource_type, resource_id, payload, created_at }, ...], next_cursor }

## 2. Успешное получение личной ленты (PAT)
- Auth: PAT
- Query: ?mode=personal
- Expected: 200
- Response: { items: [...], next_cursor }

## 3. Пустая лента
- Auth: JWT
- Query: ?mode=personal
- Preconditions: нет событий по доступным задачам
- Expected: 200
- Response: { items: [], next_cursor: null }

## 4. Видны только события по доступным задачам
- Auth: JWT
- Query: ?mode=personal
- Preconditions: есть события по задачам, к которым пользователь не имеет доступа
- Expected: 200
- Response: items содержит только события по доступным задачам

## 5. Учёт окна доступа — события до получения доступа не видны
- Auth: JWT
- Query: ?mode=personal
- Preconditions: пользователь получил доступ к задаче вчера, события по задаче были и ранее
- Expected: 200
- Response: items содержит только события после получения доступа

## 6. Пагинация — первая страница
- Auth: JWT
- Query: ?mode=personal
- Preconditions: много событий
- Expected: 200
- Response: { items: [...], next_cursor: "<non-null>" }

## 7. Пагинация — с cursor
- Auth: JWT
- Query: ?mode=personal&cursor=<next_cursor>
- Expected: 200
- Response: { items: [...], next_cursor }

## 8. Пагинация — последняя страница
- Auth: JWT
- Query: ?mode=personal&cursor=<cursor>
- Expected: 200
- Response: { items: [...], next_cursor: null }

## 9. Невалидный cursor
- Auth: JWT
- Query: ?mode=personal&cursor=invalid
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Без авторизации
- Auth: нет
- Query: ?mode=personal
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
