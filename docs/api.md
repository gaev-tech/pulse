# API

Base URL: `/api/v1`

Два способа авторизации — см. [requirements.md](requirements.md#api-authentication).

Cursor-based пагинация: `?cursor=` в запросе, `next_cursor` в ответе (null если страниц больше нет). См. подробнее в [requirements.md](requirements.md#cursor-based-pagination).

## Error format

Все ошибки возвращаются в едином формате:

```json
{
  "error": {
    "code": "PERMISSION_DENIED",
    "message": "You don't have permission to edit this task"
  }
}
```

| HTTP | code | Когда |
|---|---|---|
| 400 | `VALIDATION_ERROR` | Невалидные поля запроса |
| 401 | `UNAUTHORIZED` | Отсутствует или невалидный токен |
| 403 | `PERMISSION_DENIED` | Нет прав на действие |
| 403 | `QUOTA_EXCEEDED` | Превышен лимит плана подписки |
| 404 | `NOT_FOUND` | Ресурс не найден |
| 409 | `CONFLICT` | Конфликт (например, цикл в цепочке блокировок) |

---

## Auth

```
POST /api/v1/auth/magic-link
  Body: { email }
  Response: 200 OK
  Auth: нет

POST /api/v1/auth/magic-link/verify
  Body: { token }
  Response: { access_token, refresh_token, user: { id, email, username } }
  Auth: нет

POST /api/v1/auth/refresh
  Body: { refresh_token }
  Response: { access_token, refresh_token }
  Auth: нет

POST /api/v1/auth/logout
  Body: { refresh_token }
  Response: 204 No Content
  Auth: нет

GET /api/v1/auth/me
  Response: { id, email, username }
  Auth: JWT / PAT
```

---

## Users

```
GET /api/v1/teams/{prefix}/members/search?email=
  Response: [{ id, email, username }]
  Auth: JWT / PAT

GET /api/v1/tasks/{key}/share/search?email=
  Response: [{ id, email, username }]
  Auth: JWT / PAT
```

---

## Teams

```
POST /api/v1/teams
  Body: { name, prefix }
  Response: { id, name, prefix, owner_id, created_at }
  Auth: JWT / PAT

GET /api/v1/teams
  Response: [{ id, name, prefix, owner_id, team_permissions: number, created_at }]
  Auth: JWT / PAT

GET /api/v1/teams-for-create-task
  Response: [{ id, name, prefix }]
  Auth: JWT / PAT

GET /api/v1/teams/{prefix}
  Response: { id, name, prefix, owner_id, team_permissions: number, created_at }
  Auth: JWT / PAT

PATCH /api/v1/teams/{prefix}
  Body: { name?, prefix? }
  Response: { id, name, prefix, owner_id, created_at }
  Auth: JWT / PAT

DELETE /api/v1/teams/{prefix}
  Response: 204 No Content
  Auth: JWT / PAT

GET /api/v1/teams/{prefix}/members
  Response: [{ user: { id, email, username }, permissions: number, joined_at }]
  Auth: JWT / PAT

POST /api/v1/teams/{prefix}/members
  Body: { email, permissions: ["edit.title", "team.manage_filters", ...] }
  Response: 201 Created
  Auth: JWT / PAT

PATCH /api/v1/teams/{prefix}/members/{userID}/permissions
  Body: { permissions: ["edit.title", "team.manage_filters", ...] }
  Response: 204 No Content
  Auth: JWT / PAT

DELETE /api/v1/teams/{prefix}/members/{userID}
  Response: 204 No Content
  Auth: JWT / PAT
```

---

## Permissions

```
POST /api/v1/tasks/{key}/permissions
  Body: { user_id, permissions: ["edit.title", "share", ...] }
  Response: { permission_id }
  Auth: JWT / PAT

DELETE /api/v1/tasks/{key}/permissions/{permissionID}
  Response: 204 No Content
  Auth: JWT / PAT
```

Права приходят встроенными в ответы ресурсов:
- `GET /api/v1/tasks/{key}` → `task_permissions: number`
- `GET /api/v1/teams` → `team_permissions: number`
- `GET /api/v1/teams/{prefix}` → `team_permissions: number`

---

## Tasks

```
POST /api/v1/tasks
  Body: { title, description?, owner_type, owner_id, assignee_id?, parent_id?, label_ids? }
  Response: { ...task, task_permissions: number }
  Auth: JWT / PAT

GET /api/v1/tasks/{key}
  Response: { ...task, task_permissions: number }
  Auth: JWT / PAT

PATCH /api/v1/tasks/{key}
  Body: {
    title?,
    description?,
    status?,
    assignee_id?,        — null снимает ответственного
    parent_id?,          — null снимает родителя
    label_ids?,          — полный новый массив; [] удаляет все метки
    links?,              — полный новый массив [{ url, title? }]
    relations?,          — полный новый массив ключей связанных задач
    blocking?            — полный новый массив ключей блокируемых задач; валидируется цикл
  }
  Response: { ...task, task_permissions: number }
  Auth: JWT / PAT

DELETE /api/v1/tasks/{key}
  Response: 204 No Content
  Auth: JWT / PAT

POST /api/v1/tasks/{key}/attachments
  Body: multipart/form-data { file }
  Response: { id, name, url, size, created_at }
  Auth: JWT / PAT

DELETE /api/v1/tasks/{key}/attachments/{attachmentID}
  Response: 204 No Content
  Auth: JWT / PAT
```

---

## Labels

```
GET /api/v1/labels?owner_type=&owner_id=
  Response: [{ id, name, color }]
  Auth: JWT / PAT

POST /api/v1/labels
  Body: { owner_type, owner_id, name, color }
  Response: { id, name, color }
  Auth: JWT / PAT

PATCH /api/v1/labels/{labelID}
  Body: { name?, color? }
  Response: { id, name, color }
  Auth: JWT / PAT

DELETE /api/v1/labels/{labelID}
  Response: 204 No Content
  Auth: JWT / PAT
```

---

## Filters

```
GET /api/v1/filters
  Response: {
    personal_filters: [{ id, name }],
    teams: [{
      id, name, prefix, team_permissions: number,
      filters: [{ id, name }]
    }]
  }
  Auth: JWT / PAT

POST /api/v1/filters
  Body: { name, owner_type, owner_id, filter_mode,
          — simple mode:
          search_contains?, assignee_ids?, status?, label_ids?,
          — rsql mode:
          rsql?,
          — только simple mode, только для личных фильтров:
          team_id? }
  Response: { id, name, owner_type, owner_id, filter_mode, ...criteria fields }
  Auth: JWT / PAT

GET /api/v1/filters/{filterID}
  Response: { id, name, owner_type, owner_id, ...criteria fields }
  Auth: JWT / PAT

PATCH /api/v1/filters/{filterID}
  Body: { name?, filter_mode?,
          — simple mode:
          search_contains?, assignee_ids?, status?, label_ids?,
          — rsql mode:
          rsql?,
          — только simple mode, только для личных фильтров:
          team_id? }
  Response: { id, name, owner_type, owner_id, filter_mode, ...criteria fields }
  Auth: JWT / PAT

DELETE /api/v1/filters/{filterID}
  Response: 204 No Content
  Auth: JWT / PAT

GET /api/v1/filters/{filterID}/tasks?cursor=
  Response: { items: [{ ...task, task_permissions: number }], next_cursor }
  Auth: JWT / PAT

GET /api/v1/filters/{filterID}/settings
  Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }
  Auth: JWT / PAT

PATCH /api/v1/filters/{filterID}/settings
  Body: { columns?, sort1_column?, sort1_dir?, sort2_column?, sort2_dir? }
  Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }
  Auth: JWT / PAT
```

---

## Search

```
GET /api/v1/search?q=&cursor=
  Response: { items: [{ id, key, title, status, owner_type, owner_id }], next_cursor }
  Auth: JWT / PAT
```

---

## Automations

```
GET /api/v1/automations?owner_type=&owner_id=
  Response: [{ id, name, trigger, enabled }]
  Auth: JWT / PAT

POST /api/v1/automations
  Body: { owner_type, owner_id, name, trigger, criteria?, actions }
  Response: { id, name, trigger, criteria, actions, enabled }
  Auth: JWT / PAT

PATCH /api/v1/automations/{automationID}
  Body: { name?, trigger?, criteria?, actions?, enabled? }
  Response: { id, name, trigger, criteria, actions, enabled }
  Auth: JWT / PAT

DELETE /api/v1/automations/{automationID}
  Response: 204 No Content
  Auth: JWT / PAT
```

---

## Import

```
POST /api/v1/import/csv
  Content-Type: multipart/form-data
  Body: { file, team_id?, id_column }
  Response: { import_id }
  Auth: JWT / PAT

POST /api/v1/import/jira
  Body: { instance_url, project_key, api_token, email, team_id? }
  Response: { import_id }
  Auth: JWT / PAT

POST /api/v1/import/github
  Body: { owner, repo, access_token, team_id? }
  Response: { import_id }
  Auth: JWT / PAT

GET /api/v1/import/{importID}
  Response: { id, source, status, progress: { total, processed },
              result: { imported, updated, errors: [{ origin_id, error }] } }
  Auth: JWT / PAT

WebSocket /api/v1/ws/import/{importID}?token=
  Messages: { processed, total }
  Final message: { status: "completed", imported, updated, errors: [{ origin_id, error }] }
```

---

## Subscriptions

```
GET /api/v1/subscription
  Response: { subject_type: "user", subject_id, plan, status }
  Auth: JWT / PAT

GET /api/v1/teams/{prefix}/subscription
  Response: { subject_type: "team", subject_id, plan, status }
  Auth: JWT / PAT

PATCH /api/v1/teams/{prefix}/subscription
  Body: { plan }
  Response: { subject_type: "team", subject_id, plan, status }
  Auth: JWT / PAT
  Requires: team.manage_owners
```

---

## Private Access Token

```
POST /api/v1/pat
  Response: { token }
  Auth: JWT

DELETE /api/v1/pat
  Response: 204 No Content
  Auth: JWT
```

---

## Activity Feed

```
GET /api/v1/feed?mode=personal&cursor=
  Response: { items: [{ id, event_type, actor: { id, username }, resource_type,
              resource_id, payload, created_at }], next_cursor }
  Auth: JWT / PAT

GET /api/v1/feed?mode=team&team_prefix=&actor_ids[]=&cursor=
  Response: { items: [...], next_cursor }
  Auth: JWT / PAT

GET /api/v1/filters/{filterID}/feed?cursor=
  Response: { items: [...], next_cursor }
  Auth: JWT / PAT

WebSocket /api/v1/ws/personal?token=
  Message: { event: "feed_updated" }

WebSocket /api/v1/ws/tasks/{key}?token=
  Message: { event: "task_updated" }

WebSocket /api/v1/ws/teams/{prefix}?token=
  Message: { event: "feed_updated" }

WebSocket /api/v1/ws/filters/{filterID}?token=
  Message: { event: "feed_updated" }

WebSocket /api/v1/ws/filter-list-changes?token=
  Message: { event: "filter_list_updated" }
```

---

## User Settings

```
GET /api/v1/settings
  Response: { language, theme, sidebar_personal_open, sidebar_teams_open,
              team_states: [{ team_id, open }] }
  Auth: JWT / PAT

PATCH /api/v1/settings
  Body: { language?, theme?, sidebar_personal_open?, sidebar_teams_open?,
          team_states?: [{ team_id, open }] }
  Response: { language, theme, sidebar_personal_open, sidebar_teams_open,
              team_states: [{ team_id, open }] }
  Auth: JWT / PAT
```
