package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v74/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func main() {
	content := generatePRFileAnalysis()
	updatePullRequestDescription(content)
	GitHubActionSummary()
}

// initialize global logger
func init() {
	logger, err := initLogger()
	zap.ReplaceGlobals(zap.Must(logger, err))
}

// initLogger initializes and returns a zap logger according to the
// DEBUG environment variable. If DEBUG=="true" a development logger
// will be returned, otherwise a production logger is used.
func initLogger() (*zap.Logger, error) {
	if os.Getenv("INPUT_DEBUG") == "true" {
		return zap.NewDevelopment()
	}
	return zap.NewProduction()
}

// FileStats represents statistics for a file type
type FileStats struct {
	Language   string
	Count      int
	Percentage float64
}

// generatePRFileAnalysis generates file change analysis for the PR
func generatePRFileAnalysis() string {
	files := getPRFiles()
	if len(files) == 0 {
		return "```markdown\nNo files changed in this PR.\n```"
	}

	stats := analyzeFileTypes(files)
	return formatFileStats(stats, len(files))
}

// getPRFiles retrieves the list of changed files in the PR
func getPRFiles() []*github.CommitFile {
	owner := os.Getenv("INPUT_GITHUB_REPOSITORY_OWNER")
	fullRepoName := os.Getenv("INPUT_GITHUB_REPOSITORY")
	prNumberStr := os.Getenv("INPUT_GITHUB_PR_NUMBER")
	token := os.Getenv("INPUT_GITHUB_TOKEN")

	if owner == "" || fullRepoName == "" || prNumberStr == "" || token == "" {
		zap.L().Error("Missing required GitHub environment variables")
		return nil
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		zap.L().Error("Invalid GITHUB_REPOSITORY format")
		return nil
	}

	repo := repoParts[1]
	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		zap.L().Error("Invalid PR number", zap.Error(err))
		return nil
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get PR files
	files, _, err := client.PullRequests.ListFiles(ctx, owner, repo, prNumber, nil)
	if err != nil {
		zap.L().Error("Failed to get PR files", zap.Error(err))
		return nil
	}

	return files
}

// analyzeFileTypes analyzes file types and returns statistics
func analyzeFileTypes(files []*github.CommitFile) []FileStats {
	languageMap := make(map[string]int)
	totalFiles := len(files)

	for _, file := range files {
		if file.Filename == nil {
			continue
		}
		language := getLanguageFromExtension(*file.Filename)
		languageMap[language]++
	}

	var stats []FileStats
	for language, count := range languageMap {
		percentage := float64(count) / float64(totalFiles) * 100
		stats = append(stats, FileStats{
			Language:   language,
			Count:      count,
			Percentage: percentage,
		})
	}

	// Sort by count descending
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].Count > stats[j].Count
	})

	return stats
}

// getLanguageFromExtension maps file extensions to language names
func getLanguageFromExtension(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	extensionMap := map[string]string{
		".py":    "Python",
		".md":    "Markdown",
		".tex":   "TeX",
		".html":  "HTML",
		".htm":   "HTML",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".go":    "Go",
		".java":  "Java",
		".cpp":   "C++",
		".c":     "C",
		".cs":    "C#",
		".php":   "PHP",
		".rb":    "Ruby",
		".rs":    "Rust",
		".sh":    "Shell",
		".yaml":  "YAML",
		".yml":   "YAML",
		".json":  "JSON",
		".xml":   "XML",
		".css":   "CSS",
		".scss":  "SCSS",
		".sass":  "Sass",
		".sql":   "SQL",
		".r":     "R",
		".kt":    "Kotlin",
		".swift": "Swift",
		".dart":  "Dart",
		".vue":   "Vue",
		".jsx":   "JSX",
		".tsx":   "TSX",
	}

	if language, exists := extensionMap[ext]; exists {
		return language
	}
	return "Other"
}

// formatFileStats formats the file statistics into markdown
func formatFileStats(stats []FileStats, totalFiles int) string {
	var output strings.Builder
	levelStr := os.Getenv("INPUT_HEADING_LEVEL")
	level, err := strconv.Atoi(levelStr)
	if err != nil {
		zap.L().Warn("Invalid heading level, defaulting to 2", zap.String("level", levelStr))
		level = 2
	}
	output.WriteString(strings.Repeat("#", level) + " Pull Request Change Statistics\n\n")
	output.WriteString("```markdown\n")
	output.WriteString(fmt.Sprintf("Files changed: %d\n\n", totalFiles))

	for _, stat := range stats {
		// Create progress bar (25 characters total)
		filled := int(stat.Percentage / 4) // Each character represents 4%
		if filled > 25 {
			filled = 25
		}

		progressBar := strings.Repeat(">", filled) + strings.Repeat("-", 25-filled)

		output.WriteString(fmt.Sprintf("%-18s %2d files changed %s   %02.0f %%\n",
			stat.Language,
			stat.Count,
			progressBar,
			stat.Percentage,
		))
	}

	output.WriteString("```")
	return output.String()
}

// update PR description
func updatePullRequestDescription(content string) {
	// Get required environment variables
	owner := os.Getenv("INPUT_GITHUB_REPOSITORY_OWNER")
	fullRepoName := os.Getenv("INPUT_GITHUB_REPOSITORY") // Expected format: "owner/repo"
	prNumberStr := os.Getenv("INPUT_GITHUB_PR_NUMBER")
	token := os.Getenv("INPUT_GITHUB_TOKEN")

	if owner == "" || fullRepoName == "" || prNumberStr == "" || token == "" {
		zap.L().Error(
			"Missing required GitHub environment variables: INPUT_GITHUB_REPOSITORY_OWNER, INPUT_GITHUB_REPOSITORY, INPUT_GITHUB_PR_NUMBER, or INPUT_GITHUB_TOKEN",
		)
		return
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		zap.L().Error(
			"INPUT_GITHUB_REPOSITORY environment variable is not in the expected 'owner/repo' format",
			zap.String("repository", fullRepoName),
		)
		return
	}
	repo := repoParts[1]

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		zap.L().Error("Invalid PR number", zap.Error(err))
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
		zap.L().Error("Failed to get PR", zap.Error(err))
		return
	}

	// Replace the marker line in the PR description with a string
	body := pr.GetBody()
	marker := "<!-- github-pr-analyser-replace-line -->"
	replacement := content

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
		zap.L().Warn("Marker not found in PR description. No update performed.")
		return
	}

	updatedBody := strings.Join(lines, "\n")
	pr.Body = &updatedBody

	_, _, err = client.PullRequests.Edit(ctx, owner, repo, prNumber, pr)
	if err != nil {
		zap.L().Error("Failed to update PR description", zap.Error(err))
		return
	}
	zap.L().Info("PR description updated successfully.")
}

// Generate GitHub Action Summary
func GitHubActionSummary() {
	action := os.Getenv("RUNNING_IN_GITHUB_ACTION")
	if action == "true" {
		zap.L().Info("Running in GitHub Action, Generating Summary")
		gitHubActionSummaryFile := os.Getenv("GITHUB_STEP_SUMMARY")
		content := []byte("# Hello World")
		err := os.WriteFile(gitHubActionSummaryFile, content, 0o600)
		if err != nil {
			panic(err)
		}
		zap.L().Info("Summary Generated")
	}
}

// AddPullRequestComment adds a comment to the pull request using GitHub API
func AddPullRequestComment(comment string) {
	owner := os.Getenv("INPUT_GITHUB_REPOSITORY_OWNER")
	fullRepoName := os.Getenv("INPUT_GITHUB_REPOSITORY") // Expected format: "owner/repo"
	prNumberStr := os.Getenv("INPUT_GITHUB_PR_NUMBER")
	token := os.Getenv("INPUT_GITHUB_TOKEN")

	if owner == "" || fullRepoName == "" || prNumberStr == "" {
		zap.L().Error(
			"Missing required GitHub environment variables: INPUT_GITHUB_REPOSITORY_OWNER, INPUT_GITHUB_REPOSITORY, or INPUT_GITHUB_PR_NUMBER",
		)
		return
	}

	repoParts := strings.Split(fullRepoName, "/")
	if len(repoParts) != 2 {
		zap.L().Error(
			"INPUT_GITHUB_REPOSITORY environment variable is not in the expected 'owner/repo' format",
			zap.String("repository", fullRepoName),
		)
		return
	}
	repo := repoParts[1] // Use only the repository name

	prNumber, err := strconv.Atoi(prNumberStr)
	if err != nil {
		zap.L().Error("Error converting INPUT_GITHUB_PR_NUMBER to integer", zap.String("prNumber", prNumberStr), zap.Error(err))
		return
	}

	if token == "" {
		zap.L().Error(
			"INPUT_GITHUB_TOKEN environment variable is not set. Cannot authenticate to GitHub.",
		)
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
		zap.L().Error("Error adding comment to pull request", zap.Error(err))
	} else {
		zap.L().Info("Successfully added comment to PR", zap.Int("prNumber", prNumber), zap.String("owner", owner), zap.String("repo", repo))
	}
}
