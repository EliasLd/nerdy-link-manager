package deploy

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunDeployment() error {
	scriptPath := os.Getenv("DEPLOY_SCRIPT_PATH")
	if scriptPath == "" {
		return fmt.Errorf("[DEPLOY][ERROR] DEPLOY_SCRIPT_PATH is not set")
	}

	workDir := os.Getenv("DEPLOY_WORKDIR")
	if workDir == "" {
		fmt.Println("[DEPLOY][WARN] DEPLOY_WORKDIR is not set, pursuing with '/'")
		workDir = "/"
	}

	cmd := exec.Command(scriptPath)
	cmd.Dir = workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"[DEPLOY][ERROR] Deployment script failed: %w\nstdout:\n%s\nstderr:\n%s",
			err, stdout.String(), stderr.String(),
		)
	}

	fmt.Printf("[DEPLOY][INFO] Script output:\n%s\n", stdout.String())
	return nil
}
