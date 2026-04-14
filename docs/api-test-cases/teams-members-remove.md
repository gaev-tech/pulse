# DELETE /api/v1/teams/{prefix}/members/{userID}

## 1. Успешное удаление участника (JWT)
- Auth: JWT
- Preconditions: пользователь — владелец команды, целевой участник существует
- Expected: 204 No Content
- Side effects: каскадный отзыв прав участника и прав, выданных им другим пользователям по командным задачам

## 2. Успешное удаление участника (PAT)
- Auth: PAT
- Preconditions: пользователь — владелец команды
- Expected: 204 No Content
- Side effects: каскадный отзыв прав

## 3. Каскадное удаление прав, выданных удалённым участником
- Auth: JWT
- Preconditions: участник A расшарил задачу пользователю B; удаляем A из команды
- Expected: 204 No Content
- Side effects: права, выданные A пользователю B по командным задачам, каскадно удалены

## 4. Удаление владельца команды
- Auth: JWT
- Preconditions: userID — владелец команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 5. Участник с team.manage_owners удаляет другого участника
- Auth: JWT
- Preconditions: пользователь — участник с team.manage_owners
- Expected: 204 No Content

## 6. Участник без team.manage_owners
- Auth: JWT
- Preconditions: пользователь — участник без team.manage_owners
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 7. Не участник команды
- Auth: JWT
- Preconditions: пользователь не является участником команды
- Expected: 403
- Response: { error: { code: "PERMISSION_DENIED", message: "..." } }

## 8. Несуществующий userID
- Auth: JWT
- URL: /api/v1/teams/ACME/members/nonexistent-uuid
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 9. Несуществующий prefix команды
- Auth: JWT
- URL: /api/v1/teams/NONEXISTENT/members/{userID}
- Expected: 404
- Response: { error: { code: "NOT_FOUND", message: "..." } }

## 10. Без авторизации
- Auth: нет
- Expected: 401
- Response: { error: { code: "UNAUTHORIZED", message: "..." } }
