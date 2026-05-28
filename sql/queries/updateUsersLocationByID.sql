-- name: UpdateUsersLocationByID :one
UPDATE users
SET
    location =  ST_Point($2, $3, 4326)
WHERE id = $1
RETURNING *;
