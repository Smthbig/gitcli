package gitops

import (
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
StashSave saves current working tree changes
*/
func StashSave() {
	if !system.EnsureGitRepo() {
		return
	}

	msg := ui.Input("Stash message (optional)")
	args := []string{"stash", "push"}

	if msg != "" {
		args = append(args, "-m", msg)
	}

	if err := system.RunGit(args...); err != nil {
		ui.Error("Failed to stash changes")
		return
	}

	ui.Success("Changes stashed successfully")
}

/*
StashList shows all stashes
*/
func StashList() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("stash", "list"); err != nil {
		ui.Error("Failed to list stashes")
	}
}

/*
StashPop applies and removes latest stash
*/
func StashPop() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("stash", "pop"); err != nil {
		ui.Error("Failed to apply stash")
		return
	}

	ui.Success("Stash applied successfully")
}
