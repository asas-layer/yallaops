CREATE TYPE env_status AS ENUM ('pending', 'deploying', 'deployed', 'failed', 'blocked');

CREATE TABLE release_environments (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    release_id  TEXT        NOT NULL REFERENCES releases (id) ON DELETE CASCADE,
    env_name    TEXT        NOT NULL,
    status      env_status  NOT NULL DEFAULT 'pending',
    deployed_at TIMESTAMPTZ,
    UNIQUE (release_id, env_name)
);

CREATE INDEX idx_release_environments_release_id ON release_environments (release_id);
