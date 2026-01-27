package handlers

import (
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/EliasLd/nerdy-link-manager/internal/deploy"
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

	var deployMutex sync.Mutex

	go func() {
		// Lock the deployment script to avoid simultaneous deployments
		deployMutex.Lock()
		defer deployMutex.Unlock()

		if err := deploy.RunDeployment(); err != nil {
			log.Printf("%v\n", err)
		} else {
			log.Println("[DEPLOY] Deployment completed successfully!")
		}
	}()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Nerdy link manager deployment triggered!"))
}
