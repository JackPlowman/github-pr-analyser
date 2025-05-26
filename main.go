package main

import (
	"context"
	"os"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/go-github/v68/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func main() {
	InitLogging()
	CloneRepository()
	GitHubActionSummary()
	AddPullRequestComment("Hello World")
}

// Init logging configuration
func InitLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

}

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

// Generate GitHub Action Summary
func GitHubActionSummary() {
	action := os.Getenv("RUNNING_IN_GITHUB_ACTION")
	if action == "true" {
		log.Info("Running in GitHub Action, Generating Summary")
		gitHubActionSummaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
		content := []byte("# Hello World")
		err := os.WriteFile(gitHubActionSummaryFile, content, 0600)
		if err != nil {
			panic(err)
		}
		log.Info("Summary Generated")
	}
}

// AddPullRequestComment adds a comment to the pull request using GitHub API
func AddPullRequestComment(comment string) {
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	fullRepoName := os.Getenv("GITHUB_REPOSITORY") // Expected format: "owner/repo"
	prNumberStr := os.Getenv("GITHUB_PR_NUMBER")
	token := os.Getenv("GITHUB_TOKEN")

	if owner == "" || fullRepoName == "" || prNumberStr == "" {
		log.Error("Missing required GitHub environment variables: GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY, or GITHUB_PR_NUMBER")
		return
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		log.Errorf("GITHUB_REPOSITORY environment variable (%s) is not in the expected 'owner/repo' format.", fullRepoName)
		return
	}
	repo := repoParts[1] // Use only the repository name

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		log.Errorf("Error converting GITHUB_PR_NUMBER '%s' to integer: %v", prNumberStr, err)
		return
	}

	if token == "" {
		log.Error("GITHUB_TOKEN environment variable is not set. Cannot authenticate to GitHub.")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc) // Use authenticated client

	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{
		Body: &comment,
	})
	if err != nil {
		log.Error("Error adding comment to pull request: ", err)
	} else {
		log.Infof("Successfully added comment to PR #%d in %s/%s", prNumber, owner, repo)
	}
}
