package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// ProductRow is the result of a product query with a JOIN on categories.
type ProductRow struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	CategoryName  string         `json:"category_name"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

const getProductByID = `-- name: GetProductByID :one
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE p.id = $1 LIMIT 1
`

type GetProductByIDRow struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	CategoryName  string         `json:"category_name"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductByID(ctx context.Context, id uuid.UUID) (GetProductByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getProductByID, id)
	var i GetProductByIDRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Barcode,
		&i.Sku,
		&i.CategoryID,
		&i.CategoryName,
		&i.UnitPrice,
		&i.CostPrice,
		&i.TaxRate,
		&i.UnitOfMeasure,
		&i.ReorderLevel,
		&i.ImageURL,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getProductByBarcode = `-- name: GetProductByBarcode :one
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE p.barcode = $1 LIMIT 1
`

type GetProductByBarcodeRow struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	CategoryName  string         `json:"category_name"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) GetProductByBarcode(ctx context.Context, barcode string) (GetProductByBarcodeRow, error) {
	row := q.db.QueryRowContext(ctx, getProductByBarcode, barcode)
	var i GetProductByBarcodeRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Barcode,
		&i.Sku,
		&i.CategoryID,
		&i.CategoryName,
		&i.UnitPrice,
		&i.CostPrice,
		&i.TaxRate,
		&i.UnitOfMeasure,
		&i.ReorderLevel,
		&i.ImageURL,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listProducts = `-- name: ListProducts :many
SELECT p.id, p.name, p.barcode, p.sku, p.category_id, c.name as category_name,
       p.unit_price, p.cost_price, p.tax_rate, p.unit_of_measure,
       p.reorder_level, p.image_url, p.is_active, p.created_at, p.updated_at
FROM products p
JOIN categories c ON p.category_id = c.id
WHERE ($3::text IS NULL OR p.name ILIKE '%' || $3::text || '%' OR p.barcode ILIKE '%' || $3::text || '%')
  AND ($4::uuid IS NULL OR p.category_id = $4::uuid)
  AND ($5::boolean IS NULL OR p.is_active = $5::boolean)
ORDER BY p.name ASC
LIMIT $1 OFFSET $2
`

type ListProductsParams struct {
	Limit      int32          `json:"limit"`
	Offset     int32          `json:"offset"`
	Search     sql.NullString `json:"search"`
	CategoryID *uuid.UUID     `json:"category_id"`
	IsActive   sql.NullBool   `json:"is_active"`
}

type ListProductsRow struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	CategoryName  string         `json:"category_name"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
	IsActive      bool           `json:"is_active"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

func (q *Queries) ListProducts(ctx context.Context, arg ListProductsParams) ([]ListProductsRow, error) {
	rows, err := q.db.QueryContext(ctx, listProducts,
		arg.Limit,
		arg.Offset,
		arg.Search,
		arg.CategoryID,
		arg.IsActive,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListProductsRow
	for rows.Next() {
		var i ListProductsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Barcode,
			&i.Sku,
			&i.CategoryID,
			&i.CategoryName,
			&i.UnitPrice,
			&i.CostPrice,
			&i.TaxRate,
			&i.UnitOfMeasure,
			&i.ReorderLevel,
			&i.ImageURL,
			&i.IsActive,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const countProducts = `-- name: CountProducts :one
SELECT COUNT(*)
FROM products p
WHERE ($1::text IS NULL OR p.name ILIKE '%' || $1::text || '%' OR p.barcode ILIKE '%' || $1::text || '%')
  AND ($2::uuid IS NULL OR p.category_id = $2::uuid)
  AND ($3::boolean IS NULL OR p.is_active = $3::boolean)
`

type CountProductsParams struct {
	Search     sql.NullString `json:"search"`
	CategoryID *uuid.UUID     `json:"category_id"`
	IsActive   sql.NullBool   `json:"is_active"`
}

func (q *Queries) CountProducts(ctx context.Context, arg CountProductsParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, countProducts, arg.Search, arg.CategoryID, arg.IsActive)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createProduct = `-- name: CreateProduct :one
INSERT INTO products (
    name, barcode, sku, category_id, unit_price, cost_price,
    tax_rate, unit_of_measure, reorder_level, image_url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, name, barcode, sku, category_id, unit_price, cost_price,
          tax_rate, unit_of_measure, reorder_level, image_url, is_active,
          created_at, updated_at
`

type CreateProductParams struct {
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
}

func (q *Queries) CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, createProduct,
		arg.Name,
		arg.Barcode,
		arg.Sku,
		arg.CategoryID,
		arg.UnitPrice,
		arg.CostPrice,
		arg.TaxRate,
		arg.UnitOfMeasure,
		arg.ReorderLevel,
		arg.ImageURL,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Barcode,
		&i.Sku,
		&i.CategoryID,
		&i.UnitPrice,
		&i.CostPrice,
		&i.TaxRate,
		&i.UnitOfMeasure,
		&i.ReorderLevel,
		&i.ImageURL,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateProduct = `-- name: UpdateProduct :one
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
          created_at, updated_at
`

type UpdateProductParams struct {
	ID            uuid.UUID      `json:"id"`
	Name          string         `json:"name"`
	Barcode       string         `json:"barcode"`
	Sku           sql.NullString `json:"sku"`
	CategoryID    uuid.UUID      `json:"category_id"`
	UnitPrice     int32          `json:"unit_price"`
	CostPrice     int32          `json:"cost_price"`
	TaxRate       string         `json:"tax_rate"`
	UnitOfMeasure string         `json:"unit_of_measure"`
	ReorderLevel  int32          `json:"reorder_level"`
	ImageURL      sql.NullString `json:"image_url"`
	IsActive      bool           `json:"is_active"`
}

func (q *Queries) UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProduct,
		arg.ID,
		arg.Name,
		arg.Barcode,
		arg.Sku,
		arg.CategoryID,
		arg.UnitPrice,
		arg.CostPrice,
		arg.TaxRate,
		arg.UnitOfMeasure,
		arg.ReorderLevel,
		arg.ImageURL,
		arg.IsActive,
	)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Barcode,
		&i.Sku,
		&i.CategoryID,
		&i.UnitPrice,
		&i.CostPrice,
		&i.TaxRate,
		&i.UnitOfMeasure,
		&i.ReorderLevel,
		&i.ImageURL,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deactivateProduct = `-- name: DeactivateProduct :one
UPDATE products
SET is_active = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, barcode, sku, category_id, unit_price, cost_price,
          tax_rate, unit_of_measure, reorder_level, image_url, is_active,
          created_at, updated_at
`

type DeactivateProductParams struct {
	ID       uuid.UUID `json:"id"`
	IsActive bool      `json:"is_active"`
}

func (q *Queries) DeactivateProduct(ctx context.Context, arg DeactivateProductParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, deactivateProduct, arg.ID, arg.IsActive)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Barcode,
		&i.Sku,
		&i.CategoryID,
		&i.UnitPrice,
		&i.CostPrice,
		&i.TaxRate,
		&i.UnitOfMeasure,
		&i.ReorderLevel,
		&i.ImageURL,
		&i.IsActive,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteProduct = `-- name: DeleteProduct :one
DELETE FROM products
WHERE id = $1
RETURNING id
`

func (q *Queries) DeleteProduct(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var deletedID uuid.UUID
	err := q.db.QueryRowContext(ctx, deleteProduct, id).Scan(&deletedID)
	return deletedID, err
}
