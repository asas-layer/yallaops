CREATE TYPE release_status AS ENUM ('draft', 'running', 'deployed', 'failed');

CREATE TABLE releases (
    id         TEXT PRIMARY KEY,
    service    TEXT NOT NULL,
    version    TEXT NOT NULL,
    status     release_status NOT NULL DEFAULT 'draft',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_releases_service ON releases (service);
CREATE INDEX idx_releases_status  ON releases (status);
