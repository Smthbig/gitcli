package setup

import (
	"fmt"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
CreateOrLinkRepo
Stable, token-safe, org-safe
No fake offline mode
No hard dependency on RepoExists
*/
func CreateOrLinkRepo() {

	ui.Clear()
	ui.Header("Create / Link GitHub Repository")

	cfg := config.Load()

	// --------------------------------------------------
	// Ensure local git repo
	// --------------------------------------------------
	if !system.EnsureGitRepo() {
		return
	}

	system.EnsureSafeDirectory(cfg.WorkDir)

	// --------------------------------------------------
	// Owner / Repo
	// --------------------------------------------------
	if cfg.Owner == "" {
		cfg.Owner = ui.Input("GitHub username or organisation")
	}

	if cfg.Repo == "" {
		cfg.Repo = ui.Input("Repository name")
	}

	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("Owner and repository name are required")
		return
	}

	// Default remote if empty
	if cfg.Remote == "" {
		cfg.Remote = "origin"
	}

	// --------------------------------------------------
	// Token check (no validation yet)
	// --------------------------------------------------
	token := github.GetToken()
	if token == "" {
		ui.Error("GitHub token not configured")
		ui.Info("Run Setup first to configure token")
		return
	}

	// --------------------------------------------------
	// Try checking repo existence (non-fatal)
	// --------------------------------------------------
	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)

	if err != nil {
		ui.Warn("Could not verify repository with GitHub API")
		ui.Info("Continuing anyway...")
	} else if !exists {

		ui.Warn("Repository does not exist on GitHub")

		if ui.Confirm("Create repository on GitHub now?") {

			private := ui.Confirm("Make repository PRIVATE?")

			if err := github.CreateRepo(cfg.Owner, cfg.Repo, private); err != nil {
				ui.Error("Repository creation failed")
				ui.Info(err.Error())
				ui.Info("You may create the repository manually and retry.")
				return
			}

			ui.Success("Repository created successfully")
		}
	} else {
		ui.Success("GitHub repository exists")
	}

	// --------------------------------------------------
	// Configure remote safely
	// --------------------------------------------------
	remoteURL := fmt.Sprintf(
		"https://github.com/%s/%s.git",
		cfg.Owner,
		cfg.Repo,
	)

	// Check if remote exists
	existingURL, _ := system.GitOutput("remote", "get-url", cfg.Remote)

	if existingURL != "" {
		if existingURL == remoteURL {
			ui.Success("Remote already configured correctly")
		} else {
			if err := system.RunGit("remote", "set-url", cfg.Remote, remoteURL); err != nil {
				ui.Error("Failed to update existing remote")
				return
			}
			ui.Success("Remote updated successfully")
		}
	} else {
		if err := system.RunGit("remote", "add", cfg.Remote, remoteURL); err != nil {
			ui.Error("Failed to add remote")
			return
		}
		ui.Success("Remote added successfully")
	}

	config.Save(cfg)

	ui.Success("Repository linked successfully")
	ui.Info("Authentication will be requested during push")
	ui.Info("Use GitHub username + Personal Access Token")
}
