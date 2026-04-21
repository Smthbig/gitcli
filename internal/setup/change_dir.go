package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/gitops"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func ChangeProjectDir() bool {
	return switchProject(false)
}

func SwitchProjectRepo() bool {
	return switchProject(true)
}

func switchProject(withRepoOptions bool) bool {
	title := "Change Project Directory"
	if withRepoOptions {
		title = "Switch Project / Repo"
	}

	ui.Clear()
	ui.Header(title)

	current := config.Load().GetWorkDir()
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

		cfg, isRepo, err := activateProjectDir(resolveRecentDir(dir, recent))
		if err != nil {
			ui.Error(err.Error())
			if !ui.ConfirmDefault("Try another directory?", true) {
				return false
			}
			continue
		}

		ui.Success("Project directory updated")
		ui.Info("New project directory:")
		ui.Info(cfg.WorkDir)

		if !isRepo {
			return handleNonRepoSelection(cfg)
		}

		ui.Success("Git repository detected in new directory")
		bootstrapProjectConfig(&cfg)

		if !system.EnsureBranchSync() {
			return false
		}

		if withRepoOptions {
			return offerRepoSwitchOptions(cfg)
		}
		return true
	}
}

func resolveRecentDir(input string, recent []string) string {
	if len(input) == 1 && input[0] >= '1' && input[0] <= '5' {
		idx := int(input[0] - '1')
		if idx >= 0 && idx < len(recent) {
			return recent[idx]
		}
	}
	return input
}

func activateProjectDir(dir string) (config.Config, bool, error) {
	abs, err := filepath.Abs(strings.TrimSpace(dir))
	if err != nil {
		return config.Config{}, false, fmt.Errorf("failed to resolve directory path")
	}

	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return config.Config{}, false, fmt.Errorf("invalid directory path")
	}

	config.SetActiveWorkDir(abs)

	if system.IsGitRepo() {
		system.EnsureSafeDirectory(abs)
		cfg := config.Load()
		cfg.WorkDir = abs
		return cfg, true, nil
	}

	cfg := config.LoadForWorkDir(abs)
	cfg.WorkDir = abs
	return cfg, false, nil
}

func bootstrapProjectConfig(cfg *config.Config) {
	if cfg == nil {
		return
	}

	changed := false

	if !config.HasProjectConfig(cfg.WorkDir) {
		if branch := system.CurrentGitBranch(); branch != "" {
			cfg.Branch = branch
			changed = true
		}

		if remotes, err := system.RemoteNames(); err == nil && len(remotes) == 1 {
			cfg.Remote = remotes[0]
			changed = true
		}
	}

	if !changed {
		return
	}

	config.Save(*cfg)
	ui.Info("Loaded repo defaults for this project so setup does not need to be repeated")
}

func handleNonRepoSelection(cfg config.Config) bool {
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
	bootstrapProjectConfig(&cfg)

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

func offerRepoSwitchOptions(cfg config.Config) bool {
	ui.Info("Saved branch : " + cfg.Branch)
	ui.Info("Saved remote : " + cfg.Remote)

	if ui.Confirm("Switch branch in this repo now?") {
		if !gitops.SwitchBranch() {
			return false
		}
	}

	if ui.Confirm("Configure or switch the active remote now?") {
		if !gitops.SwitchRemote() {
			return false
		}
	}

	if !config.HasProjectConfig(cfg.WorkDir) {
		ui.Info("This repo is using auto-detected defaults right now")
		ui.Info("Run Tools -> Setup / Reconfigure later if you want to save owner/repo details too")
	}

	return true
}
