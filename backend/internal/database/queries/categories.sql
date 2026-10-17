-- name: GetCategoryByID :one
SELECT id, name, description, created_at, updated_at
FROM categories
WHERE id = $1 LIMIT 1;

-- name: ListCategories :many
SELECT id, name, description, created_at, updated_at
FROM categories
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: CreateCategory :one
INSERT INTO categories (
    name, description
) VALUES (
    $1, $2
)
RETURNING id, name, description, created_at, updated_at;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: CountCategories :one
SELECT COUNT(*) FROM categories;

-- name: CountProductsByCategory :one
SELECT COUNT(*) FROM products WHERE category_id = $1;
