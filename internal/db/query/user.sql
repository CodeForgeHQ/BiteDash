-- name: CreateUser :one
INSERT INTO users (
    ID,
    email,
    password_hash,
    created_at
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, email, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email, password_hash, created_at
FROM users
WHERE id = $1;