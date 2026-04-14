# GET /api/v1/settings

## 1. Успешное получение настроек (JWT)
- Auth: JWT
- Expected: 200
- Response: { language, theme, sidebar_personal_open, sidebar_teams_open, team_states: [{ team_id, open }] }

## 2. Успешное получение настроек (PAT)
- Auth: PAT
- Expected: 200
- Response: { language, theme, sidebar_personal_open, sidebar_teams_open, team_states: [...] }

## 3. Настройки по умолчанию (новый пользователь)
- Auth: JWT
- Preconditions: пользователь только что зарегистрирован, настройки не изменялись
- Expected: 200
- Response: { language, theme, sidebar_personal_open, sidebar_teams_open, team_states: [] }

## 4. Настройки с team_states
- Auth: JWT
- Preconditions: пользователь состоит в командах, настраивал раскрытие sidebar
- Expected: 200
- Response: { team_states: [{ team_id: "<id>", open: true }, { team_id: "<id>", open: false }] }

## 5. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 6. Невалидный токен
- Auth: JWT (invalid)
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
