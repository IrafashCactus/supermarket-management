package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supermarket/backend/internal/auth"
	"github.com/supermarket/backend/internal/database"
	"github.com/supermarket/backend/internal/database/sqlc"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := database.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
	}{
		{"Admin User", "admin@supermarket.rw", "admin123", "admin"},
		{"Manager User", "manager@supermarket.rw", "manager123", "manager"},
		{"Accountant User", "accountant@supermarket.rw", "accountant123", "accountant"},
		{"Storekeeper User", "storekeeper@supermarket.rw", "storekeeper123", "storekeeper"},
		{"Cashier User", "cashier@supermarket.rw", "cashier123", "cashier"},
	}

	for _, u := range users {
		hash, err := auth.HashPassword(u.Password)
		if err != nil {
			log.Fatalf("Failed to hash password for %s: %v", u.Email, err)
		}

		_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
			Name:         u.Name,
			Email:        u.Email,
			PasswordHash: hash,
			Role:         u.Role,
		})
		
		if err != nil {
			log.Printf("Could not create user %s (might already exist): %v", u.Email, err)
		} else {
			log.Printf("Successfully created user: %s (%s)", u.Name, u.Role)
		}
	}

	log.Println("Seeding complete.")
}
