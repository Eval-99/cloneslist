-- name: PostsByCategory :many
SELECT *
FROM categories
JOIN posts ON categories.post_id = posts.id
WHERE name =
$1;
