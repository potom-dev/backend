-- name: CreateGroup :one
INSERT INTO groups (id, created_at, updated_at, name, author_id)
VALUES (
    gen_random_uuid(),
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    $1,
    $2
)
RETURNING *;
