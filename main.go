package main

import (
	"context" // Added import
	"os"
	"strconv"

	"github.com/google/go-github/v68/github"
	log "github.com/sirupsen/logrus"
)

func main() {
	initLogging()
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

// AddPullRequestComment adds a comment to the pull request using GitHub API
func AddPullRequestComment(comment string) {
	owner := os.Getenv("GITHUB_REPOSITORY_OWNER")
	repo := os.Getenv("GITHUB_REPOSITORY")
	number := os.Getenv("GITHUB_PR_NUMBER")
	prNumber, err := strconv.Atoi(number)
	client := github.NewClient(nil)

	_, _, err = client.Issues.CreateComment(context.Background(), owner, repo, prNumber, &github.IssueComment{ // Changed client to context.Background()
		Body: &comment,
	})
	if err != nil {
		log.Error("Error adding comment to pull request: ", err)
	}
}
