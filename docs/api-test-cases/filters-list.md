# GET /api/v1/filters

## 1. Успешное получение фильтров (JWT)
- Auth: JWT
- Preconditions: у пользователя есть личные фильтры и он состоит в командах с фильтрами
- Expected: 200
- Response: { personal_filters: [{ id, name }], teams: [{ id, name, prefix, team_permissions: number, filters: [{ id, name }] }] }

## 2. Успешное получение фильтров (PAT)
- Auth: PAT
- Expected: 200
- Response: { personal_filters: [...], teams: [...] }

## 3. Нет личных фильтров и не в командах
- Auth: JWT
- Preconditions: пользователь не имеет фильтров и не участник команд
- Expected: 200
- Response: { personal_filters: [], teams: [] }

## 4. Есть личные фильтры, нет команд
- Auth: JWT
- Preconditions: у пользователя есть личные фильтры, но он не в командах
- Expected: 200
- Response: { personal_filters: [{ id, name }, ...], teams: [] }

## 5. Проверка team_permissions в ответе
- Auth: JWT
- Preconditions: пользователь — участник команды с правами [task.create, team.manage_filters]
- Expected: 200
- Response: teams[].team_permissions содержит соответствующие биты

## 6. Команда без фильтров
- Auth: JWT
- Preconditions: пользователь — участник команды, у команды нет фильтров
- Expected: 200
- Response: teams содержит команду с filters: []

## 7. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
