package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/supermarket/backend/internal/database/sqlc"
	"github.com/supermarket/backend/internal/dto"
)

type ProductService struct {
	queries sqlc.Querier
}

func NewProductService(queries sqlc.Querier) *ProductService {
	return &ProductService{queries: queries}
}

func (s *ProductService) List(ctx context.Context, page, perPage int, search, categoryID, active string) (dto.ListResponse[dto.ProductResponse], error) {
	page, perPage = normalizePagination(page, perPage)
	params := sqlc.ListProductsParams{
		Limit:  int32(perPage),
		Offset: int32((page - 1) * perPage),
		Search: nullableString(search),
	}
	if categoryID != "" {
		id, err := uuid.Parse(categoryID)
		if err != nil {
			return dto.ListResponse[dto.ProductResponse]{}, ErrInvalidCatalogRequest
		}
		params.CategoryID = &id
	}
	if active != "" {
		value, err := strconv.ParseBool(active)
		if err != nil {
			return dto.ListResponse[dto.ProductResponse]{}, ErrInvalidCatalogRequest
		}
		params.IsActive = sql.NullBool{Bool: value, Valid: true}
	}
	products, err := s.queries.ListProducts(ctx, params)
	if err != nil {
		return dto.ListResponse[dto.ProductResponse]{}, err
	}
	total, err := s.queries.CountProducts(ctx, sqlc.CountProductsParams{
		Search: params.Search, CategoryID: params.CategoryID, IsActive: params.IsActive,
	})
	if err != nil {
		return dto.ListResponse[dto.ProductResponse]{}, err
	}
	items := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		item, err := mapProduct(product.ID, product.Name, product.Barcode, product.Sku, product.CategoryID, product.CategoryName, product.UnitPrice, product.CostPrice, product.TaxRate, product.UnitOfMeasure, product.ReorderLevel, product.ImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt)
		if err != nil {
			return dto.ListResponse[dto.ProductResponse]{}, err
		}
		items = append(items, item)
	}
	return dto.ListResponse[dto.ProductResponse]{
		Data: items,
		Meta: dto.PaginationMeta{
			Page: page, PerPage: perPage, Total: total,
			TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
		},
	}, nil
}

func (s *ProductService) Get(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	product, err := s.queries.GetProductByID(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, catalogDBError(err)
	}
	return mapProduct(product.ID, product.Name, product.Barcode, product.Sku, product.CategoryID, product.CategoryName, product.UnitPrice, product.CostPrice, product.TaxRate, product.UnitOfMeasure, product.ReorderLevel, product.ImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt)
}

func (s *ProductService) GetByBarcode(ctx context.Context, barcode string) (dto.ProductResponse, error) {
	product, err := s.queries.GetProductByBarcode(ctx, barcode)
	if err != nil {
		return dto.ProductResponse{}, catalogDBError(err)
	}
	return mapProduct(product.ID, product.Name, product.Barcode, product.Sku, product.CategoryID, product.CategoryName, product.UnitPrice, product.CostPrice, product.TaxRate, product.UnitOfMeasure, product.ReorderLevel, product.ImageURL, product.IsActive, product.CreatedAt, product.UpdatedAt)
}

func (s *ProductService) Create(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	values, err := validateProduct(req)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	product, err := s.queries.CreateProduct(ctx, values)
	if err != nil {
		return dto.ProductResponse{}, catalogDBError(err)
	}
	return s.Get(ctx, product.ID)
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	values, err := validateProduct(req.CreateProductRequest)
	if err != nil {
		return dto.ProductResponse{}, err
	}
	product, err := s.queries.UpdateProduct(ctx, sqlc.UpdateProductParams{
		ID: id, Name: values.Name, Barcode: values.Barcode, Sku: values.Sku,
		CategoryID: values.CategoryID, UnitPrice: values.UnitPrice, CostPrice: values.CostPrice,
		TaxRate: values.TaxRate, UnitOfMeasure: values.UnitOfMeasure,
		ReorderLevel: values.ReorderLevel, ImageURL: values.ImageURL, IsActive: req.IsActive,
	})
	if err != nil {
		return dto.ProductResponse{}, catalogDBError(err)
	}
	return s.Get(ctx, product.ID)
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.queries.DeleteProduct(ctx, id)
	if err != nil {
		return catalogDBError(err)
	}
	return nil
}

func validateProduct(req dto.CreateProductRequest) (sqlc.CreateProductParams, error) {
	categoryID, err := uuid.Parse(req.CategoryID)
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Barcode) == "" ||
		strings.TrimSpace(req.UnitOfMeasure) == "" || err != nil ||
		req.UnitPrice < 0 || req.CostPrice < 0 || req.ReorderLevel < 0 ||
		math.IsNaN(req.TaxRate) || math.IsInf(req.TaxRate, 0) || req.TaxRate < 0 || req.TaxRate > 100 {
		return sqlc.CreateProductParams{}, ErrInvalidCatalogRequest
	}
	return sqlc.CreateProductParams{
		Name: strings.TrimSpace(req.Name), Barcode: strings.TrimSpace(req.Barcode),
		Sku: nullableString(req.SKU), CategoryID: categoryID,
		UnitPrice: req.UnitPrice, CostPrice: req.CostPrice,
		TaxRate:       strconv.FormatFloat(req.TaxRate, 'f', 2, 64),
		UnitOfMeasure: strings.TrimSpace(req.UnitOfMeasure),
		ReorderLevel:  req.ReorderLevel, ImageURL: nullableString(req.ImageURL),
	}, nil
}

func mapProduct(id uuid.UUID, name, barcode string, sku sql.NullString, categoryID uuid.UUID, categoryName string, unitPrice, costPrice int32, taxRate, unitOfMeasure string, reorderLevel int32, imageURL sql.NullString, active bool, createdAt, updatedAt time.Time) (dto.ProductResponse, error) {
	tax, err := strconv.ParseFloat(taxRate, 64)
	if err != nil {
		return dto.ProductResponse{}, fmt.Errorf("parse product tax rate: %w", err)
	}
	response := dto.ProductResponse{
		ID: id, Name: name, Barcode: barcode, CategoryID: categoryID, CategoryName: categoryName,
		UnitPrice: unitPrice, CostPrice: costPrice, TaxRate: tax, UnitOfMeasure: unitOfMeasure,
		ReorderLevel: reorderLevel, IsActive: active, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
	if sku.Valid {
		response.SKU = sku.String
	}
	if imageURL.Valid {
		response.ImageURL = imageURL.String
	}
	return response, nil
}

func catalogDBError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505", "23503":
			return ErrConflict
		}
	}
	return err
}
