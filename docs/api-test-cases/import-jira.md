# POST /api/v1/import/jira

## 1. Успешный импорт личных задач (JWT)
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Expected: 200
- Response: { import_id }
- Side effects: импорт запущен асинхронно

## 2. Успешный импорт личных задач (PAT)
- Auth: PAT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Expected: 200
- Response: { import_id }

## 3. Импорт в команду — с правом team.import
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь имеет team.import
- Expected: 200
- Response: { import_id }

## 4. Импорт в команду — без права team.import
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь — участник без team.import
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Импорт в команду — не участник команды
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Повторный импорт — обновление существующих задач
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Preconditions: задачи с теми же origin_id (jira:mycompany.atlassian.net:PROJ-*) уже существуют
- Expected: 200
- Side effects: существующие задачи обновлены, поля перезаписаны

## 7. Маппинг статусов Jira
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Preconditions: проект содержит issues со статусами "To Do", "In Progress", "Done"
- Expected: 200
- Side effects: "Done" → closed, "To Do" и "In Progress" → opened

## 8. Маппинг assignee по email
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>", "team_id": "<team_id>" }
- Preconditions: часть assignee из Jira совпадает по email с участниками команды, часть — нет
- Expected: 200
- Side effects: совпавшие assignee назначены, остальные — без assignee

## 9. Маппинг связей — parent-child, blocking, relations
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Preconditions: issues содержат parent links, Blocks links, другие типы issuelinks
- Expected: 200
- Side effects: parent-child и blocking восстановлены; остальные issuelinks → relations

## 10. Невалидный api_token — ошибка доступа к Jira API
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "invalid" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 11. Несуществующий project_key
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "NONEXISTENT", "email": "user@example.com", "api_token": "<token>" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 12. Недоступный instance_url
- Auth: JWT
- Body: { "instance_url": "nonexistent.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Отсутствует поле instance_url
- Auth: JWT
- Body: { "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Отсутствует поле project_key
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "email": "user@example.com", "api_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Отсутствует поле email
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "api_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Отсутствует поле api_token
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Несуществующий team_id
- Auth: JWT
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>", "team_id": "<nonexistent_id>" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 18. Без авторизации
- Auth: нет
- Body: { "instance_url": "mycompany.atlassian.net", "project_key": "PROJ", "email": "user@example.com", "api_token": "<token>" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
