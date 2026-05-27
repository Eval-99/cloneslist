-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, location)
VALUES (
    gen_random_uuid(),
    now(),
    now(),
    $1,
    $2,
    ST_Point($3, $4, 4326)
)
RETURNING *;
