package main

import (
	"log"
	"net/http"
	"time"

	"github.com/EliasLd/nerdy-link-manager/config"
	"github.com/EliasLd/nerdy-link-manager/internal/auth"
	"github.com/EliasLd/nerdy-link-manager/internal/db"
	"github.com/EliasLd/nerdy-link-manager/internal/handlers"
	"github.com/EliasLd/nerdy-link-manager/internal/middleware"
	"github.com/EliasLd/nerdy-link-manager/internal/repositories"
	"github.com/EliasLd/nerdy-link-manager/internal/router"
	"github.com/EliasLd/nerdy-link-manager/internal/services"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("[ERROR] Database connection failed: %v", err)
	}
	defer database.DB.Close()

	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepo)

	linkRepo := repositories.NewLinkRepository(database.DB)
	linkService := services.NewLinkService(linkRepo)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, time.Hour*3600)

	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	linkHandler := handlers.NewLinkHandler(linkService)

	authMiddleware := middleware.AuthMiddleware(*jwtManager)

	r := router.New(authHandler, linkHandler, authMiddleware)

	log.Println("[INFO] Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
