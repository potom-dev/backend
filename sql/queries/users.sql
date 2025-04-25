-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users
ORDER BY created_at ASC;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;