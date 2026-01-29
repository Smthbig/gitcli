package setup

import (
	"fmt"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
CreateOrLinkRepo creates OR links a GitHub repository
Android-safe, Org-safe, Token-safe
Used in:
Tools → Create / Link GitHub Repository
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

	// --------------------------------------------------
	// Token check (NOT validating ownership here)
	// --------------------------------------------------
	token := github.GetToken()
	if token == "" {
		ui.Error("GitHub token not configured")
		ui.Info("Run: Tools → Setup / Reconfigure")
		return
	}

	// --------------------------------------------------
	// Internet check (reliable now)
	// --------------------------------------------------
	if !system.Online {
		ui.Warn("No internet connection detected")
		ui.Info("Cannot verify or create GitHub repository")
		return
	}

	// --------------------------------------------------
	// Repo existence check
	// --------------------------------------------------
	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
	if err != nil {
		ui.Error("Failed to communicate with GitHub API")
		ui.Info("Possible reasons:")
		ui.Info("• Invalid token")
		ui.Info("• GitHub rate limit")
		ui.Info("• Organisation access denied")
		system.LogError("repo exists check failed", err)
		return
	}

	// --------------------------------------------------
	// Create repo if missing
	// --------------------------------------------------
	if !exists {
		ui.Warn("GitHub repository does not exist")

		if !ui.Confirm("Create repository on GitHub now?") {
			ui.Warn("Repository creation skipped")
			return
		}

		private := ui.Confirm("Make repository PRIVATE?")

		err := github.CreateRepo(cfg.Owner, cfg.Repo, private)
		if err != nil {
			ui.Error("Repository creation failed")

			ui.Info("Common causes:")
			ui.Info("• You are NOT an owner/admin of organisation: " + cfg.Owner)
			ui.Info("• Token lacks required permissions")
			ui.Info("• Organisation restricts repo creation")

			ui.Info("Fix:")
			ui.Info("• Ask org admin to create repo manually")
			ui.Info("• OR push to your personal account")

			system.LogError("repo creation failed", err)
			return
		}

		ui.Success("GitHub repository created successfully")
	} else {
		ui.Success("GitHub repository already exists")
	}

	// --------------------------------------------------
	// Configure remote (SECURE, NO TOKEN IN URL)
	// --------------------------------------------------
	remoteURL := fmt.Sprintf(
		"https://github.com/%s/%s.git",
		cfg.Owner,
		cfg.Repo,
	)

	_ = system.RunGit("remote", "remove", cfg.Remote)

	if err := system.RunGit("remote", "add", cfg.Remote, remoteURL); err != nil {
		ui.Error("Failed to configure git remote")
		system.LogError("remote add failed", err)
		return
	}

	config.Save(cfg)

	ui.Success("GitHub repository linked successfully")
	ui.Info("You will be asked for authentication on first push")
	ui.Info("Use your GitHub username + token as password")
}
