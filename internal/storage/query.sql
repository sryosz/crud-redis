-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: CreateUser :exec
INSERT INTO users(email, password)
VALUES ($1, $2);

-- name: GetUsers :many
SELECT *
FROM users;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;