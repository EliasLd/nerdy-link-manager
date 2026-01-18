package router

import (
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/internal/handlers"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheck)

	return mux
}
