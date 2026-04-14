# Database Schema

PostgreSQL 16. Все первичные ключи — UUID, генерируются через `gen_random_uuid()`.

---

## users

```sql
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL UNIQUE,
    username    TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## refresh_tokens

```sql
CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON refresh_tokens (user_id);
```

---

## private_access_tokens

```sql
CREATE TABLE private_access_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON private_access_tokens (user_id);
```

---

## teams

```sql
CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    prefix      TEXT NOT NULL UNIQUE,
    owner_id    UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## team_members

```sql
CREATE TABLE team_members (
    team_id    UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (team_id, user_id)
);

CREATE INDEX ON team_members (user_id);
```

---

## labels

```sql
CREATE TABLE labels (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type  TEXT NOT NULL,
    owner_id    UUID NOT NULL,
    name        TEXT NOT NULL,
    color       TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (owner_type, owner_id, name)
);

CREATE INDEX ON labels (owner_type, owner_id);
```

---

## task_sequences

```sql
CREATE TABLE task_sequences (
    owner_type   TEXT   NOT NULL,
    owner_id     UUID   NOT NULL,
    last_number  BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (owner_type, owner_id)
);
```

Инкремент атомарный: `UPDATE task_sequences SET last_number = last_number + 1 WHERE owner_type = $1 AND owner_id = $2 RETURNING last_number` внутри транзакции создания задачи.

---

## tasks

```sql
CREATE TABLE tasks (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    key_number   BIGINT NOT NULL,
    owner_type   TEXT   NOT NULL,
    owner_id     UUID   NOT NULL,
    title        TEXT   NOT NULL,
    description  TEXT,
    status       TEXT   NOT NULL DEFAULT 'opened',
    assignee_id  UUID   REFERENCES users(id),
    parent_id    UUID   REFERENCES tasks(id),
    created_by   UUID   NOT NULL REFERENCES users(id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    origin_id    TEXT UNIQUE,
    UNIQUE (owner_type, owner_id, key_number)
);

CREATE INDEX ON tasks (owner_type, owner_id);
CREATE INDEX ON tasks (assignee_id);
CREATE INDEX ON tasks (parent_id);
```

---

## task_labels

```sql
CREATE TABLE task_labels (
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, label_id)
);
```

---

## task_links

```sql
CREATE TABLE task_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    title      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON task_links (task_id);
```

---

## task_relations

```sql
CREATE TABLE task_relations (
    task_id         UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    related_task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, related_task_id),
    CHECK (task_id < related_task_id)
);

CREATE INDEX ON task_relations (related_task_id);
```

`CHECK (task_id < related_task_id)` исключает дубли в обратном порядке. При запросе связей задачи ищем по обоим столбцам.

---

## task_blocking

```sql
CREATE TABLE task_blocking (
    blocker_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

CREATE INDEX ON task_blocking (blocked_id);
```

---

## task_attachments

```sql
CREATE TABLE task_attachments (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id     UUID   NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    name        TEXT   NOT NULL,
    url         TEXT   NOT NULL,
    size        BIGINT NOT NULL,
    created_by  UUID   NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON task_attachments (task_id);
```

---

## permissions

```sql
CREATE TABLE permissions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_type  TEXT NOT NULL,
    subject_id    UUID NOT NULL,
    resource_type TEXT NOT NULL,
    resource_id   UUID,
    action        TEXT NOT NULL,
    granted_by    UUID REFERENCES users(id),
    level         TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON permissions (subject_type, subject_id, resource_type, resource_id, action);
```

- `subject_type` — `user` / `team`
- `resource_type` — `task` / `team` / `filter`
- `resource_id` — NULL для scope-level прав (на все ресурсы команды)
- `level` — `owner` / `team` / `direct`
- `granted_by` — NULL для автоматически выданных прав (владелец при создании)

---

## filters

```sql
CREATE TABLE filters (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type       TEXT NOT NULL,
    owner_id         UUID NOT NULL,
    name             TEXT NOT NULL,
    filter_mode      TEXT NOT NULL DEFAULT 'simple',
    search_contains  TEXT,
    assignee_ids     UUID[],
    status           TEXT,
    label_ids        UUID[],
    rsql             TEXT,
    team_id          UUID REFERENCES teams(id),
    created_by       UUID NOT NULL REFERENCES users(id),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON filters (owner_type, owner_id);
```

- `filter_mode` — `simple` или `rsql`; определяет какие поля критериев применяются
- `search_contains` — поиск по ключу, названию и описанию (contains); только в simple mode
- `assignee_ids`, `status`, `label_ids` — только в simple mode
- `rsql` — RSQL-выражение; только в rsql mode
- `team_id` — только для личных фильтров, только в simple mode (AND к основным критериям); в rsql mode выражается внутри RSQL-выражения

---

## task_opens

```sql
CREATE TABLE task_opens (
    task_id   UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, user_id)
);

CREATE INDEX ON task_opens (user_id, opened_at DESC);
```

При открытии задачи: `INSERT ... ON CONFLICT (task_id, user_id) DO UPDATE SET opened_at = now()`.

---

## filter_settings

```sql
CREATE TABLE filter_settings (
    filter_id     UUID PRIMARY KEY REFERENCES filters(id) ON DELETE CASCADE,
    columns       TEXT[] NOT NULL DEFAULT '{}',
    sort1_column  TEXT,
    sort1_dir     TEXT,
    sort2_column  TEXT,
    sort2_dir     TEXT,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `columns` — упорядоченный массив включённых необязательных колонок
- `sort1_column` / `sort1_dir` — первичная сортировка
- `sort2_column` / `sort2_dir` — вторичная сортировка (применяется при равенстве первичной)

---

## events

```sql
CREATE TABLE events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type    TEXT NOT NULL,
    actor_id      UUID NOT NULL REFERENCES users(id),
    resource_type TEXT NOT NULL,
    resource_id   UUID NOT NULL,
    team_id       UUID REFERENCES teams(id),
    payload       JSONB NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON events (actor_id, created_at DESC);
CREATE INDEX ON events (team_id, created_at DESC) WHERE team_id IS NOT NULL;
CREATE INDEX ON events (resource_id, created_at DESC);
CREATE INDEX ON events USING GIN (payload);
```

- Таблица только для вставки, записи не удаляются и не обновляются
- `team_id` — NULL для личных событий
- `payload` — `{"old": ..., "new": ...}` с изменившимися полями; структура зависит от `event_type`

---

## automations

```sql
CREATE TABLE automations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_type  TEXT    NOT NULL,
    owner_id    UUID    NOT NULL,
    name        TEXT    NOT NULL,
    trigger     TEXT    NOT NULL,
    criteria    JSONB   NOT NULL DEFAULT '{}',
    actions     JSONB   NOT NULL DEFAULT '[]',
    enabled     BOOLEAN NOT NULL DEFAULT true,
    created_by  UUID    NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON automations (owner_type, owner_id);
```

- `trigger` — тип события из `events.event_type`, например `task.status_changed`
- `criteria` — условия фильтрации по полям payload, например `{"new_status": "closed"}`
- `actions` — массив вызовов API: `[{"url": "...", "method": "POST", "headers": {...}, "body": "..."}]`

---

## user_settings

```sql
CREATE TABLE user_settings (
    user_id                UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    language               TEXT    NOT NULL DEFAULT 'en',
    theme                  TEXT    NOT NULL DEFAULT 'system',
    sidebar_personal_open  BOOLEAN NOT NULL DEFAULT true,
    sidebar_teams_open     BOOLEAN NOT NULL DEFAULT true,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

- `language` — `en` / `ru`
- `theme` — `light` / `dark` / `system`

---

## user_sidebar_team_states

```sql
CREATE TABLE user_sidebar_team_states (
    user_id  UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id  UUID    NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    open     BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, team_id)
);
```

Запись создаётся при первом изменении состояния. Если записи нет — считается `open = true`.

---

## imports

```sql
CREATE TABLE imports (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id),
    team_id             UUID REFERENCES teams(id),
    source              TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'in_progress',
    progress_total      INT  NOT NULL DEFAULT 0,
    progress_processed  INT  NOT NULL DEFAULT 0,
    result              JSONB,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON imports (user_id);
```

- `source` — `csv`, `jira`, `github`
- `status` — `in_progress`, `completed`
- `result` — `{"imported": N, "updated": N, "errors": [{"origin_id": "...", "error": "..."}]}`
