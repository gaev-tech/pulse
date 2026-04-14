package v1

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authmiddleware "github.com/gaevivan/pulse/internal/handler/middleware"
)

type Deps struct {
	Auth    *AuthHandler
	AuthMW  *authmiddleware.Auth
}

func NewRouter(deps Deps) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Get("/health", healthHandler)

	router.Route("/v1", func(router chi.Router) {
		router.Route("/auth", func(router chi.Router) {
			router.Post("/magic-link", deps.Auth.SendMagicLink)
			router.Post("/magic-link/verify", deps.Auth.VerifyMagicLink)
			router.Post("/refresh", deps.Auth.Refresh)
			router.Post("/logout", deps.Auth.Logout)

			router.With(deps.AuthMW.Required).Get("/me", deps.Auth.Me)
		})
	})

	return router
}

func healthHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(map[string]string{"status": "ok"})
}
