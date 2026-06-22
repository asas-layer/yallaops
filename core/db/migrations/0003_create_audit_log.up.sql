CREATE TABLE audit_log (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    release_id TEXT        NOT NULL REFERENCES releases (id) ON DELETE CASCADE,
    actor      TEXT        NOT NULL,
    action     TEXT        NOT NULL,
    detail     JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_log_release_id ON audit_log (release_id);
CREATE INDEX idx_audit_log_created_at ON audit_log (created_at DESC);
