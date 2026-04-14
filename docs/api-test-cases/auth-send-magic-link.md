# POST /api/v1/auth/magic-link

## 1. Успешная отправка magic-link на существующий email
- Auth: нет
- Body: { "email": "user@example.com" }
- Preconditions: пользователь с email user@example.com существует
- Expected: 200 OK
- Side effects: отправлено письмо со ссылкой

## 2. Успешная отправка magic-link на несуществующий email
- Auth: нет
- Body: { "email": "newuser@example.com" }
- Preconditions: пользователя с таким email нет
- Expected: 200 OK
- Side effects: отправлено письмо со ссылкой (пользователь создаётся при верификации токена)

## 3. Повторная отправка magic-link на тот же email
- Auth: нет
- Body: { "email": "user@example.com" }
- Preconditions: magic-link уже отправлен на этот email
- Expected: 200 OK
- Side effects: предыдущий токен инвалидирован, отправлено новое письмо

## 4. Отсутствует поле email
- Auth: нет
- Body: {}
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 5. Пустая строка в email
- Auth: нет
- Body: { "email": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 6. Невалидный формат email
- Auth: нет
- Body: { "email": "not-an-email" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 7. Email с невалидным типом (число)
- Auth: нет
- Body: { "email": 123 }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Email с невалидным типом (null)
- Auth: нет
- Body: { "email": null }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Пустое тело запроса
- Auth: нет
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }
