package router

import (
	"net/http"
	"os"

	"github.com/EliasLd/nerdy-link-manager/internal/handlers"
	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
)

func New(
	authHandler *handlers.AuthHandler,
	authMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	if os.Getenv("ENABLE_DEPLOYMENT_ENDPOINT") == "true" {
		mux.HandleFunc("/deploy", handlers.DeployHandler)
	}

	return middleware.CORS(mux)
}
