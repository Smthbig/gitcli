package setup

import (
	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Run performs the full setup / reconfigure flow.
It is the main entry point for configuring Git Genius.
*/
func Run() bool {
	ui.Clear()
	ui.BoxHeader("Git Genius Setup")

	cfg := config.Load()

	if !selectWorkDir(&cfg) {
		return false
	}
	config.Save(cfg)

	if err := system.EnsureGitInstalled(); err != nil {
		ui.Error("Git is required for setup")
		return false
	}

	if !system.EnsureGitRepo() {
		return false
	}

	system.EnsureSafeDirectory(cfg.GetWorkDir())

	if !system.EnsureBranchSync() {
		return false
	}

	if !ensureGitIdentity(cfg.GetWorkDir()) {
		return false
	}

	setupGitBasics(&cfg)

	if !ensureConfiguredBranch(&cfg) {
		return false
	}

	if !setupRepo(&cfg) {
		return false
	}

	if !setupGitHubToken() {
		return false
	}

	_ = configureGitAuthIfNeeded()

	ensureGitHubRepo(&cfg)

	if err := configureRemote(&cfg); err != nil {
		ui.Error("Failed to configure git remote")
		return false
	}

	offerFirstPush(&cfg)
	config.Save(cfg)

	ui.BoxHeader("Setup Summary")
	ui.Success("Setup completed successfully")
	return true
}
