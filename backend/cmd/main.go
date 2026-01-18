package main

import (
	"log"
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/config"
	"github.com/EliasLd/nerdy-link-manager/internal/db"
	"github.com/EliasLd/nerdy-link-manager/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("[ERROR] Database connection failed: %v", err)
	}
	defer database.DB.Close()

	r := router.New()

	log.Println("[INFO] Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
