CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT NOT NULL UNIQUE,
    username    TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON refresh_tokens (user_id);

CREATE TABLE private_access_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON private_access_tokens (user_id);

CREATE TABLE teams (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    prefix      TEXT NOT NULL UNIQUE,
    owner_id    UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE team_members (
    team_id    UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (team_id, user_id)
);

CREATE INDEX ON team_members (user_id);

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

CREATE TABLE task_sequences (
    owner_type   TEXT   NOT NULL,
    owner_id     UUID   NOT NULL,
    last_number  BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (owner_type, owner_id)
);

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

CREATE TABLE task_labels (
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    label_id   UUID NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, label_id)
);

CREATE TABLE task_links (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id    UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    url        TEXT NOT NULL,
    title      TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX ON task_links (task_id);

CREATE TABLE task_relations (
    task_id         UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    related_task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, related_task_id),
    CHECK (task_id < related_task_id)
);

CREATE INDEX ON task_relations (related_task_id);

CREATE TABLE task_blocking (
    blocker_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    blocked_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (blocker_id, blocked_id),
    CHECK (blocker_id != blocked_id)
);

CREATE INDEX ON task_blocking (blocked_id);

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

CREATE TABLE task_opens (
    task_id   UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    user_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (task_id, user_id)
);

CREATE INDEX ON task_opens (user_id, opened_at DESC);

CREATE TABLE filter_settings (
    filter_id     UUID PRIMARY KEY REFERENCES filters(id) ON DELETE CASCADE,
    columns       TEXT[] NOT NULL DEFAULT '{}',
    sort1_column  TEXT,
    sort1_dir     TEXT,
    sort2_column  TEXT,
    sort2_dir     TEXT,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

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

CREATE TABLE user_settings (
    user_id                UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    language               TEXT    NOT NULL DEFAULT 'en',
    theme                  TEXT    NOT NULL DEFAULT 'system',
    sidebar_personal_open  BOOLEAN NOT NULL DEFAULT true,
    sidebar_teams_open     BOOLEAN NOT NULL DEFAULT true,
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_sidebar_team_states (
    user_id  UUID    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    team_id  UUID    NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    open     BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, team_id)
);

CREATE TABLE subscriptions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subject_type  TEXT NOT NULL,
    subject_id    UUID NOT NULL,
    plan          TEXT NOT NULL DEFAULT 'free',
    status        TEXT NOT NULL DEFAULT 'active',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (subject_type, subject_id)
);

CREATE INDEX ON subscriptions (subject_type, subject_id);

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
