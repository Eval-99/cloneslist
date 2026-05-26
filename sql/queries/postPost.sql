-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, user_id, title, description, price, location)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3,
    $4,
    ST_Point($5, $6, 4326)
)
RETURNING *;
