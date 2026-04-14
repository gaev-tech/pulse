# POST /api/v1/auth/logout

## 1. Успешный logout
- Auth: нет
- Body: { "refresh_token": "<valid_refresh_token>" }
- Preconditions: refresh_token валиден
- Expected: 204 No Content
- Side effects: refresh_token отозван, повторный refresh с этим токеном вернёт 401

## 2. Logout с уже отозванным refresh_token
- Auth: нет
- Body: { "refresh_token": "<already_revoked_token>" }
- Preconditions: токен уже был отозван
- Expected: 204 No Content

## 3. Logout с невалидным refresh_token
- Auth: нет
- Body: { "refresh_token": "invalid-token" }
- Expected: 204 No Content

## 4. Отсутствует поле refresh_token
- Auth: нет
- Body: {}
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 5. Пустая строка в refresh_token
- Auth: нет
- Body: { "refresh_token": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 6. Пустое тело запроса
- Auth: нет
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
