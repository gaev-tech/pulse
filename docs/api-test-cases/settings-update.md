# PATCH /api/v1/settings

## 1. Успешное изменение language (JWT)
- Auth: JWT
- Body: { "language": "en" }
- Expected: 200
- Response: { language: "en", theme, sidebar_personal_open, sidebar_teams_open, team_states }

## 2. Успешное изменение настроек (PAT)
- Auth: PAT
- Body: { "language": "ru" }
- Expected: 200

## 3. Изменение theme
- Auth: JWT
- Body: { "theme": "dark" }
- Expected: 200
- Response: { theme: "dark", ... }

## 4. Изменение sidebar_personal_open
- Auth: JWT
- Body: { "sidebar_personal_open": false }
- Expected: 200
- Response: { sidebar_personal_open: false, ... }

## 5. Изменение sidebar_teams_open
- Auth: JWT
- Body: { "sidebar_teams_open": true }
- Expected: 200
- Response: { sidebar_teams_open: true, ... }

## 6. Изменение team_states
- Auth: JWT
- Body: { "team_states": [{ "team_id": "<team_id>", "open": false }] }
- Expected: 200
- Response: { team_states: [{ team_id: "<team_id>", open: false }], ... }

## 7. Изменение нескольких настроек одновременно
- Auth: JWT
- Body: { "language": "en", "theme": "light", "sidebar_personal_open": true }
- Expected: 200
- Response: обновлены все переданные поля

## 8. Частичное обновление — не переданные поля не меняются
- Auth: JWT
- Body: { "language": "en" }
- Preconditions: theme = "dark"
- Expected: 200
- Response: { language: "en", theme: "dark", ... } — theme не изменился

## 9. team_states с несуществующим team_id
- Auth: JWT
- Body: { "team_states": [{ "team_id": "nonexistent-uuid", "open": true }] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Невалидное значение language
- Auth: JWT
- Body: { "language": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Невалидное значение theme
- Auth: JWT
- Body: { "theme": "invalid" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 12. sidebar_personal_open с невалидным типом (строка)
- Auth: JWT
- Body: { "sidebar_personal_open": "true" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Без авторизации
- Auth: нет
- Body: { "language": "en" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 14. Пустое тело запроса
- Auth: JWT
- Body: {}
- Expected: 200
- Response: { language, theme, sidebar_personal_open, sidebar_teams_open, team_states } (без изменений)
