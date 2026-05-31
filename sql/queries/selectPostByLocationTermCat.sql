-- name: SelectPostsByLocationTermCat :many
SELECT p.*
FROM posts p
JOIN users u ON p.user_id = u.id
JOIN categories c ON p.id = c.post_id
WHERE ST_DWITHIN(
    $1,
    u.location,
    $2 * 1609.34
)
AND TO_TSVECTOR(p.title || ' ' || p.description) @@ TO_TSQUERY($3)
AND c.name = $4;
