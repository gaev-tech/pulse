# Automations

Автоматизации срабатывают на события в системе. Каждая автоматизация состоит из:
- **Триггер** — тип события (`event_type` из таблицы `events`)
- **Критерии** — необязательные условия фильтрации по полям payload триггера
- **Действия** — один или несколько вызовов API (Public API системы или сторонний API) с доступом к полям payload

---

## Task triggers

| Триггер | Когда срабатывает |
|---|---|
| `task.created` | Создана задача |
| `task.deleted` | Удалена задача |
| `task.title_changed` | Изменено название |
| `task.description_changed` | Изменено описание |
| `task.status_changed` | Изменён статус |
| `task.assignee_changed` | Изменён ответственный |
| `task.labels_changed` | Изменены метки |
| `task.links_changed` | Изменены ссылки |
| `task.relations_changed` | Изменены связи |
| `task.blocking_changed` | Изменено блокирование |
| `task.parent_changed` | Изменена родительская задача |
| `task.attachments_changed` | Изменены вложения |
| `task.shared` | Выдан доступ к задаче |

Payload для всех task-триггеров:

```json
{
  "task": { "id", "key", "title", "status", "owner_type", "owner_id" },
  "actor": { "id", "email", "username" },
  "old_<field>": "...",
  "new_<field>": "..."
}
```

Для `task.created` и `task.deleted` — только `task` и `actor`, без `old_`/`new_`.

Примеры:

```json
// task.status_changed
{
  "task": { "id": "...", "key": "ACME-42", "title": "...", "status": "closed", "owner_type": "team", "owner_id": "..." },
  "actor": { "id": "...", "email": "...", "username": "..." },
  "old_status": "opened",
  "new_status": "closed"
}
```

---

## Team triggers

| Триггер | Когда срабатывает |
|---|---|
| `team.member_added` | Добавлен участник команды |
| `team.member_removed` | Удалён участник команды |
| `team.member_permissions_changed` | Изменены права участника |

Payload:

```json
// team.member_added
{
  "team": { "id", "name", "prefix" },
  "actor": { "id", "email", "username" },
  "user": { "id", "email", "username" },
  "permissions": ["edit.title", "team.manage_filters"]
}

// team.member_removed
{
  "team": { "id", "name", "prefix" },
  "actor": { "id", "email", "username" },
  "user": { "id", "email", "username" }
}

// team.member_permissions_changed
{
  "team": { "id", "name", "prefix" },
  "actor": { "id", "email", "username" },
  "user": { "id", "email", "username" },
  "old_permissions": ["edit.title"],
  "new_permissions": ["edit.title", "team.manage_filters"]
}
```

---

## Filter triggers

| Триггер | Когда срабатывает |
|---|---|
| `filter.created` | Создан фильтр |
| `filter.updated` | Изменены критерии или название фильтра |
| `filter.deleted` | Удалён фильтр |
| `filter.task_entered` | Задача вошла в результаты фильтра вследствие изменения задачи |
| `filter.task_left` | Задача вышла из результатов фильтра вследствие изменения задачи |

Payload:

```json
// filter.created, filter.deleted
{
  "filter": { "id", "name", "owner_type", "owner_id" },
  "actor": { "id", "email", "username" }
}

// filter.updated
{
  "filter": { "id", "name", "owner_type", "owner_id" },
  "actor": { "id", "email", "username" },
  "old_criteria": { "..." },
  "new_criteria": { "..." }
}

// filter.task_entered, filter.task_left
{
  "filter": { "id", "name", "owner_type", "owner_id" },
  "task": { "id", "key", "title", "status", "owner_type", "owner_id" }
}
```

---

## Label triggers

| Триггер | Когда срабатывает |
|---|---|
| `label.renamed` | Переименована метка |
| `label.deleted` | Удалена метка |

Payload:

```json
// label.renamed
{
  "label": { "id", "owner_type", "owner_id" },
  "actor": { "id", "email", "username" },
  "old_name": "...",
  "new_name": "..."
}

// label.deleted
{
  "label": { "id", "owner_type", "owner_id" },
  "actor": { "id", "email", "username" },
  "name": "..."
}
```
