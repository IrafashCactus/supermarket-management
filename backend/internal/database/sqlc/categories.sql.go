package sqlc

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getCategoryByID = `-- name: GetCategoryByID :one
SELECT id, name, description, created_at, updated_at
FROM categories
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetCategoryByID(ctx context.Context, id uuid.UUID) (Category, error) {
	row := q.db.QueryRowContext(ctx, getCategoryByID, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT id, name, description, created_at, updated_at
FROM categories
ORDER BY name ASC
LIMIT $1 OFFSET $2
`

type ListCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, listCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
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

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (
    name, description
) VALUES (
    $1, $2
)
RETURNING id, name, description, created_at, updated_at
`

type CreateCategoryParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, createCategory, arg.Name, arg.Description)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories
SET name = $2,
    description = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING id, name, description, created_at, updated_at
`

type UpdateCategoryParams struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.db.QueryRowContext(ctx, updateCategory, arg.ID, arg.Name, arg.Description)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteCategory, id)
	return err
}

const countCategories = `-- name: CountCategories :one
SELECT COUNT(*) FROM categories
`

func (q *Queries) CountCategories(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countCategories)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countProductsByCategory = `-- name: CountProductsByCategory :one
SELECT COUNT(*) FROM products WHERE category_id = $1
`

func (q *Queries) CountProductsByCategory(ctx context.Context, categoryID uuid.UUID) (int64, error) {
	row := q.db.QueryRowContext(ctx, countProductsByCategory, categoryID)
	var count int64
	err := row.Scan(&count)
	return count, err
}
