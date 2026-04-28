package main

import (
	"context"
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

	folderRepo := repositories.NewFolderRepository(database.DB)
	folderService := services.NewFolderService(folderRepo)

	linkRepo := repositories.NewLinkRepository(database.DB)
	linkService := services.NewLinkService(linkRepo, folderRepo)

	jwtManager := auth.NewJWTManager(cfg.JWTSecret, time.Hour*3600)

	authHandler := handlers.NewAuthHandler(userService, jwtManager)
	folderHandler := handlers.NewFolderHandler(folderService)
	linkHandler := handlers.NewLinkHandler(linkService)

	created, err := userService.BootstrapInitialUser(
		context.Background(),
		cfg.InitialAdminEmail,
		cfg.InitialAdminPassword,
	)
	if err != nil {
		log.Fatalf("[ERROR] Failed to bootstrap intial user: %v", err)
	}
	if created {
		log.Println("[INFO] Initial admin user created")
	} else {
		log.Println("[INFO] Initial admin bootstrap skipped")
	}

	authMiddleware := middleware.AuthMiddleware(*jwtManager)

	r := router.New(authHandler, linkHandler, folderHandler, authMiddleware)

	log.Println("[INFO] Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
