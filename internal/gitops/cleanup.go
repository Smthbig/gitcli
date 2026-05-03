package gitops

import (
	"strings"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// CleanupMergedBranches identifies and deletes local branches merged into HEAD.
func CleanupMergedBranches() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ui.BoxHeader("Smart Branch Cleanup")

	// Get branches merged into current HEAD
	out, err := system.GitOutput("branch", "--merged")
	if err != nil {
		ui.Error("Failed to list merged branches")
		return false
	}

	current := system.CurrentGitBranch()
	lines := strings.Split(out, "\n")
	var toDelete []string

	for _, line := range lines {
		name := strings.TrimSpace(strings.Replace(line, "*", "", 1))
		if name == "" || name == current || name == "main" || name == "master" || name == "develop" {
			continue
		}
		toDelete = append(toDelete, name)
	}

	if len(toDelete) == 0 {
		ui.Success("No merged branches to clean up")
		return true
	}

	ui.Info("The following branches are already merged into " + current + ":")
	for _, b := range toDelete {
		ui.Warn(" - " + b)
	}

	if !ui.Confirm("Delete these " + string(rune(len(toDelete)+'0')) + " branches?") {
		ui.Info("Cleanup cancelled")
		return false
	}

	for _, b := range toDelete {
		if err := system.RunGit("branch", "-d", b); err != nil {
			ui.Error("Failed to delete branch " + b)
		} else {
			ui.Success("Deleted " + b)
		}
	}

	return true
}
