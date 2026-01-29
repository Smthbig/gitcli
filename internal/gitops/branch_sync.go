package gitops

import (
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
EnsureBranchSync keeps config branch and actual git branch in sync
Handles master/main mismatch cleanly
*/
func EnsureBranchSync() {
	if !system.IsGitRepo() {
		return
	}

	cfg := config.Load()

	current, err := system.GitOutput("branch", "--show-current")
	if err != nil || current == "" {
		// No commits yet → nothing to sync
		return
	}

	current = strings.TrimSpace(current)

	// Already synced
	if current == cfg.Branch {
		return
	}

	ui.Warn("Branch mismatch detected")
	ui.Info("Git branch    : " + current)
	ui.Info("Config branch : " + cfg.Branch)

	if !ui.Confirm("Sync config to current git branch?") {
		ui.Warn("Keeping config branch unchanged")
		return
	}

	cfg.Branch = current
	config.Save(cfg)

	ui.Success("Branch synced to: " + current)
}
