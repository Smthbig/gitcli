package gitops

import (
	"fmt"
	"strings"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// ShowRepoStats displays high-level repository statistics.
func ShowRepoStats() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ui.BoxHeader("Repository Statistics")

	// Total commits
	total, _ := system.GitOutput("rev-list", "--count", "HEAD")
	if total == "" {
		total = "0"
	}
	ui.Info("Total Commits    : " + total)

	// Last commit date
	lastDate, _ := system.GitOutput("log", "-1", "--format=%cd", "--date=relative")
	if lastDate != "" {
		ui.Info("Last Activity    : " + lastDate)
	}

	// Contributors count
	authors, _ := system.GitOutput("shortlog", "-s", "-n", "HEAD")
	authorLines := strings.Split(strings.TrimSpace(authors), "\n")
	ui.Info(fmt.Sprintf("Contributors     : %d", len(authorLines)))

	// Top contributor
	if len(authorLines) > 0 && authorLines[0] != "" {
		top := strings.TrimSpace(authorLines[0])
		ui.Info("Top Contributor  : " + top)
	}

	// Branch count
	branches, _ := system.GitOutput("branch", "-a")
	branchLines := strings.Split(strings.TrimSpace(branches), "\n")
	ui.Info(fmt.Sprintf("Total Branches   : %d", len(branchLines)))

	return true
}
