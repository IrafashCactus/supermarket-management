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

type ProductHandler struct {
	productService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))
	result, err := h.productService.List(r.Context(), page, perPage, query.Get("search"), query.Get("category_id"), query.Get("is_active"))
	if err != nil {
		writeCatalogError(w, err, "LIST_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, result)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}
	product, err := h.productService.Get(r.Context(), id)
	if err != nil {
		writeProductReadError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.ProductResponse]{Data: product})
}

func (h *ProductHandler) GetByBarcode(w http.ResponseWriter, r *http.Request) {
	product, err := h.productService.GetByBarcode(r.Context(), chi.URLParam(r, "barcode"))
	if err != nil {
		writeProductReadError(w, err)
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.ProductResponse]{Data: product})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	product, err := h.productService.Create(r.Context(), req)
	if err != nil {
		writeCatalogError(w, err, "CREATE_FAILED")
		return
	}
	respondJSON(w, http.StatusCreated, dto.SuccessResponse[dto.ProductResponse]{Data: product})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}
	var req dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}
	product, err := h.productService.Update(r.Context(), id, req)
	if err != nil {
		writeCatalogError(w, err, "UPDATE_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.ProductResponse]{Data: product})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID")
		return
	}
	if err := h.productService.Delete(r.Context(), id); err != nil {
		writeCatalogError(w, err, "DELETE_FAILED")
		return
	}
	respondJSON(w, http.StatusOK, dto.SuccessResponse[any]{Message: "Product deleted successfully"})
}

func writeProductReadError(w http.ResponseWriter, err error) {
	if errors.Is(err, service.ErrNotFound) {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "Product not found")
		return
	}
	respondError(w, http.StatusInternalServerError, "LOOKUP_FAILED", "Failed to look up product")
}
