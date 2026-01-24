package handlers

import (
	"log"
	"net/http"
	"os"
)

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	enabled := os.Getenv("ENABLE_DEPLOYMENT_ENDPOINT")
	if enabled != "true" {
		http.Error(w, "deployment endpoint disabled", http.StatusForbidden)
		return
	}

	log.Println("[DEPLOY] Deployment request received")

	// TODO: Call deployment script
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Nerdy link manager deployment triggered!"))
}
