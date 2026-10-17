package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCategoryRequest CreateCategoryRequest

type CategoryResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	ProductCount int64     `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name          string  `json:"name"`
	Barcode       string  `json:"barcode"`
	SKU           string  `json:"sku"`
	CategoryID    string  `json:"category_id"`
	UnitPrice     int32   `json:"unit_price"`
	CostPrice     int32   `json:"cost_price"`
	TaxRate       float64 `json:"tax_rate"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	ReorderLevel  int32   `json:"reorder_level"`
	ImageURL      string  `json:"image_url"`
}

type UpdateProductRequest struct {
	CreateProductRequest
	IsActive bool `json:"is_active"`
}

type ProductResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Barcode       string    `json:"barcode"`
	SKU           string    `json:"sku"`
	CategoryID    uuid.UUID `json:"category_id"`
	CategoryName  string    `json:"category_name"`
	UnitPrice     int32     `json:"unit_price"`
	CostPrice     int32     `json:"cost_price"`
	TaxRate       float64   `json:"tax_rate"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	ReorderLevel  int32     `json:"reorder_level"`
	ImageURL      string    `json:"image_url"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
