package setup

import (
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/gitops"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func setupGitBasics(cfg *config.Config) {
	cfg.Branch = strings.TrimSpace(ui.InputDefault("Default branch", cfg.Branch))
	cfg.Remote = strings.TrimSpace(ui.InputDefault("Remote name", cfg.Remote))
}

func ensureConfiguredBranch(cfg *config.Config) bool {
	desired := cfg.Branch
	current := system.CurrentGitBranch()
	if current == desired {
		return true
	}
	if current == "" {
		_ = system.PrepareBranch(desired)
		return true
	}
	if system.HasLocalBranch(desired) {
		if ui.Confirm("Switch to the configured branch?") {
			_ = system.SwitchToBranch(desired)
		}
	} else if ui.Confirm("Create and switch to the configured branch?") {
		_ = system.CreateBranch(desired)
	}
	return true
}

func setupRepo(cfg *config.Config) bool {
	ui.BoxHeader("GitHub Repository")
	defaultRepo := filepath.Base(cfg.GetWorkDir())
	cfg.Owner = strings.TrimSpace(ui.InputDefault("GitHub owner", cfg.Owner))
	cfg.Repo = strings.TrimSpace(ui.InputDefault("Repository name", defaultRepo))
	return true
}

func offerFirstPush(cfg *config.Config) {
	if ui.Confirm("Push current code to GitHub now?") {
		_ = gitops.Push("Initial commit")
	}
}
