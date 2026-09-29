package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/supermarket/backend/internal/auth"
	"github.com/supermarket/backend/internal/config"
	"github.com/supermarket/backend/internal/database"
	"github.com/supermarket/backend/internal/database/sqlc"
	"github.com/supermarket/backend/internal/handler"
	"github.com/supermarket/backend/internal/router"
	"github.com/supermarket/backend/internal/service"
	
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	_ = godotenv.Load()

	cfg := config.LoadConfig()

	// Run migrations
	runMigrations(cfg.DatabaseURL)

	db, err := database.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTAccessExpiry, cfg.JWTRefreshExpiry)

	authService := service.NewAuthService(queries, jwtManager)
	userService := service.NewUserService(queries)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	r := router.SetupRouter(cfg.FrontendURL, jwtManager, authHandler, userHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}

func runMigrations(dbURL string) {
	log.Println("Running database migrations...")
	m, err := migrate.New("file://internal/database/migrations", dbURL)
	if err != nil {
		log.Printf("Failed to init migrate: %v", err)
		return
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations applied successfully")
}
