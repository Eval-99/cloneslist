-- name: PostCategoryByID :one
SELECT *
FROM categories
WHERE post_id =
$1;
