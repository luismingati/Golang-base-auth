-- name: InsertUser :one
INSERT INTO users (id, email, password, username)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: FindUserByEmail :one
SELECT id, email, password, username FROM users WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
SET password = $1
WHERE email = $2;