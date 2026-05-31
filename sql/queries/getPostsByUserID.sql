-- name: PostsByUserID :many
SELECT * FROM posts
WHERE user_id =
$1;
