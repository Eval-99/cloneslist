-- name: AddToCategory :one
INSERT INTO categories (name, post_id)
VALUES (
    $1,
    $2
)
RETURNING *;

