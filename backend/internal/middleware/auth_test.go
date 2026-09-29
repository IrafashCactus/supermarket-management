package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/supermarket/backend/internal/auth"
)

func TestRequireRole_Allowed(t *testing.T) {
	jwtManager := auth.NewJWTManager("testsecret12345678901234567890123456", time.Hour, time.Hour)
	handler := RequireRole("admin", "manager")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Admin role
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, &auth.Claims{
		UserID: uuid.New(),
		Email:  "admin@test.rw",
		Role:   "admin",
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", rr.Code)
	}
	_ = jwtManager
}

func TestRequireRole_Forbidden(t *testing.T) {
	handler := RequireRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Cashier role trying to access admin endpoint
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), UserContextKey, &auth.Claims{
		UserID: uuid.New(),
		Email:  "cashier@test.rw",
		Role:   "cashier",
	})

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req.WithContext(ctx))

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden, got %d", rr.Code)
	}
}
