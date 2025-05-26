package git
import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"github.com/google/go-github/v68/github"
	log "github.com/sirupsen/logrus"
)


// Clone a GitHub repository
func CloneRepository() {
	fullRepoName := os.Getenv("GITHUB_REPOSITORY") // Expected format: "owner/repo"
	repoURL := fmt.Sprintf("https://github.com/%s.git", fullRepoName)
	cloneDir := MakeTemporaryDirectory()


	if repoURL == "" || cloneDir == "" {
		log.Error("Missing required environment variables: GITHUB_REPOSITORY_URL or GITHUB_CLONE_DIR")
		return
	}

	cmd := exec.Command("git", "clone", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Errorf("Failed to clone repository: %v", err)
	} else {
		log.Infof("Repository cloned to %s", cloneDir)
	}
}

func MakeTemporaryDirectory() string {
	tempDir, err := os.MkdirTemp("", "github-action-")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %v", err)
	}
	log.Infof("Temporary directory created: %s", tempDir)
	return tempDir
}
