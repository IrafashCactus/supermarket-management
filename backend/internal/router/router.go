package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/supermarket/backend/internal/auth"
	"github.com/supermarket/backend/internal/handler"
	"github.com/supermarket/backend/internal/middleware"
)

func SetupRouter(
	frontendURL string,
	jwtManager *auth.JWTManager,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	categoryHandler *handler.CategoryHandler,
	productHandler *handler.ProductHandler,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logging)
	r.Use(chiMiddleware.Recoverer)
	r.Use(middleware.CORS(frontendURL))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", authHandler.Login)
			r.Post("/refresh", authHandler.Refresh)

			r.Group(func(r chi.Router) {
				r.Use(middleware.Authenticator(jwtManager))
				r.Post("/logout", authHandler.Logout)
				r.Get("/me", authHandler.Me)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Authenticator(jwtManager))
			r.Use(middleware.RequireRole("admin"))

			r.Get("/", userHandler.List)
			r.Post("/", userHandler.Create)
			r.Get("/{id}", userHandler.Get)
			r.Put("/{id}", userHandler.Update)
			r.Put("/{id}/password", userHandler.ResetPassword)
			r.Post("/{id}/reset-password", userHandler.ResetPassword)
			r.Put("/{id}/reset-password", userHandler.ResetPassword)
		})

		r.Route("/categories", func(r chi.Router) {
			r.Use(middleware.Authenticator(jwtManager))
			r.Use(middleware.RequireRole("admin", "manager"))

			r.Get("/", categoryHandler.List)
			r.Post("/", categoryHandler.Create)
			r.Get("/{id}", categoryHandler.Get)
			r.Put("/{id}", categoryHandler.Update)
			r.Delete("/{id}", categoryHandler.Delete)
		})

		r.Route("/products", func(r chi.Router) {
			r.Use(middleware.Authenticator(jwtManager))
			r.Use(middleware.RequireRole("admin", "manager", "storekeeper"))

			r.Get("/", productHandler.List)
			r.Post("/", productHandler.Create)
			r.Get("/barcode/{barcode}", productHandler.GetByBarcode)
			r.Get("/{id}", productHandler.Get)
			r.Put("/{id}", productHandler.Update)
			r.Delete("/{id}", productHandler.Delete)
		})
	})

	return r
}
