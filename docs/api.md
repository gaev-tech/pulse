# API

Base URL: `/api/v1`

Два способа авторизации — см. [requirements.md](requirements.md#api-authentication).

Cursor-based пагинация: `?cursor=` в запросе, `next_cursor` в ответе (null если страниц больше нет).

---

## Auth

```
POST /api/v1/auth/magic-link
  Body: { email }
  Response: 200 OK
  Auth: нет
  — отправляет письмо со ссылкой; если пользователя нет — создаёт его

POST /api/v1/auth/magic-link/verify
  Body: { token }
  Response: { access_token, refresh_token, user: { id, email, username } }
  Auth: нет
  — верифицирует токен из ссылки, возвращает JWT-пару

POST /api/v1/auth/refresh
  Body: { refresh_token }
  Response: { access_token, refresh_token }
  Auth: нет

POST /api/v1/auth/logout
  Body: { refresh_token }
  Response: 204 No Content
  Auth: нет
  — отзывает конкретный refresh token (выход из одной сессии)

GET /api/v1/auth/me
  Response: { id, email, username }
  Auth: JWT / PAT
  — возвращает данные текущего пользователя по токену
```

---

## Users

```
GET /api/v1/teams/{prefix}/members/search?email=
  Response: [{ id, email, username }]
  Auth: JWT / PAT
  — поиск для приглашения в команду, исключает уже участников

GET /api/v1/tasks/{key}/share/search?email=
  Response: [{ id, email, username }]
  Auth: JWT / PAT
  — поиск для шаринга задачи, исключает тех, у кого уже есть доступ
```

---

## Teams

```
POST /api/v1/teams
  Body: { name, prefix }
  Response: { id, name, prefix, owner_id, created_at }
  Auth: JWT / PAT
  — Team Create Modal

GET /api/v1/teams
  Response: [{ id, name, prefix, owner_id, team_permissions: number, created_at }]
  Auth: JWT / PAT
  — команды, в которых состоит текущий пользователь

GET /api/v1/teams-for-create-task
  Response: [{ id, name, prefix }]
  Auth: JWT / PAT
  — команды, в которых пользователь имеет право `task.create`;
    используется в Task Create Modal

GET /api/v1/teams/{prefix}
  Response: { id, name, prefix, owner_id, team_permissions: number, created_at }
  Auth: JWT / PAT
  — Team Activity Screen, Team Edit Modal

PATCH /api/v1/teams/{prefix}
  Body: { name?, prefix? }
  Response: { id, name, prefix, owner_id, created_at }
  Auth: JWT / PAT
  — Team Edit Modal: изменение названия и prefix

DELETE /api/v1/teams/{prefix}
  Response: 204 No Content
  Auth: JWT / PAT
  — Team Edit Modal: удаление команды

GET /api/v1/teams/{prefix}/members
  Response: [{ user: { id, email, username }, permissions: number, joined_at }]
  Auth: JWT / PAT
  — Team Edit Modal: список участников с их правами

POST /api/v1/teams/{prefix}/members
  Body: { email, permissions: ["edit.title", "team.manage_filters", ...] }
  Response: 201 Created
  Auth: JWT / PAT
  — Team Edit Modal: приглашение пользователя; право `view` выдаётся автоматически;
    permissions может быть пустым массивом

PATCH /api/v1/teams/{prefix}/members/{userID}/permissions
  Body: { permissions: ["edit.title", "team.manage_filters", ...] }
  Response: 204 No Content
  Auth: JWT / PAT
  — Team Edit Modal: замена прав участника; право `view` выдаётся автоматически;
    permissions может быть пустым массивом

DELETE /api/v1/teams/{prefix}/members/{userID}
  Response: 204 No Content
  Auth: JWT / PAT
  — Team Edit Modal: удаление участника, каскадный отзыв его прав
    и прав, выданных им другим пользователям по командным задачам
```

---

## Permissions

```
POST /api/v1/tasks/{key}/permissions
  Body: { user_id, permissions: ["edit.title", "share", ...] }
  Response: { permission_id }
  Auth: JWT / PAT
  — шаринг задачи: выдача прав пользователю; право `view` выдаётся автоматически;
    permissions может быть пустым массивом

DELETE /api/v1/tasks/{key}/permissions/{permissionID}
  Response: 204 No Content
  Auth: JWT / PAT
  — отзыв права; каскадно удаляет права, выданные через эту цепочку шаринга
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
  — Task Create Modal; key_number генерируется атомарно на бэкенде

GET /api/v1/tasks/{key}
  Response: { ...task, task_permissions: number }
  Auth: JWT / PAT
  — Task Screen, Task Modal; фиксирует факт открытия в task_opens

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
  — Task Screen, Task Modal

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
  — личные метки или метки команды; используется в Task Screen/Modal
    и Automation Modal

POST /api/v1/labels
  Body: { owner_type, owner_id, name, color }
  Response: { id, name, color }
  Auth: JWT / PAT
  — модалка личных меток, Team Edit Modal

PATCH /api/v1/labels/{labelID}
  Body: { name?, color? }
  Response: { id, name, color }
  Auth: JWT / PAT
  — переименование фиксирует событие во всех задачах с этой меткой

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
  — Navigation Sidebar: личные фильтры + команды с их фильтрами

POST /api/v1/filters
  Body: { name, owner_type, owner_id, author_id?, assignee_id?, title_contains?,
          key_contains?, desc_contains?, status?, label_ids?,
          — null = не фильтровать; [] = задачи без меток; [id,...] = задачи с любой из меток
          updated_after?,
          updated_before?, created_after?, created_before?, has_links?,
          has_relation?, relation_task_key?, has_blocking?, blocking_task_key?,
          parent_task_key?, attachment_name?, team_id? }
  Response: { id, name, owner_type, owner_id, ...criteria fields }
  Auth: JWT / PAT
  — Filter Create Modal

GET /api/v1/filters/{filterID}
  Response: { id, name, owner_type, owner_id, ...criteria fields }
  Auth: JWT / PAT
  — Filter Results Screen: загрузка критериев при открытии фильтра

PATCH /api/v1/filters/{filterID}
  Body: { name?, author_id?, assignee_id?, title_contains?, key_contains?,
          desc_contains?, status?, label_ids?,
          — null = не фильтровать; [] = задачи без меток; [id,...] = задачи с любой из меток
          updated_after?, updated_before?,
          created_after?, created_before?, has_links?, has_relation?, relation_task_key?,
          has_blocking?, blocking_task_key?, parent_task_key?,
          attachment_name?, team_id? }
  Response: { id, name, owner_type, owner_id, ...criteria fields }
  Auth: JWT / PAT
  — Filter Results Screen: сохранение изменённых критериев

DELETE /api/v1/filters/{filterID}
  Response: 204 No Content
  Auth: JWT / PAT
  — Filter Popup

GET /api/v1/filters/{filterID}/tasks?cursor=
  Response: { items: [{ ...task, task_permissions: number }], next_cursor }
  Auth: JWT / PAT
  — Filter Results Screen: результаты фильтра

GET /api/v1/filters/{filterID}/settings
  Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }
  Auth: JWT / PAT

PATCH /api/v1/filters/{filterID}/settings
  Body: { columns?, sort1_column?, sort1_dir?, sort2_column?, sort2_dir? }
  Response: { columns, sort1_column, sort1_dir, sort2_column, sort2_dir }
  Auth: JWT / PAT
  — Table Settings Popup
```

---

## Search

```
GET /api/v1/search?q=&cursor=
  Response: { items: [{ id, key, title, status, owner_type, owner_id }], next_cursor }
  Auth: JWT / PAT
  — Search Modal: поиск по ключу, названию и описанию (contains);
    сортировка: сначала точное совпадение ключа или названия,
    затем по дате последнего открытия задачи пользователем
```

---

## Automations

```
GET /api/v1/automations?owner_type=&owner_id=
  Response: [{ id, name, trigger, enabled }]
  Auth: JWT / PAT
  — Automation Modal: список автоматизаций

POST /api/v1/automations
  Body: { owner_type, owner_id, name, trigger, criteria?, actions }
  Response: { id, name, trigger, criteria, actions, enabled }
  Auth: JWT / PAT
  — Automation Modal: создание

PATCH /api/v1/automations/{automationID}
  Body: { name?, trigger?, criteria?, actions?, enabled? }
  Response: { id, name, trigger, criteria, actions, enabled }
  Auth: JWT / PAT
  — Automation Modal: редактирование и включение/выключение

DELETE /api/v1/automations/{automationID}
  Response: 204 No Content
  Auth: JWT / PAT
```

---

## Private Access Token

```
POST /api/v1/pat
  Response: { token }
  Auth: JWT
  — генерация нового PAT; инвалидирует предыдущий;
    токен возвращается единожды; только JWT

DELETE /api/v1/pat
  Response: 204 No Content
  Auth: JWT
  — отзыв текущего PAT
```

---

## Activity Feed

```
GET /api/v1/feed?mode=personal&cursor=
  Response: { items: [{ id, event_type, actor: { id, username }, resource_type,
              resource_id, payload, created_at }], next_cursor }
  Auth: JWT / PAT
  — Personal Activity Screen

GET /api/v1/feed?mode=team&team_prefix=&actor_ids[]=&cursor=
  Response: { items: [...], next_cursor }
  Auth: JWT / PAT
  — Team Activity Screen; actor_ids — опциональный фильтр по авторам изменений

GET /api/v1/filters/{filterID}/feed?cursor=
  Response: { items: [...], next_cursor }
  Auth: JWT / PAT
  — Filter Activity Screen

WebSocket /api/v1/ws/personal?token=
  — сигнал { event: "feed_updated" } при появлении новых личных событий

WebSocket /api/v1/ws/tasks/{key}?token=
  — сигнал { event: "task_updated" } при изменении конкретной задачи

WebSocket /api/v1/ws/teams/{prefix}?token=
  — сигнал { event: "feed_updated" } при появлении новых событий команды

WebSocket /api/v1/ws/filters/{filterID}?token=
  — сигнал { event: "feed_updated" } при появлении новых событий по фильтру

WebSocket /api/v1/ws/filter-list-changes?token=
  — сигнал { event: "filter_list_updated" } при изменении списка фильтров
    или команд пользователя; клиент делает re-fetch GET /api/v1/filters
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
