# POST /api/v1/auth/refresh

## 1. Успешное обновление токенов
- Auth: нет
- Body: { "refresh_token": "<valid_refresh_token>" }
- Preconditions: refresh_token валиден и не отозван
- Expected: 200
- Response: { access_token: "...", refresh_token: "..." }

## 2. Невалидный refresh_token
- Auth: нет
- Body: { "refresh_token": "invalid-token" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 3. Отозванный refresh_token (после logout)
- Auth: нет
- Body: { "refresh_token": "<revoked_token>" }
- Preconditions: токен был отозван через POST /auth/logout
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 4. Истёкший refresh_token (30 дней)
- Auth: нет
- Body: { "refresh_token": "<expired_token>" }
- Preconditions: токен старше 30 дней
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 5. Отсутствует поле refresh_token
- Auth: нет
- Body: {}
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 6. Пустая строка в refresh_token
- Auth: нет
- Body: { "refresh_token": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Пустое тело запроса
- Auth: нет
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
