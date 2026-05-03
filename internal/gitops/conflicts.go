package gitops

import (
	"strings"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// ShowConflicts checks for unmerged files and displays them.
func ShowConflicts() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	out, err := system.GitOutput("diff", "--name-only", "--diff-filter=U")
	if err != nil {
		ui.Error("Failed to check for conflicts")
		return false
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || lines[0] == "" {
		ui.Success("No merge conflicts detected")
		return true
	}

	ui.BoxHeader("Merge Conflicts Detected")
	ui.Warn("The following files have unresolved conflicts:")
	for _, f := range lines {
		ui.Info(" ! " + f)
	}
	ui.Divider()
	ui.Info("Resolve the conflicts in your editor, then add the files and commit.")
	ui.Info("Or run 'git merge --abort' / 'git rebase --abort' to cancel.")
	
	return true
}
