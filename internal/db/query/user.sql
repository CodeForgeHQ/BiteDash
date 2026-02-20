-- name: CreateUser :exec
INSERT INTO users (
    id,
    email,
    password,
    created_at
) VALUES (
    $1, $2, $3, $4
);

-- name: GetUserByEmail :one
SELECT id, email, password, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password, created_at
FROM users
WHERE id = $1;