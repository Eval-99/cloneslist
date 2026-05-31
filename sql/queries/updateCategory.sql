-- name: UpdatePostCategory :exec
UPDATE categories
SET
name =
$2
WHERE post_id = $1
RETURNING *;
