-- name: SelectPostsByLocation :many
SELECT p.*
FROM posts p
JOIN users u ON p.user_id = u.id
WHERE ST_DWITHIN(
    $1,
    u.location,
    $2 * 1609.34
);
