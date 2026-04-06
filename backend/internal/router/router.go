package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/EliasLd/nerdy-link-manager/internal/handlers"
	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
)

func New(
	authHandler *handlers.AuthHandler,
	LinkHandler *handlers.LinkHandler,
	authMiddleware func(http.Handler) http.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("POST /api/login", authHandler.Login)

	// Protected auth routes
	mux.Handle("POST /api/register", authMiddleware(http.HandlerFunc(authHandler.Register)))

	// link routes (protected)
	mux.Handle("/api/links", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			LinkHandler.GetAllLinks(w, r)
		case http.MethodPost:
			LinkHandler.CreateLink(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// link detail routes: /links/{id}
	mux.Handle("/api/links/", authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// /links/{id}/click
		if strings.HasSuffix(r.URL.Path, "/click") {
			if r.Method == http.MethodPost {
				LinkHandler.TrackClick(w, r)
			} else {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// /links/{id}
		switch r.Method {
		case http.MethodGet:
			LinkHandler.GetLink(w, r)
		case http.MethodPut:
			LinkHandler.UpdateLink(w, r)
		case http.MethodDelete:
			LinkHandler.DeleteLink(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	if os.Getenv("ENABLE_DEPLOYMENT_ENDPOINT") == "true" {
		mux.HandleFunc("/deploy", handlers.DeployHandler)
	}

	// Apply global CORS middleware
	return middleware.CORS(mux)
}
