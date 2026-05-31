-- name: UpdatePost :one
UPDATE posts
SET
    updated_at = now(),
    title = $2,
    description = $3,
    price = $4,
    status = $5
WHERE id = $1
RETURNING *;
