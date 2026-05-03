package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/ui"
)

func selectWorkDir(cfg *config.Config) bool {
	cwd, _ := os.Getwd()
	if cfg.WorkDir == "" {
		cfg.WorkDir = cwd
	}

	ui.Info("Current directory: " + cwd)
	recent := config.RecentWorkDirs()

	if !ui.Confirm("Use current project directory?") {
		for {
			ui.Info("Recent project directories:")
			for i, dir := range recent {
				if i >= 3 {
					break
				}
				ui.Info(fmt.Sprintf("  %d) %s", i+1, dir))
			}
			dir := strings.TrimSpace(ui.Input("Enter full path of project directory (or 1/2/3 for recent)"))
			if dir == "" {
				ui.Error("Directory path cannot be empty")
				continue
			}

			if len(dir) == 1 && dir[0] >= '1' && dir[0] <= '3' {
				idx := int(dir[0] - '1')
				if idx < len(recent) {
					dir = recent[idx]
				}
			}

			abs, err := filepath.Abs(dir)
			if err != nil {
				ui.Error("Failed to resolve directory path")
				if !ui.ConfirmDefault("Try another directory?", true) {
					return false
				}
				continue
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
			break
		}
	}

	config.SetActiveWorkDir(cfg.WorkDir)
	return true
}
