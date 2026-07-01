CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS taskana;

CREATE TYPE taskana.project_status AS ENUM ('active', 'completed', 'dropped');
CREATE TYPE taskana.task_status AS ENUM ('open', 'completed', 'canceled');
CREATE TYPE taskana.task_bucket AS ENUM ('inbox', 'today', 'anytime', 'someday');
CREATE TYPE taskana.recurrence_type AS ENUM ('fixed', 'from_completion');
CREATE TYPE taskana.target_bucket AS ENUM ('today', 'inbox');

CREATE TABLE taskana.users
(
    id            UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version       bigint       not null default 1,
    first_name    varchar(100) not null check (
        char_length(first_name) between 3 and 100
        ),
    last_name     varchar(100) not null check (
        char_length(last_name) between 3 and 100
        ),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    timezone      VARCHAR(64)  NOT NULL DEFAULT 'UTC',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE taskana.refresh_tokens
(
    id         UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    token_hash TEXT        NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON taskana.refresh_tokens (user_id);

CREATE TABLE taskana.areas
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version       bigint       not null default 1,
    user_id    UUID         NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    position   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_areas_user_id ON taskana.areas (user_id);

CREATE TABLE taskana.projects
(
    id           UUID PRIMARY KEY                DEFAULT gen_random_uuid(),
    user_id      UUID                   NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    area_id      UUID                   REFERENCES taskana.areas (id) ON DELETE SET NULL,
    title        VARCHAR(255)           NOT NULL,
    notes        TEXT,
    status       taskana.project_status NOT NULL DEFAULT 'active',
    position     INTEGER                NOT NULL DEFAULT 0,
    deadline     TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ            NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ            NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_user_id ON taskana.projects (user_id);
CREATE INDEX idx_projects_area_id ON taskana.projects (area_id);

CREATE TABLE taskana.headings
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    project_id UUID         NOT NULL REFERENCES taskana.projects (id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    position   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_headings_project_id ON taskana.headings (project_id);

CREATE TABLE taskana.task_templates
(
    id                  UUID PRIMARY KEY                 DEFAULT gen_random_uuid(),
    user_id             UUID                    NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    project_id          UUID REFERENCES taskana.projects (id) ON DELETE CASCADE,
    heading_id          UUID                    REFERENCES taskana.headings (id) ON DELETE SET NULL,
    title               VARCHAR(255)            NOT NULL,
    notes               TEXT,
    recurrence_rule     TEXT                    NOT NULL,
    recurrence_type     taskana.recurrence_type NOT NULL,
    target_bucket       taskana.target_bucket   NOT NULL DEFAULT 'inbox',
    next_execution_date DATE                    NOT NULL,
    is_time_tracked     BOOLEAN                 NOT NULL DEFAULT FALSE,
    estimated_pomodoros INTEGER                 NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_task_templates_user_id ON taskana.task_templates (user_id);
CREATE INDEX idx_task_templates_next_execution ON taskana.task_templates (next_execution_date)
    WHERE recurrence_type = 'fixed';

CREATE TABLE taskana.tasks
(
    id                  UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    user_id             UUID                NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    project_id          UUID REFERENCES taskana.projects (id) ON DELETE CASCADE,
    heading_id          UUID                REFERENCES taskana.headings (id) ON DELETE SET NULL,
    template_id         UUID                REFERENCES taskana.task_templates (id) ON DELETE SET NULL,
    title               VARCHAR(255)        NOT NULL,
    notes               TEXT,
    status              taskana.task_status NOT NULL DEFAULT 'open',
    bucket              taskana.task_bucket NOT NULL DEFAULT 'inbox',
    start_date          DATE,
    deadline            DATE,
    position            INTEGER             NOT NULL DEFAULT 0,
    is_time_tracked     BOOLEAN             NOT NULL DEFAULT FALSE,
    estimated_pomodoros INTEGER             NOT NULL DEFAULT 0,
    completed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ         NOT NULL DEFAULT NOW(),

    CHECK (
        (status IN ('open') AND completed_at IS NULL)
            OR
        (status IN ('completed', 'canceled') AND completed_at IS NOT NULL)
        )
);

CREATE INDEX idx_tasks_user_id ON taskana.tasks (user_id);
CREATE INDEX idx_tasks_project_id ON taskana.tasks (project_id);
CREATE INDEX idx_tasks_heading_id ON taskana.tasks (heading_id);
CREATE INDEX idx_tasks_template_id ON taskana.tasks (template_id);
CREATE INDEX idx_tasks_start_date ON taskana.tasks (start_date);
CREATE INDEX idx_tasks_bucket ON taskana.tasks (user_id, bucket);

CREATE TABLE taskana.checklist_items
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    task_id      UUID         NOT NULL REFERENCES taskana.tasks (id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    is_completed BOOLEAN      NOT NULL DEFAULT FALSE,
    position     INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_checklist_items_task_id ON taskana.checklist_items (task_id);

CREATE TABLE taskana.pomodoro_sessions
(
    id               UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    user_id          UUID        NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    task_id          UUID        NOT NULL REFERENCES taskana.tasks (id) ON DELETE CASCADE,
    start_time       TIMESTAMPTZ NOT NULL,
    end_time         TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER     NOT NULL CHECK (duration_seconds >= 0),
    is_interrupted   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CHECK (end_time >= start_time)
);

CREATE INDEX idx_pomodoro_sessions_user_id ON taskana.pomodoro_sessions (user_id);
CREATE INDEX idx_pomodoro_sessions_task_id ON taskana.pomodoro_sessions (task_id);
CREATE INDEX idx_pomodoro_sessions_start_time ON taskana.pomodoro_sessions (user_id, start_time);
