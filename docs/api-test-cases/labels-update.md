# PATCH /api/v1/labels/{labelID}

## 1. Успешное изменение name личной метки (JWT)
- Auth: JWT
- Body: { "name": "Critical Bug" }
- Preconditions: метка принадлежит текущему пользователю
- Expected: 200
- Response: { id, name: "Critical Bug", color }
- Side effects: событие label.renamed генерируется во всех задачах с этой меткой

## 2. Успешное изменение name метки (PAT)
- Auth: PAT
- Body: { "name": "Critical Bug" }
- Preconditions: метка принадлежит текущему пользователю
- Expected: 200
- Response: { id, name: "Critical Bug", color }

## 3. Успешное изменение color
- Auth: JWT
- Body: { "color": "#00ff00" }
- Preconditions: метка принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, color: "#00ff00" }

## 4. Изменение name и color одновременно
- Auth: JWT
- Body: { "name": "New Name", "color": "#00ff00" }
- Expected: 200
- Response: { id, name: "New Name", color: "#00ff00" }

## 5. Переименование командной метки — с правом team.manage_labels
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: метка принадлежит команде, пользователь имеет team.manage_labels
- Expected: 200
- Response: { id, name: "New Name", color }
- Side effects: событие label.renamed во всех задачах с этой меткой

## 6. Переименование командной метки — без права team.manage_labels
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь — участник без team.manage_labels
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Изменение чужой личной метки
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: метка принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Командная метка — не участник команды
- Auth: JWT
- Body: { "name": "New Name" }
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 9. Пустая строка в name
- Auth: JWT
- Body: { "name": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Пустая строка в color
- Auth: JWT
- Body: { "color": "" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 11. Несуществующий labelID
- Auth: JWT
- URL: /api/v1/labels/nonexistent-uuid
- Body: { "name": "New" }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 12. Без авторизации
- Auth: нет
- Body: { "name": "New" }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 13. Пустое тело запроса
- Auth: JWT
- Body: {}
- Preconditions: метка принадлежит текущему пользователю
- Expected: 200
- Response: { id, name, color } (без изменений)
