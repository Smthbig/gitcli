package gitops

import (
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
UndoLastCommit undoes the last commit but keeps changes staged
Uses: git reset --soft HEAD~1
*/
func UndoLastCommit() {
	if !system.EnsureGitRepo() {
		return
	}

	if !hasAnyCommit() {
		ui.Warn("No commits found to undo")
		return
	}

	if !ui.Confirm("Undo last commit? (changes will be kept)") {
		ui.Warn("Undo cancelled")
		return
	}

	if err := system.RunGit("reset", "--soft", "HEAD~1"); err != nil {
		ui.Error("Failed to undo last commit")
		return
	}

	ui.Success("Last commit undone (changes preserved)")
}
