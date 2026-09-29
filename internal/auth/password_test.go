package auth

import (
	"testing"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "securepassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}

	if hash == password {
		t.Errorf("expected hash to differ from raw password")
	}

	// Verify valid password
	if err := CheckPassword(password, hash); err != nil {
		t.Errorf("expected password to match hash, got error: %v", err)
	}

	// Verify invalid password
	if err := CheckPassword("wrongpassword", hash); err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}
}
