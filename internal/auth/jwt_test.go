package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTManager_GenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key-12345678901234567890"
	manager := NewJWTManager(secret, 15*time.Minute, 7*24*time.Hour)

	userID := uuid.New()
	email := "cashier@supermarket.rw"
	role := "cashier"

	tokenStr, err := manager.GenerateAccessToken(userID, email, role)
	if err != nil {
		t.Fatalf("unexpected error generating access token: %v", err)
	}

	claims, err := manager.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("unexpected error validating token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %v, got %v", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestJWTManager_ExpiredToken(t *testing.T) {
	secret := "test-secret-key-12345678901234567890"
	// Expired immediately (-1 second)
	manager := NewJWTManager(secret, -1*time.Second, 7*24*time.Hour)

	userID := uuid.New()
	tokenStr, err := manager.GenerateAccessToken(userID, "expired@test.rw", "admin")
	if err != nil {
		t.Fatalf("unexpected error generating token: %v", err)
	}

	_, err = manager.ValidateToken(tokenStr)
	if err == nil {
		t.Errorf("expected validation error for expired token, got nil")
	}
}
