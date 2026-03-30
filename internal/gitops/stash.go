package gitops

import (
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
StashSave saves current working tree changes
*/
func StashSave() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	msg := ui.Input("Stash message (optional)")
	args := []string{"stash", "push"}

	if msg != "" {
		args = append(args, "-m", msg)
	}

	if err := system.RunGit(args...); err != nil {
		ui.Error("Failed to stash changes")
		return false
	}

	ui.Success("Changes stashed successfully")
	return true
}

/*
StashList shows all stashes
*/
func StashList() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	if err := system.RunGit("stash", "list"); err != nil {
		ui.Error("Failed to list stashes")
		return false
	}
	return true
}

/*
StashPop applies and removes latest stash
*/
func StashPop() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	if err := system.RunGit("stash", "pop"); err != nil {
		ui.Error("Failed to apply stash")
		return false
	}

	ui.Success("Stash applied successfully")
	return true
}
