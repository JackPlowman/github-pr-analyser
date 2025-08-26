package main

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-github/v74/github"
	"github.com/stretchr/testify/assert"
)

func TestGitHubActionSummary(t *testing.T) {
	// Arrange
	os.Setenv("RUNNING_IN_GITHUB_ACTION", "true")
	os.Setenv("GITHUB_STEP_SUMMARY", "summary.txt")
	// Act
	GitHubActionSummary()
	// Assert
	content, err := os.ReadFile("summary.txt")
	assert.NoError(t, err)
	assert.Equal(
		t,
		"# Hello World",
		string(content),
		"The content of the file should be '# Hello World'",
	)
	// Clean up
	defer os.Remove("summary.txt")
	os.Unsetenv("RUNNING_IN_GITHUB_ACTION")
	os.Unsetenv("GITHUB_STEP_SUMMARY")
}

func TestGetLanguageFromExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"main.py", "Python"},
		{"README.md", "Markdown"},
		{"document.tex", "TeX"},
		{"index.html", "HTML"},
		{"script.js", "JavaScript"},
		{"app.ts", "TypeScript"},
		{"main.go", "Go"},
		{"App.java", "Java"},
		{"program.cpp", "C++"},
		{"code.c", "C"},
		{"Program.cs", "C#"},
		{"script.php", "PHP"},
		{"app.rb", "Ruby"},
		{"main.rs", "Rust"},
		{"script.sh", "Shell"},
		{"config.yaml", "YAML"},
		{"data.json", "JSON"},
		{"styles.css", "CSS"},
		{"unknown.xyz", "Other"},
		{"no-extension", "Other"},
	}

	for _, test := range tests {
		t.Run(test.filename, func(t *testing.T) {
			result := getLanguageFromExtension(test.filename)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestAnalyzeFileTypes(t *testing.T) {
	// Create mock files
	files := []*github.CommitFile{
		{Filename: stringPtr("main.py")},
		{Filename: stringPtr("test.py")},
		{Filename: stringPtr("util.py")},
		{Filename: stringPtr("README.md")},
		{Filename: stringPtr("docs.md")},
		{Filename: stringPtr("style.css")},
		{Filename: stringPtr("script.js")},
		{Filename: stringPtr("config.json")},
		{Filename: stringPtr("unknown.xyz")},
		{Filename: stringPtr("no-ext")},
	}

	stats := analyzeFileTypes(files)

	// Verify results - should have 6 unique languages: Python, Markdown, CSS, JavaScript, JSON, Other
	assert.Len(t, stats, 6)

	// Find Python stats (should be first due to sorting by count)
	pythonStats := stats[0]
	assert.Equal(t, "Python", pythonStats.Language)
	assert.Equal(t, 3, pythonStats.Count)
	assert.Equal(t, float64(30), pythonStats.Percentage)

	// Find Other stats
	var otherStats FileStats
	for _, stat := range stats {
		if stat.Language == "Other" {
			otherStats = stat
			break
		}
	}
	assert.Equal(t, "Other", otherStats.Language)
	assert.Equal(t, 2, otherStats.Count)
	assert.Equal(t, float64(20), otherStats.Percentage)
}

func TestFormatFileStats(t *testing.T) {
	stats := []FileStats{
		{Language: "Python", Count: 5, Percentage: 50.0},
		{Language: "Markdown", Count: 3, Percentage: 30.0},
		{Language: "JavaScript", Count: 2, Percentage: 20.0},
	}

	result := formatFileStats(stats, 10)

	// Check that the output contains expected elements
	assert.Contains(t, result, "Files changed: 10")
	assert.Contains(t, result, "Python")
	assert.Contains(t, result, "5 files changed")
	assert.Contains(t, result, "50 %")
	assert.Contains(t, result, "Markdown")
	assert.Contains(t, result, "3 files changed")
	assert.Contains(t, result, "30 %")
	assert.Contains(t, result, "JavaScript")
	assert.Contains(t, result, "2 files changed")
	assert.Contains(t, result, "20 %")

	// Check that it's wrapped in markdown code block
	assert.True(t, strings.HasPrefix(result, "## Pull Request Change Statistics\n\n```markdown\n"))
	assert.True(t, strings.HasSuffix(result, "```"))

	// Check progress bars are present (50% = 12-13 '>' characters)
	assert.Contains(t, result, ">>>>>>>>>>>>") // 50% should have 12+ >
	assert.Contains(t, result, "-------")      // Should have some -
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
