package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/supermarket/backend/internal/dto"
	"github.com/supermarket/backend/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, dto.SuccessResponse[dto.UserResponse]{
		Data:    user,
		Message: "User created successfully",
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))

	listResp, err := h.userService.ListUsers(r.Context(), page, perPage)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, listResp)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	user, err := h.userService.GetUser(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "NOT_FOUND", "User not found")
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.UserResponse]{
		Data: user,
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	user, err := h.userService.UpdateUser(r.Context(), id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.UserResponse]{
		Data:    user,
		Message: "User updated successfully",
	})
}

func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_ID", "Invalid user ID")
		return
	}

	var req dto.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	if err := h.userService.ResetPassword(r.Context(), id, req); err != nil {
		respondError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse[any]{
		Message: "Password reset successfully",
	})
}
