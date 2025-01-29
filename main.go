package main

import (
	"os"
<<<<<<< HEAD
	"path/filepath"
=======
	"strconv"
	"context"
>>>>>>> 72b4a8b (update)

	log "github.com/sirupsen/logrus"
)
import "github.com/google/go-github/v68/github"
func main() {
	initLogging()
<<<<<<< HEAD
	RunHello()
	CreateTempFolders()
=======
>>>>>>> 72b4a8b (update)
	GitHubActionSummary()
	AddPullRequestComment("Hello World")
}

// Init logging configuration
func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
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

<<<<<<< HEAD
// Create a temp folder with two subfolders: one for the repo default branch and one for the commit being reviewed
func CreateTempFolders() {
	tempDir, err := os.MkdirTemp("", "github-pr-analyser")
	if err != nil {
		log.Fatal("Failed to create temp directory: ", err)
	}

	defaultBranchDir := filepath.Join(tempDir, "default-branch")
	commitDir := filepath.Join(tempDir, "commit")

	err = os.Mkdir(defaultBranchDir, 0755)
	if err != nil {
		log.Fatal("Failed to create default branch directory: ", err)
	}

	err = os.Mkdir(commitDir, 0755)
	if err != nil {
		log.Fatal("Failed to create commit directory: ", err)
	}

	log.Infof("Created temp directories: %s, %s", defaultBranchDir, commitDir)
=======
func AddPullRequestComment(comment string) {
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	repo := os.Getenv("GITHUB_REPOSITORY")
	number := os.Getenv("GITHUB_PR_NUMBER")
	prNumber, err := strconv.Atoi(number)
	client := github.NewClient(nil)

	_, _, err = client.Issues.CreateComment(client, owner, repo, prNumber, &github.IssueComment{
		Body: &comment,
	})
	if err != nil {
		log.Error("Error adding comment to pull request: ", err)
	}
>>>>>>> 72b4a8b (update)
}
