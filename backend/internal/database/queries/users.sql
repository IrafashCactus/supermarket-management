-- name: GetUserByID :one
SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT id, name, email, password_hash, role, is_active, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (
    name, email, password_hash, role
) VALUES (
    $1, $2, $3, $4
)
RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at;

-- name: UpdateUser :one
UPDATE users
SET name = $2,
    email = $3,
    role = $4,
    is_active = $5,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at;

-- name: UpdateUserPassword :one
UPDATE users
SET password_hash = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
