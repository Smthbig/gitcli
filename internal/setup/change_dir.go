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

/*
SwitchProject is the unified entry point for changing the active project directory.
It detects if the target is a Git repo and shows context before switching.
*/
func SwitchProject() bool {
	ui.Clear()
	ui.Header("Switch Project")

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

		targetDir := resolveRecentDir(dir, recent)
		abs, err := filepath.Abs(targetDir)
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

		// Load config for the target directory to show context
		targetCfg := config.LoadForWorkDir(abs)
		isRepo := system.IsGitRepoAt(abs)

		ui.Divider()
		ui.Info("Target Directory: " + abs)
		if isRepo {
			ui.Success("Git repository detected")
			showRepoContext(abs, targetCfg)
		} else {
			ui.Warn("Not a git repository")
		}
		ui.Divider()

		if !ui.ConfirmDefault("Switch to this project?", true) {
			if ui.Confirm("Try another directory?") {
				continue
			}
			return false
		}

		// Perform the switch
		config.SetActiveWorkDir(abs)
		ui.Success("Project directory updated")

		if !isRepo {
			return handleNonRepoSelection(targetCfg)
		}

		// Sync branch/remote defaults if it's a new repo for Git Genius
		bootstrapProjectConfig(&targetCfg)

		// Offer post-switch options if Git is available
		if system.CommandExists("git") {
			if !system.EnsureBranchSync() {
				return false
			}
			return offerRepoSwitchOptions(targetCfg)
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

func showRepoContext(dir string, cfg config.Config) {
	// We need to temporarily set the work dir to inspect it, or use GitOutputAt
	branch, _ := system.GitOutputAt(dir, "branch", "--show-current")
	if branch == "" {
		branch = cfg.Branch
	}
	if branch == "" {
		branch = "-"
	}

	remote := cfg.Remote
	if remote == "" {
		// Try to auto-detect a remote if not in config
		if remotes, err := system.RemoteNames(); err == nil && len(remotes) > 0 {
			remote = remotes[0]
		}
	}
	if remote == "" {
		remote = "-"
	}

	ui.Info("Branch : " + branch)
	ui.Info("Remote : " + remote)
	if cfg.Owner != "" && cfg.Repo != "" {
		ui.Info("Repo   : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	}
}

func bootstrapProjectConfig(cfg *config.Config) {
	if cfg == nil {
		return
	}

	changed := false

	if !config.HasProjectConfig(cfg.WorkDir) {
		if branch, _ := system.GitOutputAt(cfg.WorkDir, "branch", "--show-current"); branch != "" {
			cfg.Branch = branch
			changed = true
		}

		if remotes, err := system.GitOutputAt(cfg.WorkDir, "remote"); err == nil && remotes != "" {
			lines := strings.Split(remotes, "\n")
			if len(lines) > 0 && lines[0] != "" {
				cfg.Remote = strings.TrimSpace(lines[0])
				changed = true
			}
		}
	}

	if !changed {
		return
	}

	config.Save(*cfg)
	ui.Info("Loaded repo defaults for this project")
}

func handleNonRepoSelection(cfg config.Config) bool {
	if !system.CommandExists("git") {
		ui.Warn("Git is not installed. Operations will be limited.")
		return true
	}

	if !ui.Confirm("Initialize a git repository here?") {
		ui.Warn("Git operations will be limited until repo is initialized")
		return true
	}

	if err := system.RunGitAt(cfg.WorkDir, "init"); err != nil {
		ui.Error("Failed to initialize git repository")
		return false
	}

	ui.Success("Git repository initialized")
	bootstrapProjectConfig(&cfg)

	if cfg.Branch != "" {
		// We can't use PrepareBranch easily here because it uses config.Load()
		// but we just initialized a repo. Let's just try to create/checkout.
		_ = system.RunGitAt(cfg.WorkDir, "checkout", "-b", cfg.Branch)
	}

	return true
}

func offerRepoSwitchOptions(cfg config.Config) bool {
	ui.Info("Current branch : " + cfg.Branch)
	ui.Info("Current remote : " + cfg.Remote)

	if ui.Confirm("Switch branch now?") {
		if !gitops.SwitchBranch() {
			return false
		}
	}

	if ui.Confirm("Configure or switch remote now?") {
		if !gitops.SwitchRemote() {
			return false
		}
	}

	return true
}

// Deprecated: used for backward compatibility during refactor
func ChangeProjectDir() bool {
	return SwitchProject()
}

// Deprecated: used for backward compatibility during refactor
func SwitchProjectRepo() bool {
	return SwitchProject()
}
