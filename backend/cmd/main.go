package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/EliasLd/nerdy-link-manager/config"
)

func main() {
	cfg := config.LoadConfig()

	fmt.Println("Starting server on port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
