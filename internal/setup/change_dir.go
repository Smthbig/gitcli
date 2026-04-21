package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
ChangeProjectDir switches the active project directory safely.
*/
func ChangeProjectDir() bool {
	ui.Clear()
	ui.Header("Change Project Directory")

	cfg := config.Load()
	current := cfg.GetWorkDir()
	ui.Info("Current project directory:")
	ui.Info(current)

	recent := config.RecentWorkDirs()
	if len(recent) > 0 {
		ui.Info("Recent project directories:")
		max := len(recent)
		if max > 5 {
			max = 5
		}
		for i := 0; i < max; i++ {
			ui.Info(fmt.Sprintf("  %d) %s", i+1, recent[i]))
		}
	}

	for {
		dir := strings.TrimSpace(ui.Input("Enter full path of new project directory (or 1-5 for recent)"))
		if dir == "" {
			ui.Error("Directory path cannot be empty")
			return false
		}

		if len(dir) == 1 && dir[0] >= '1' && dir[0] <= '5' {
			idx := int(dir[0] - '1')
			if idx >= 0 && idx < len(recent) {
				dir = recent[idx]
			}
		}

		abs, err := filepath.Abs(dir)
		if err != nil {
			ui.Error("Failed to resolve directory path")
			return false
		}

		info, err := os.Stat(abs)
		if err != nil || !info.IsDir() {
			ui.Error("Invalid directory path")
			if !ui.ConfirmDefault("Try another directory?", true) {
				return false
			}
			continue
		}

		cfg.WorkDir = abs
		config.Save(cfg)

		ui.Success("Project directory updated")
		ui.Info("New project directory:")
		ui.Info(abs)

		if system.IsGitRepo() {
			ui.Success("Git repository detected in new directory")
			if !system.EnsureBranchSync() {
				return false
			}
			return true
		}

		ui.Warn("Selected directory is not a git repository")
		if !ui.Confirm("Initialize a git repository here?") {
			ui.Warn("Git operations will be limited until repo is initialized")
			return true
		}

		if err := system.RunGit("init"); err != nil {
			ui.Error("Failed to initialize git repository")
			return false
		}

		ui.Success("Git repository initialized")
		if cfg.Branch != "" {
			if err := system.PrepareBranch(cfg.Branch); err != nil {
				ui.Warn("Could not prepare configured branch")
				ui.Info(err.Error())
			} else {
				ui.Success("Branch ready: " + cfg.Branch)
			}
		}

		return true
	}
}
