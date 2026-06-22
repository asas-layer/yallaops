-- name: GetEnvironment :one
SELECT id, release_id, env_name, status, deployed_at
FROM release_environments
WHERE release_id = $1 AND env_name = $2;

-- name: ListEnvironmentsByRelease :many
SELECT id, release_id, env_name, status, deployed_at
FROM release_environments
WHERE release_id = $1
ORDER BY env_name;

-- name: UpsertEnvironment :one
INSERT INTO release_environments (release_id, env_name, status)
VALUES ($1, $2, $3)
ON CONFLICT (release_id, env_name)
DO UPDATE SET status = EXCLUDED.status, deployed_at = CASE
    WHEN EXCLUDED.status = 'deployed' THEN now()
    ELSE release_environments.deployed_at
END
RETURNING id, release_id, env_name, status, deployed_at;
