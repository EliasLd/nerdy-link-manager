package router

import (
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/internal/handlers"
	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)

	return middlware.CORS(mux)
}
