# POST /api/v1/tasks/{key}/permissions

## 1. Успешный шаринг задачи с правами (JWT)
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title", "edit.description"] }
- Preconditions: текущий пользователь имеет права share, edit.title, edit.description на задачу
- Expected: 200
- Response: { permission_id }
- Side effects: целевой пользователь получает права edit.title, edit.description + автоматический view

## 2. Успешный шаринг задачи с правами (PAT)
- Auth: PAT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title"] }
- Preconditions: текущий пользователь имеет права share и edit.title
- Expected: 200
- Response: { permission_id }

## 3. Шаринг с пустым массивом permissions — выдаётся только view
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": [] }
- Preconditions: текущий пользователь имеет право share
- Expected: 200
- Response: { permission_id }
- Side effects: целевой пользователь получает только view

## 4. Попытка выдать право, которого нет у текущего пользователя
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title", "edit.status"] }
- Preconditions: у текущего пользователя есть share и edit.title, но нет edit.status
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Нет права share
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title"] }
- Preconditions: у текущего пользователя есть edit.title, но нет share
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Нет доступа к задаче
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title"] }
- Preconditions: текущий пользователь не имеет доступа к задаче
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Невалидные значения permissions
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["invalid_permission"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 8. Permissions с невалидным типом (строка вместо массива)
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": "edit.title" }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 9. Отсутствует поле user_id
- Auth: JWT
- Body: { "permissions": ["edit.title"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 10. Отсутствует поле permissions
- Auth: JWT
- Body: { "user_id": "<target_user_id>" }
- Expected: 200
- Response: { permission_id }
- Side effects: выдаётся только view

## 11. Несуществующий user_id
- Auth: JWT
- Body: { "user_id": "nonexistent-uuid", "permissions": ["edit.title"] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 12. Шаринг самому себе
- Auth: JWT
- Body: { "user_id": "<current_user_id>", "permissions": ["edit.title"] }
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 13. Повторный шаринг пользователю, у которого уже есть доступ
- Auth: JWT
- Body: { "user_id": "<target_user_id>", "permissions": ["edit.title"] }
- Preconditions: у целевого пользователя уже есть доступ к задаче
- Expected: 409
- Response: { error: { code: "CONFLICT", message: "..." } }

## 14. Несуществующий ключ задачи
- Auth: JWT
- URL: /api/v1/tasks/NONEXISTENT-999/permissions
- Body: { "user_id": "<target_user_id>", "permissions": [] }
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 15. Без авторизации
- Auth: нет
- Body: { "user_id": "<target_user_id>", "permissions": [] }
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }

## 16. Пустое тело запроса
- Auth: JWT
- Body: нет
- Expected: 400
- Response: { error: { code: "VALIDATION_ERROR", message: "..." } }

## 17. Шаринг — цепочка прав (A шарит B, B шарит C)
- Auth: JWT (пользователь B)
- Body: { "user_id": "<user_C_id>", "permissions": ["edit.title"] }
- Preconditions: A выдал B права share и edit.title; B шарит C
- Expected: 200
- Response: { permission_id }
- Side effects: C получает edit.title + view; цепочка A→B→C сохраняется
