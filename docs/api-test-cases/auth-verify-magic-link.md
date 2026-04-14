# POST /api/v1/auth/magic-link/verify

## 1. Успешная верификация токена
- Auth: нет
- Body: { "token": "<valid_magic_link_token>" }
- Preconditions: magic-link отправлен, токен не использован, не истёк
- Expected: 200
- Response: { access_token: "...", refresh_token: "...", user: { id, email, username } }

## 2. Невалидный токен
- Auth: нет
- Body: { "token": "invalid-token-value" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 3. Использование токена повторно
- Auth: нет
- Body: { "token": "<already_used_token>" }
- Preconditions: токен уже был использован для верификации
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 4. Истёкший токен
- Auth: нет
- Body: { "token": "<expired_token>" }
- Preconditions: токен просрочен
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 5. Токен после повторной отправки magic-link (старый токен)
- Auth: нет
- Body: { "token": "<old_token>" }
- Preconditions: после отправки этого токена был запрошен новый magic-link на тот же email
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 6. Отсутствует поле token
- Auth: нет
- Body: {}
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Пустая строка в token
- Auth: нет
- Body: { "token": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Token с невалидным типом (число)
- Auth: нет
- Body: { "token": 123 }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Пустое тело запроса
- Auth: нет
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
