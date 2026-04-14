# GET /api/v1/filters/{filterID}/settings

## 1. Успешное получение настроек личного фильтра (JWT)
- Auth: JWT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }

## 2. Успешное получение настроек (PAT)
- Auth: PAT
- Preconditions: фильтр принадлежит текущему пользователю
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }

## 3. Командный фильтр — участник команды
- Auth: JWT
- Preconditions: фильтр принадлежит команде, пользователь — участник
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }

## 4. Настройки по умолчанию (новый фильтр)
- Auth: JWT
- Preconditions: фильтр только создан, настройки не изменялись
- Expected: 200
- Response: { columns, sort1_column, sort1_dir, sort2_column: null, sort2_dir: null }

## 5. Командный фильтр — не участник команды
- Auth: JWT
- Preconditions: пользователь не участник команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 6. Чужой личный фильтр
- Auth: JWT
- Preconditions: фильтр принадлежит другому пользователю
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Несуществующий filterID
- Auth: JWT
- URL: /api/v1/filters/nonexistent-uuid/settings
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 8. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
