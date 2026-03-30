package setup

import (
	"path/filepath"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
	"os"
)

/*
ChangeProjectDir allows switching the active project directory safely
Used from: Tools → Change Project Directory
*/
func ChangeProjectDir() {
	ui.Clear()
	ui.Header("Change Project Directory")

	cfg := config.Load()

	// Show current directory
	current := cfg.GetWorkDir()
	ui.Info("Current project directory:")
	ui.Info(current)

	// Ask new directory
	dir := ui.Input("Enter full path of NEW project directory")
	if dir == "" {
		ui.Error("Directory path cannot be empty")
		return
	}

	abs, err := filepath.Abs(dir)
	if err != nil {
		ui.Error("Failed to resolve directory path")
		return
	}

	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		ui.Error("Invalid directory path")
		return
	}

	// Save new workdir
	cfg.WorkDir = abs
	config.Save(cfg)

	ui.Success("Project directory updated")
	ui.Info("New project directory:")
	ui.Info(abs)

	// ---------------- Git repo check ----------------
	if system.IsGitRepo() {
		ui.Success("Git repository detected in new directory")

		// Sync branch config if needed
		system.EnsureBranchSync()
		return
	}

	ui.Warn("Selected directory is NOT a git repository")

	if !ui.Confirm("Initialize git repository here?") {
		ui.Warn("Git operations will be limited until repo is initialized")
		return
	}

	if err := system.RunGit("init"); err != nil {
		ui.Error("Failed to initialize git repository")
		return
	}

	ui.Success("Git repository initialized")

	// Prepare branch
	cfg = config.Load()
	if cfg.Branch != "" {
		// Avoid force-reset behavior of -B on existing repos.
		if err := system.RunGit("checkout", cfg.Branch); err != nil {
			if err := system.RunGit("checkout", "-b", cfg.Branch); err != nil {
				ui.Warn("Could not prepare configured branch: " + cfg.Branch)
				return
			}
		}
		ui.Success("Branch ready: " + cfg.Branch)
	}
}
