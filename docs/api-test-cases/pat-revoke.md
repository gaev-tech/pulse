# DELETE /api/v1/pat

## 1. Успешный отзыв PAT (JWT)
- Auth: JWT
- Preconditions: у пользователя есть активный PAT
- Expected: 204 No Content
- Side effects: PAT больше не работает

## 2. Отзыв при отсутствии активного PAT
- Auth: JWT
- Preconditions: у пользователя нет активного PAT
- Expected: 204 No Content

## 3. Попытка отзыва через PAT — запрещено
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

## 6. После отзыва — запросы с PAT возвращают 401
- Auth: JWT (для отзыва), затем PAT (для проверки)
- Preconditions: PAT отозван через DELETE /pat
- Expected: любой запрос с отозванным PAT → 401
