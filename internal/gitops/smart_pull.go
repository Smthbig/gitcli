package gitops

import (
	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
SmartPull performs:
1. Detect dirty working tree
2. Auto-stash changes (optional)
3. Pull latest changes
4. Restore stash (if created)

This prevents pull failures due to local changes.
*/
func SmartPull() {
	if !system.EnsureGitRepo() {
		return
	}

	cfg := config.Load()
	stashed := false

	// Step 1: Detect uncommitted changes
	if isWorkingTreeDirty() {
		ui.Warn("Uncommitted changes detected")

		if !ui.Confirm("Auto-stash changes and continue pull?") {
			ui.Warn("Smart pull cancelled")
			return
		}

		if err := system.RunGit("stash", "push", "-m", "git-genius-auto-stash"); err != nil {
			ui.Error("Failed to auto-stash changes")
			return
		}

		stashed = true
		ui.Success("Changes stashed temporarily")
	}

	// Step 2: Pull latest changes
	ui.Info("Pulling latest changes...")
	if err := system.RunGit("pull", cfg.Remote, cfg.Branch); err != nil {
		ui.Error("Pull failed")

		// Try restoring stash if pull failed
		if stashed {
			_ = system.RunGit("stash", "pop")
		}
		return
	}

	// Step 3: Restore stash if created
	if stashed {
		ui.Info("Restoring stashed changes...")
		if err := system.RunGit("stash", "pop"); err != nil {
			ui.Warn("Auto-stash could not be applied cleanly")
			ui.Info("Resolve conflicts manually if needed")
			return
		}
		ui.Success("Stashed changes restored")
	}

	ui.Success("Smart pull completed successfully")
}
