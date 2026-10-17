package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/supermarket/backend/internal/dto"
	"github.com/supermarket/backend/internal/service"
)

type CategoryHandler struct {
	categoryService *service.CategoryService
}

func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	result, err := h.categoryService.List(r.Context(), page, perPage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "LIST_FAILED", "Failed to list categories")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	category, err := h.categoryService.Create(r.Context(), req)
	if err != nil {
		writeCatalogError(w, err, "CREATE_FAILED")
		return
	}
	respondJSON(w, http.StatusCreated, dto.SuccessResponse[dto.CategoryResponse]{Data: category})
}

func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseCatalogID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid category ID")
		return
	}
	category, err := h.categoryService.Get(r.Context(), id)
	if err != nil {
		writeCatalogError(w, err, "LOOKUP_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.CategoryResponse]{Data: category})
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseCatalogID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid category ID")
		return
	}
	var req dto.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	category, err := h.categoryService.Update(r.Context(), id, req)
	if err != nil {
		writeCatalogError(w, err, "UPDATE_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.CategoryResponse]{Data: category})
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseCatalogID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid category ID")
		return
	}
	if err := h.categoryService.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "DELETE_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[any]{Message: "Category deleted successfully"})
}

func parseCatalogID(r *http.Request) (uuid.UUID, error) {
	return uuid.Parse(chi.URLParam(r, "id"))
}

func writeCatalogError(w http.ResponseWriter, err error, code string) {
	switch {
	case errors.Is(err, service.ErrInvalidCatalogRequest):
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid category or product fields")
	case errors.Is(err, service.ErrCategoryHasProducts):
		respondError(w, http.StatusConflict, "CATEGORY_HAS_PRODUCTS", err.Error())
	case errors.Is(err, service.ErrNotFound):
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Category or product not found")
	case errors.Is(err, service.ErrConflict):
		respondError(w, http.StatusConflict, "CONFLICT", "A category or product with these values already exists")
	default:
		respondError(w, http.StatusInternalServerError, code, "Catalog operation failed")
	}
}
