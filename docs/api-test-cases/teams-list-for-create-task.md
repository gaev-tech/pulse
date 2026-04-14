# GET /api/v1/teams-for-create-task

## 1. Успешное получение списка команд для создания задачи (JWT)
- Auth: JWT
- Preconditions: пользователь состоит в командах, в некоторых имеет право task.create
- Expected: 200
- Response: [{ id, name, prefix }, ...] — только команды с правом task.create

## 2. Успешное получение списка команд для создания задачи (PAT)
- Auth: PAT
- Preconditions: пользователь имеет право task.create в нескольких командах
- Expected: 200
- Response: [{ id, name, prefix }, ...]

## 3. Пользователь не имеет права task.create ни в одной команде
- Auth: JWT
- Preconditions: пользователь состоит в командах, но без права task.create
- Expected: 200
- Response: []

## 4. Пользователь не состоит ни в одной команде
- Auth: JWT
- Preconditions: пользователь не участник ни одной команды
- Expected: 200
- Response: []

## 5. Владелец команды всегда имеет task.create
- Auth: JWT
- Preconditions: пользователь — владелец команды
- Expected: 200
- Response: команда присутствует в списке

## 6. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
