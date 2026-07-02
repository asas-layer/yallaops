-- name: GetRelease :one
SELECT id, service, version, status, created_at, updated_at
FROM releases
WHERE id = $1;

-- name: ListReleases :many
SELECT id, service, version, status, created_at, updated_at
FROM releases
WHERE service = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateRelease :one
INSERT INTO releases (id, service, version, status)
VALUES ($1, $2, $3, 'draft')
RETURNING id, service, version, status, created_at, updated_at;

-- name: UpdateReleaseStatus :one
UPDATE releases
SET status = $2, updated_at = now()
WHERE id = $1
RETURNING id, service, version, status, created_at, updated_at;

-- name: CountReleasesByService :one
SELECT count(*) FROM releases WHERE service = $1;
