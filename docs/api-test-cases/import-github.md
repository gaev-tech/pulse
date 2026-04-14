# POST /api/v1/import/github

## 1. Успешный импорт личных задач (JWT)
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Expected: 200
- Response: { import_id }
- Side effects: импорт запущен асинхронно

## 2. Успешный импорт личных задач (PAT)
- Auth: PAT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Expected: 200
- Response: { import_id }

## 3. Импорт в команду — с правом team.import
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь имеет team.import
- Expected: 200
- Response: { import_id }

## 4. Импорт в команду — без права team.import
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь — участник без team.import
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Импорт в команду — не участник команды
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>", "team_id": "<team_id>" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Повторный импорт — обновление существующих задач
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Preconditions: задачи с теми же origin_id (github:acme/app:*) уже существуют
- Expected: 200
- Side effects: существующие задачи обновлены, поля перезаписаны

## 7. Маппинг статусов GitHub
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Preconditions: репозиторий содержит issues с state "open" и "closed"
- Expected: 200
- Side effects: "closed" → closed, "open" → opened

## 8. Маппинг assignee по email
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>", "team_id": "<team_id>" }
- Preconditions: часть assignees из GitHub имеет email, совпадающий с участником команды
- Expected: 200
- Side effects: совпавшие assignee назначены, остальные — без assignee

## 9. Маппинг labels
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Preconditions: issues содержат labels, часть из которых не существует в Pulse
- Expected: 200
- Side effects: несуществующие метки созданы, задачи получили метки

## 10. Маппинг sub-issues (parent-child)
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Preconditions: issues содержат sub-issues
- Expected: 200
- Side effects: parent-child связи восстановлены

## 11. Маппинг tracked-by/tracks (blocking)
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Preconditions: issues содержат tracked-by / tracks relationships
- Expected: 200
- Side effects: blocking связи восстановлены

## 12. Невалидный access_token — ошибка доступа к GitHub API
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "invalid" }
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 13. Несуществующий репозиторий
- Auth: JWT
- Body: { "owner": "acme", "repo": "nonexistent", "access_token": "<token>" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 14. Отсутствует поле owner
- Auth: JWT
- Body: { "repo": "app", "access_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. Отсутствует поле repo
- Auth: JWT
- Body: { "owner": "acme", "access_token": "<token>" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Отсутствует поле access_token
- Auth: JWT
- Body: { "owner": "acme", "repo": "app" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Несуществующий team_id
- Auth: JWT
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>", "team_id": "<nonexistent_id>" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 18. Без авторизации
- Auth: нет
- Body: { "owner": "acme", "repo": "app", "access_token": "<token>" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
