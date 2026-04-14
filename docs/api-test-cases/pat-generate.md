# POST /api/v1/pat

## 1. Успешная генерация PAT (JWT)
- Auth: JWT
- Expected: 200
- Response: { token: "..." }
- Side effects: если был предыдущий PAT — он инвалидирован

## 2. Генерация PAT инвалидирует предыдущий
- Auth: JWT
- Preconditions: у пользователя уже есть активный PAT
- Expected: 200
- Response: { token: "..." } (новый токен)
- Side effects: предыдущий PAT больше не работает; запросы с ним возвращают 401

## 3. Попытка генерации через PAT — запрещено
- Auth: PAT
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 4. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 5. Невалидный JWT
- Auth: JWT (invalid)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 6. Истёкший JWT
- Auth: JWT (expired)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 7. Токен возвращается единожды — повторного получения нет
- Auth: JWT
- Preconditions: PAT уже сгенерирован
- Expected: при повторном POST генерируется НОВЫЙ PAT (старый инвалидируется)
