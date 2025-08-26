package main

import (
	"os"
	"testing"

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
