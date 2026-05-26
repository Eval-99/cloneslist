-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, user_id, title, description, price, category, location)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    $3,
    $4,
    $5,
    ST_Point($6, $7, 4326)
)
RETURNING *;
