package gitops

import (
	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
SwitchBranch switches (or creates) a branch
*/
func SwitchBranch() {
	if !system.EnsureGitRepo() {
		return
	}

	name := ui.Input("New branch name")
	if name == "" {
		ui.Error("Branch name cannot be empty")
		return
	}

	if err := system.RunGit("checkout", "-B", name); err != nil {
		ui.Error("Failed to switch branch")
		return
	}

	cfg := config.Load()
	cfg.Branch = name
	config.Save(cfg)

	ui.Success("Switched to branch: " + name)
}

/*
SwitchRemote changes git remote
*/
func SwitchRemote() {
	if !system.EnsureGitRepo() {
		return
	}

	name := ui.Input("Remote name")
	url := ui.Input("Remote URL")

	if name == "" || url == "" {
		ui.Error("Remote name and URL are required")
		return
	}

	_ = system.RunGit("remote", "remove", name)
	if err := system.RunGit("remote", "add", name, url); err != nil {
		ui.Error("Failed to add remote")
		return
	}

	cfg := config.Load()
	cfg.Remote = name
	config.Save(cfg)

	ui.Success("Remote updated: " + name)
}
