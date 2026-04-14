# GET /api/v1/feed?mode=team&team_prefix=&actor_ids[]=&cursor=

## 1. Успешное получение командной ленты (JWT)
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME
- Preconditions: пользователь — участник команды ACME, есть события
- Expected: 200
- Response: { items: [{ id, event_type, actor: { id, username }, resource_type, resource_id, payload, created_at }, ...], next_cursor }

## 2. Успешное получение командной ленты (PAT)
- Auth: PAT
- Query: ?mode=team&team_prefix=ACME
- Expected: 200
- Response: { items: [...], next_cursor }

## 3. Фильтрация по actor_ids — один автор
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME&actor_ids[]=<user_id>
- Expected: 200
- Response: items содержит только события от указанного актора

## 4. Фильтрация по actor_ids — несколько авторов
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME&actor_ids[]=<user_id_1>&actor_ids[]=<user_id_2>
- Expected: 200
- Response: items содержит события только от user_id_1 и user_id_2

## 5. Без actor_ids — все события команды
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME
- Expected: 200
- Response: items содержит все события по задачам команды

## 6. Пустая командная лента
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME
- Preconditions: нет событий в команде
- Expected: 200
- Response: { items: [], next_cursor: null }

## 7. Пагинация — с cursor
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME&cursor=<next_cursor>
- Expected: 200
- Response: { items: [...], next_cursor }

## 8. Не участник команды
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 9. Несуществующий team_prefix
- Auth: JWT
- Query: ?mode=team&team_prefix=NONEXISTENT
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Отсутствует team_prefix
- Auth: JWT
- Query: ?mode=team
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Невалидный cursor
- Auth: JWT
- Query: ?mode=team&team_prefix=ACME&cursor=invalid
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. Без авторизации
- Auth: нет
- Query: ?mode=team&team_prefix=ACME
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
