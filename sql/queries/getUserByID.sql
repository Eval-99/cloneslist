-- name: UsersByID :one
SELECT
    *
FROM
    users
WHERE
    id
=
$1;

