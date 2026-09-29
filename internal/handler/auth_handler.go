package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/supermarket/backend/internal/dto"
	"github.com/supermarket/backend/internal/middleware"
	"github.com/supermarket/backend/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	resp, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "AUTH_FAILED", err.Error())
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false, // Set to true in production
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/auth",
		MaxAge:   7 * 24 * 3600, // 7 days
	})

	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.LoginResponse]{
		Data:    resp,
		Message: "Login successful",
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		respondError(w, http.StatusUnauthorized, "MISSING_REFRESH_TOKEN", "Refresh token is missing")
		return
	}

	accessToken, err := h.authService.RefreshToken(r.Context(), cookie.Value)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", err.Error())
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse[dto.RefreshResponse]{
		Data: dto.RefreshResponse{
			AccessToken: accessToken,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // Set true in production
		SameSite: http.SameSiteStrictMode,
		Path:     "/api/v1/auth",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	respondJSON(w, http.StatusOK, dto.SuccessResponse[any]{
		Message: "Logged out successfully",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetUserFromContext(r.Context())
	if claims == nil {
		respondError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Not authenticated")
		return
	}

	respondJSON(w, http.StatusOK, dto.SuccessResponse[any]{
		Data: claims,
	})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dto.ErrorResponse{
		Error: dto.ErrorDetails{
			Code:    code,
			Message: message,
		},
	})
}
