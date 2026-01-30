package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DeployHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		log.Println("[INFO] Illegal request received (not POST)")
		return
	}

	enabled := os.Getenv("ENABLE_DEPLOYMENT_ENDPOINT")
	if enabled != "true" {
		http.Error(w, "deployment endpoint disabled", http.StatusForbidden)
		return
	}

	expectedToken := os.Getenv("DEPLOY_TOKEN")
	if r.Header.Get("X-Deploy-Token") != expectedToken {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		log.Println("[WARN] Deployment request with invalid token received")
		return
	}

	log.Println("[DEPLOY] Deployment request received, triggering deployment script...")

	triggerDir := os.Getenv("DEPLOY_TRIGGER_DIR")
	if triggerDir == "" {
		triggerDir = "/deploy-trigger"
	}

	if err := os.MkdirAll(triggerDir, 0755); err != nil {
		log.Printf("[DEPLOY][ERROR] Could not create trigger dir: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	triggerFile := filepath.Join(triggerDir, fmt.Sprintf("trigger-%d", time.Now().UnixNano()))
	if err := os.WriteFile(triggerFile, []byte{}, 0644); err != nil {
		log.Printf("[DEPLOY][ERROR] Could not create trigger file: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("[DEPLOY] Trigger file created: %s", triggerFile)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deployment trigger created"))
}
