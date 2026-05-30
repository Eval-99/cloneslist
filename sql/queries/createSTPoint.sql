-- name: CreateSTPoint :one
SELECT ST_Point($1, $2, 4326);
