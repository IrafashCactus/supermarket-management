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
		})
	})

	return r
}
