-- name: InsertAuditEntry :one
INSERT INTO audit_log (release_id, actor, action, detail)
VALUES ($1, $2, $3, $4)
RETURNING id, release_id, actor, action, detail, created_at;

-- name: ListAuditEntriesByRelease :many
SELECT id, release_id, actor, action, detail, created_at
FROM audit_log
WHERE release_id = $1
ORDER BY created_at DESC;
