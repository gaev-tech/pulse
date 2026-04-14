# POST /api/v1/import/csv

## 1. Успешный импорт личных задач (JWT)
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id" }
- Expected: 200
- Response: { import_id }
- Side effects: импорт запущен асинхронно

## 2. Успешный импорт личных задач (PAT)
- Auth: PAT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id" }
- Expected: 200
- Response: { import_id }

## 3. Импорт в команду — с правом team.import
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id", team_id: "<team_id>" }
- Preconditions: пользователь имеет team.import
- Expected: 200
- Response: { import_id }

## 4. Импорт в команду — без права team.import
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id", team_id: "<team_id>" }
- Preconditions: пользователь — участник без team.import
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Импорт в команду — не участник команды
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id", team_id: "<team_id>" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Повторный импорт — обновление существующих задач
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: updated.csv, id_column: "id" }
- Preconditions: задачи с теми же origin_id уже существуют от предыдущего импорта
- Expected: 200
- Side effects: существующие задачи обновлены, поля перезаписаны

## 7. Импорт с маппингом assignee по email
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: with_assignees.csv, id_column: "id", team_id: "<team_id>" }
- Preconditions: CSV содержит колонку assignee_email, один email совпадает с участником команды, другой — нет
- Expected: 200
- Side effects: задача с совпавшим email получает assignee, остальные — без assignee

## 8. Импорт с labels — создание несуществующих меток
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: with_labels.csv, id_column: "id" }
- Preconditions: CSV содержит колонку labels, часть меток не существует
- Expected: 200
- Side effects: несуществующие метки созданы, задачи получили метки

## 9. Импорт с parent-child связями
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: with_parents.csv, id_column: "id" }
- Preconditions: CSV содержит колонку parent со ссылками на id_column других строк
- Expected: 200
- Side effects: parent-child связи между задачами восстановлены

## 10. Частичный импорт — строки без title пропускаются
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: partial.csv, id_column: "id" }
- Preconditions: CSV содержит строки без обязательного поля title
- Expected: 200
- Side effects: валидные задачи импортированы, невалидные попали в errors результата

## 11. Маппинг статусов — closed и остальные
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: statuses.csv, id_column: "id" }
- Preconditions: CSV содержит строки со статусами "closed", "opened", "in progress", "done"
- Expected: 200
- Side effects: "closed" → closed, все остальные → opened

## 12. Отсутствует файл
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { id_column: "id" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Отсутствует поле id_column
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 14. Указанная id_column отсутствует в CSV
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "nonexistent_column" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 15. CSV без колонки title
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: no_title.csv, id_column: "id" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 16. Пустой CSV (только заголовки)
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: empty.csv, id_column: "id" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Несуществующий team_id
- Auth: JWT
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id", team_id: "<nonexistent_id>" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 18. Без авторизации
- Auth: нет
- Content-Type: multipart/form-data
- Body: { file: valid.csv, id_column: "id" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
