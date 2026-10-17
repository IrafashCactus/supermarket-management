package service

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"

	"github.com/google/uuid"
	"github.com/supermarket/backend/internal/database/sqlc"
	"github.com/supermarket/backend/internal/dto"
)

type CategoryService struct {
	queries sqlc.Querier
}

func NewCategoryService(queries sqlc.Querier) *CategoryService {
	return &CategoryService{queries: queries}
}

func (s *CategoryService) List(ctx context.Context, page, perPage int) (dto.ListResponse[dto.CategoryResponse], error) {
	page, perPage = normalizePagination(page, perPage)
	categories, err := s.queries.ListCategories(ctx, sqlc.ListCategoriesParams{
		Limit:  int32(perPage),
		Offset: int32((page - 1) * perPage),
	})
	if err != nil {
		return dto.ListResponse[dto.CategoryResponse]{}, err
	}
	total, err := s.queries.CountCategories(ctx)
	if err != nil {
		return dto.ListResponse[dto.CategoryResponse]{}, err
	}

	items := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		productCount, err := s.queries.CountProductsByCategory(ctx, category.ID)
		if err != nil {
			return dto.ListResponse[dto.CategoryResponse]{}, err
		}
		items = append(items, mapCategory(category, productCount))
	}
	return dto.ListResponse[dto.CategoryResponse]{
		Data: items,
		Meta: dto.PaginationMeta{
			Page: page, PerPage: perPage, Total: total,
			TotalPages: int(math.Ceil(float64(total) / float64(perPage))),
		},
	}, nil
}

func (s *CategoryService) Create(ctx context.Context, req dto.CreateCategoryRequest) (dto.CategoryResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return dto.CategoryResponse{}, ErrInvalidCatalogRequest
	}
	category, err := s.queries.CreateCategory(ctx, sqlc.CreateCategoryParams{
		Name: name, Description: nullableString(req.Description),
	})
	if err != nil {
		return dto.CategoryResponse{}, catalogDBError(err)
	}
	return mapCategory(category, 0), nil
}

func (s *CategoryService) Get(ctx context.Context, id uuid.UUID) (dto.CategoryResponse, error) {
	category, err := s.queries.GetCategoryByID(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, catalogDBError(err)
	}
	count, err := s.queries.CountProductsByCategory(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	return mapCategory(category, count), nil
}

func (s *CategoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateCategoryRequest) (dto.CategoryResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return dto.CategoryResponse{}, ErrInvalidCatalogRequest
	}
	category, err := s.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID: id, Name: name, Description: nullableString(req.Description),
	})
	if err != nil {
		return dto.CategoryResponse{}, catalogDBError(err)
	}
	count, err := s.queries.CountProductsByCategory(ctx, id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}
	return mapCategory(category, count), nil
}

func (s *CategoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.queries.GetCategoryByID(ctx, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	count, err := s.queries.CountProductsByCategory(ctx, id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryHasProducts
	}
	return catalogDBError(s.queries.DeleteCategory(ctx, id))
}

func mapCategory(category sqlc.Category, productCount int64) dto.CategoryResponse {
	description := ""
	if category.Description.Valid {
		description = category.Description.String
	}
	return dto.CategoryResponse{
		ID: category.ID, Name: category.Name, Description: description,
		ProductCount: productCount, CreatedAt: category.CreatedAt, UpdatedAt: category.UpdatedAt,
	}
}

func nullableString(value string) sql.NullString {
	value = strings.TrimSpace(value)
	return sql.NullString{String: value, Valid: value != ""}
}

func normalizePagination(page, perPage int) (int, int) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	return page, perPage
}

var (
	ErrInvalidCatalogRequest = errors.New("invalid catalog request")
	ErrCategoryHasProducts   = errors.New("category has products")
	ErrNotFound              = errors.New("not found")
	ErrConflict              = errors.New("conflict")
)
