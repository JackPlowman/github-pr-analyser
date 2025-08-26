package main

import (
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/google/go-github/v74/github"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func main() {
	initLogging()
	updatePullRequestDescription()
	GitHubActionSummary()
}

// Init logging configuration
func initLogging() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
}

// update PR description
func updatePullRequestDescription() {
	// Get required environment variables
	owner := os.Getenv("INPUT_GITHUB_REPOSITORY_OWNER")
	fullRepoName := os.Getenv("INPUT_GITHUB_REPOSITORY") // Expected format: "owner/repo"
	prNumberStr := os.Getenv("INPUT_GITHUB_PR_NUMBER")
	token := os.Getenv("INPUT_GITHUB_TOKEN")

	if owner == "" || fullRepoName == "" || prNumberStr == "" || token == "" {
		log.Error(
			"Missing required GitHub environment variables: INPUT_GITHUB_REPOSITORY_OWNER, INPUT_GITHUB_REPOSITORY, INPUT_GITHUB_PR_NUMBER, or INPUT_GITHUB_TOKEN",
		)
		return
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		log.Errorf(
			"GITHUB_REPOSITORY environment variable (%s) is not in the expected 'owner/repo' format.",
			fullRepoName,
		)
		return
	}
	repo := repoParts[1]

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		log.Errorf("Invalid PR number: %v", err)
		return
	}

	// Setup GitHub client
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get current PR
	pr, _, err := client.PullRequests.Get(ctx, owner, repo, prNumber)
	if err != nil {
		log.Errorf("Failed to get PR: %v", err)
		return
	}

	// Replace the marker line in the PR description with a string
	body := pr.GetBody()
	marker := "<!-- github-pr-analyser-replace-line -->"
	replacement := "hello world" // default replacement string

	// Perform line-based replacement to ensure we replace whole marker lines
	lines := strings.Split(body, "\n")
	replaced := false
	for i, ln := range lines {
		if strings.TrimSpace(ln) == marker {
			lines[i] = replacement
			replaced = true
		}
	}

	if !replaced {
		log.Warn("Marker not found in PR description. No update performed.")
		return
	}

	updatedBody := strings.Join(lines, "\n")
	pr.Body = &updatedBody

	_, _, err = client.PullRequests.Edit(ctx, owner, repo, prNumber, pr)
	if err != nil {
		log.Errorf("Failed to update PR description: %v", err)
		return
	}
	log.Info("PR description updated successfully.")
}

// Generate GitHub Action Summary
func GitHubActionSummary() {
	action := os.Getenv("RUNNING_IN_GITHUB_ACTION")
	if action == "true" {
		log.Info("Running in GitHub Action, Generating Summary")
		gitHubActionSummaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
		content := []byte("# Hello World")
		err := os.WriteFile(gitHubActionSummaryFile, content, 0o600)
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
		log.Error(
			"Missing required GitHub environment variables: GITHUB_REPOSITORY_OWNER, GITHUB_REPOSITORY, or GITHUB_PR_NUMBER",
		)
		return
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		log.Errorf(
			"GITHUB_REPOSITORY environment variable (%s) is not in the expected 'owner/repo' format.",
			fullRepoName,
		)
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
