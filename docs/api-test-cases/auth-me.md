# GET /api/v1/auth/me

## 1. Успешное получение данных текущего пользователя (JWT)
- Auth: JWT (valid access_token)
- Expected: 200
- Response: { id, email, username }

## 2. Успешное получение данных текущего пользователя (PAT)
- Auth: PAT (valid token)
- Expected: 200
- Response: { id, email, username }

## 3. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 4. Невалидный JWT
- Auth: JWT (invalid token)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 5. Истёкший JWT (15 минут)
- Auth: JWT (expired access_token)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 6. Невалидный PAT
- Auth: PAT (invalid token)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 7. Отозванный PAT
- Auth: PAT (revoked token)
- Preconditions: PAT был отозван через DELETE /pat или перегенерирован через POST /pat
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
