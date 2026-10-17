-- name: GetProductByID :one
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE p.id = $1 LIMIT 1;

-- name: GetProductByBarcode :one
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE p.barcode = $1 LIMIT 1;

-- name: ListProducts :many
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE (sqlc.narg('search')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('search')::text || '%' OR p.barcode ILIKE '%' || sqlc.narg('search')::text || '%')
  AND (sqlc.narg('category_id')::uuid IS NULL OR p.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('is_active')::boolean IS NULL OR p.is_active = sqlc.narg('is_active')::boolean)
ORDER BY p.name ASC
LIMIT $1 OFFSET $2;

-- name: CountProducts :one
SELECT COUNT(*)
FROM products p
WHERE (sqlc.narg('search')::text IS NULL OR p.name ILIKE '%' || sqlc.narg('search')::text || '%' OR p.barcode ILIKE '%' || sqlc.narg('search')::text || '%')
  AND (sqlc.narg('category_id')::uuid IS NULL OR p.category_id = sqlc.narg('category_id')::uuid)
  AND (sqlc.narg('is_active')::boolean IS NULL OR p.is_active = sqlc.narg('is_active')::boolean);

-- name: CreateProduct :one
INSERT INTO products (
    name, barcode, sku, category_id, unit_price, cost_price,
    tax_rate, unit_of_measure, reorder_level, image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, name, barcode, sku, category_id, unit_price, cost_price,
          tax_rate, unit_of_measure, reorder_level, image_url, is_active,
          created_at, updated_at;

-- name: UpdateProduct :one
UPDATE products
SET name = $2,
    barcode = $3,
    sku = $4,
    category_id = $5,
    unit_price = $6,
    cost_price = $7,
    tax_rate = $8,
    unit_of_measure = $9,
    reorder_level = $10,
    image_url = $11,
    is_active = $12,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, barcode, sku, category_id, unit_price, cost_price,
          tax_rate, unit_of_measure, reorder_level, image_url, is_active,
          created_at, updated_at;

-- name: DeactivateProduct :one
UPDATE products
SET is_active = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, barcode, sku, category_id, unit_price, cost_price,
          tax_rate, unit_of_measure, reorder_level, image_url, is_active,
          created_at, updated_at;

-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING id;
