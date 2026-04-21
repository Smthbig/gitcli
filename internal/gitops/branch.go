package gitops

import (
	"fmt"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func saveBranchSelection(name string) {
	cfg := config.Load()
	cfg.Branch = name
	config.Save(cfg)
}

func saveRemoteSelection(name string) {
	cfg := config.Load()
	cfg.Remote = name
	config.Save(cfg)
}

/*
SwitchBranch switches to an existing local branch safely.
*/
func SwitchBranch() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	branches, err := system.LocalBranches()
	if err != nil {
		ui.Error("Failed to list local branches")
		ui.Info(err.Error())
		return false
	}

	if len(branches) == 0 {
		ui.Warn("No local branches exist yet")
		ui.Info("Create a branch first")
		return false
	}

	current := system.CurrentGitBranch()
	options := make([]string, 0, len(branches)+1)
	for _, branch := range branches {
		label := branch
		if branch == current {
			label += " (current)"
		}
		options = append(options, label)
	}
	options = append(options, "Cancel")

	choice := ui.Select("Select existing branch", options)
	if choice == len(options) {
		ui.Warn("Branch switch cancelled")
		return false
	}

	name := branches[choice-1]
	if err := system.SwitchToBranch(name); err != nil {
		ui.Error("Failed to switch branch")
		ui.Info(err.Error())
		return false
	}

	saveBranchSelection(name)
	ui.Success("Switched to branch: " + name)
	return true
}

/*
CreateBranch creates a new branch and switches to it.
*/
func CreateBranch() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	cfg := config.Load()
	name := strings.TrimSpace(ui.InputDefault("New branch name", cfg.Branch))
	if name == "" {
		ui.Error("Branch name cannot be empty")
		return false
	}

	if system.HasLocalBranch(name) {
		ui.Warn("Branch already exists: " + name)
		if !ui.ConfirmDefault("Switch to that branch instead?", true) {
			return false
		}
		if err := system.SwitchToBranch(name); err != nil {
			ui.Error("Failed to switch branch")
			ui.Info(err.Error())
			return false
		}
		saveBranchSelection(name)
		ui.Success("Switched to branch: " + name)
		return true
	}

	if err := system.CreateBranch(name); err != nil {
		ui.Error("Failed to create branch")
		ui.Info(err.Error())
		return false
	}

	saveBranchSelection(name)
	ui.Success("Created and switched to branch: " + name)
	return true
}

/*
SwitchRemote selects an existing remote or safely creates/updates one.
*/
func SwitchRemote() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	cfg := config.Load()
	defaultRemote := cfg.Remote
	if defaultRemote == "" {
		defaultRemote = "origin"
	}

	remotes, err := system.RemoteNames()
	if err != nil {
		ui.Error("Failed to list remotes")
		ui.Info(err.Error())
		return false
	}

	if len(remotes) == 0 {
		ui.Warn("No remotes configured yet")
		name := strings.TrimSpace(ui.InputDefault("Remote name", defaultRemote))
		url := strings.TrimSpace(ui.Input("Remote URL"))
		if name == "" || url == "" {
			ui.Error("Remote name and URL are required")
			return false
		}

		if err := system.EnsureRemote(name, url); err != nil {
			ui.Error("Failed to configure remote")
			ui.Info(err.Error())
			return false
		}

		saveRemoteSelection(name)
		ui.Success("Remote added: " + name)
		return true
	}

	options := make([]string, 0, len(remotes)+2)
	for _, remote := range remotes {
		label := remote
		if url, err := system.RemoteURL(remote); err == nil && url != "" {
			label = fmt.Sprintf("%s -> %s", remote, url)
		}
		if remote == cfg.Remote {
			label += " (active)"
		}
		options = append(options, label)
	}
	options = append(options, "Create new remote", "Cancel")

	choice := ui.Select("Select remote", options)
	switch {
	case choice == len(options):
		ui.Warn("Remote update cancelled")
		return false
	case choice == len(options)-1:
		name := strings.TrimSpace(ui.InputDefault("Remote name", defaultRemote))
		url := strings.TrimSpace(ui.Input("Remote URL"))
		if name == "" || url == "" {
			ui.Error("Remote name and URL are required")
			return false
		}
		if err := system.EnsureRemote(name, url); err != nil {
			ui.Error("Failed to configure remote")
			ui.Info(err.Error())
			return false
		}
		saveRemoteSelection(name)
		ui.Success("Remote configured: " + name)
		return true
	default:
		name := remotes[choice-1]
		currentURL, _ := system.RemoteURL(name)
		if currentURL != "" {
			ui.Info("Selected remote URL: " + currentURL)
		}

		if ui.Confirm("Update this remote URL?") {
			url := strings.TrimSpace(ui.InputDefault("Remote URL", currentURL))
			if url == "" {
				ui.Error("Remote URL cannot be empty")
				return false
			}
			if err := system.EnsureRemote(name, url); err != nil {
				ui.Error("Failed to update remote")
				ui.Info(err.Error())
				return false
			}
			ui.Success("Remote updated: " + name)
		}

		saveRemoteSelection(name)
		ui.Success("Active remote set to: " + name)
		return true
	}
}
