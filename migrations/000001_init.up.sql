CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE SCHEMA IF NOT EXISTS taskana;

-- =========================================================================
-- DATA TYPES (ENUMS)
-- =========================================================================
CREATE TYPE taskana.project_status AS ENUM ('active', 'completed', 'dropped');
CREATE TYPE taskana.task_status AS ENUM ('open', 'completed', 'canceled');
CREATE TYPE taskana.task_bucket AS ENUM ('inbox', 'today', 'anytime', 'someday');
CREATE TYPE taskana.recurrence_type AS ENUM ('fixed', 'from_completion');
CREATE TYPE taskana.target_bucket AS ENUM ('today', 'inbox');

-- =========================================================================
-- TABLE: USERS
-- =========================================================================
CREATE TABLE taskana.users
(
    id            UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version       BIGINT       NOT NULL DEFAULT 1,
    first_name    VARCHAR(100) NOT NULL CHECK (char_length(first_name) BETWEEN 3 AND 100),
    last_name     VARCHAR(100) NOT NULL CHECK (char_length(last_name) BETWEEN 3 AND 100),
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    timezone      VARCHAR(64)  NOT NULL DEFAULT 'UTC',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- =========================================================================
-- TABLE: REFRESH TOKENS
-- =========================================================================
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

-- =========================================================================
-- TABLE: AREAS
-- =========================================================================
CREATE TABLE taskana.areas
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version    BIGINT       NOT NULL DEFAULT 1,
    user_id    UUID         NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    position   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_areas_user_position ON taskana.areas (user_id, position);

-- =========================================================================
-- TABLE: PROJECTS
-- =========================================================================
CREATE TABLE taskana.projects
(
    id           UUID PRIMARY KEY                DEFAULT gen_random_uuid(),
    version      BIGINT                 NOT NULL DEFAULT 1,
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
CREATE INDEX idx_projects_area_position ON taskana.projects (area_id, position);

-- =========================================================================
-- TABLE: HEADINGS
-- =========================================================================
CREATE TABLE taskana.headings
(
    id         UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version    BIGINT       NOT NULL DEFAULT 1,
    user_id    UUID         NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    project_id UUID         NOT NULL REFERENCES taskana.projects (id) ON DELETE CASCADE,
    title      VARCHAR(255) NOT NULL,
    position   INTEGER      NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_headings_project_id_id UNIQUE (project_id, id)
);

CREATE INDEX idx_headings_user_id ON taskana.headings (user_id);
CREATE INDEX idx_headings_project_position ON taskana.headings (project_id, position);

-- =========================================================================
-- TABLE: PERIODIC TASK TEMPLATES
-- =========================================================================
CREATE TABLE taskana.task_templates
(
    id                  UUID PRIMARY KEY                 DEFAULT gen_random_uuid(),
    version             BIGINT                  NOT NULL DEFAULT 1,
    user_id             UUID                    NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    project_id          UUID REFERENCES taskana.projects (id) ON DELETE CASCADE,
    heading_id          UUID,
    title               VARCHAR(255)            NOT NULL,
    notes               TEXT,
    recurrence_rule     TEXT                    NOT NULL,
    recurrence_type     taskana.recurrence_type NOT NULL,
    target_bucket       taskana.target_bucket   NOT NULL DEFAULT 'inbox',
    next_execution_date DATE                    NOT NULL,
    is_time_tracked     BOOLEAN                 NOT NULL DEFAULT FALSE,
    estimated_pomodoros INTEGER                 NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ             NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_task_templates_heading FOREIGN KEY (project_id, heading_id)
        REFERENCES taskana.headings (project_id, id) ON DELETE SET NULL,

    CONSTRAINT check_template_heading_requires_project CHECK (heading_id IS NULL OR project_id IS NOT NULL)
);

CREATE INDEX idx_task_templates_user_id ON taskana.task_templates (user_id);
CREATE INDEX idx_task_templates_next_execution ON taskana.task_templates (next_execution_date);

-- =========================================================================
-- TABLE: TASKS
-- =========================================================================
CREATE TABLE taskana.tasks
(
    id                  UUID PRIMARY KEY             DEFAULT gen_random_uuid(),
    version             BIGINT              NOT NULL DEFAULT 1,
    user_id             UUID                NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    project_id          UUID REFERENCES taskana.projects (id) ON DELETE CASCADE,
    heading_id          UUID,
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

    CONSTRAINT fk_tasks_heading FOREIGN KEY (project_id, heading_id)
        REFERENCES taskana.headings (project_id, id) ON DELETE SET NULL,

    CONSTRAINT check_task_heading_requires_project CHECK (heading_id IS NULL OR project_id IS NOT NULL),

    CONSTRAINT check_task_status_completed_at CHECK (
        (status = 'open' AND completed_at IS NULL)
            OR
        (status IN ('completed', 'canceled') AND completed_at IS NOT NULL)
        ),

    CONSTRAINT check_task_dates CHECK (deadline IS NULL OR start_date IS NULL OR
                                       deadline >= start_date)
);

CREATE INDEX idx_tasks_user_id ON taskana.tasks (user_id);
CREATE INDEX idx_tasks_template_id ON taskana.tasks (template_id);
CREATE INDEX idx_tasks_start_date ON taskana.tasks (start_date);
CREATE INDEX idx_tasks_project_position ON taskana.tasks (project_id, position);
CREATE INDEX idx_tasks_heading_position ON taskana.tasks (heading_id, position);
CREATE INDEX idx_tasks_bucket_position ON taskana.tasks (user_id, bucket, position);

-- =========================================================================
-- TABLE: CHECKLIST ITEMS
-- =========================================================================
CREATE TABLE taskana.checklists
(
    id           UUID PRIMARY KEY      DEFAULT gen_random_uuid(),
    version      BIGINT       NOT NULL DEFAULT 1,
    user_id      UUID         NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    task_id      UUID         NOT NULL REFERENCES taskana.tasks (id) ON DELETE CASCADE,
    title        VARCHAR(255) NOT NULL,
    is_completed BOOLEAN      NOT NULL DEFAULT FALSE,
    position     INTEGER      NOT NULL DEFAULT 0,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_checklists_user_id ON taskana.checklists (user_id);
CREATE INDEX idx_checklists_task_position ON taskana.checklists (task_id,
                                                                position);

-- =========================================================================
-- TABLE: POMODORO SESSION
-- =========================================================================
CREATE TABLE taskana.pomodoro_sessions
(
    id               UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    version          BIGINT      NOT NULL DEFAULT 1,
    user_id          UUID        NOT NULL REFERENCES taskana.users (id) ON DELETE CASCADE,
    task_id          UUID        NOT NULL REFERENCES taskana.tasks (id) ON DELETE CASCADE,
    start_time       TIMESTAMPTZ NOT NULL,
    end_time         TIMESTAMPTZ NOT NULL,
    duration_seconds INTEGER     NOT NULL CHECK (duration_seconds >= 0),
    is_interrupted   BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT check_session_times CHECK (end_time >= start_time)
);

CREATE INDEX idx_pomodoro_sessions_user_id ON taskana.pomodoro_sessions (user_id);
CREATE INDEX idx_pomodoro_sessions_task_id ON taskana.pomodoro_sessions (task_id);
CREATE INDEX idx_pomodoro_sessions_start_time ON taskana.pomodoro_sessions (user_id, start_time);