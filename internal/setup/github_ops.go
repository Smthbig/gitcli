package setup

import (
	"fmt"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func setupGitHubToken() bool {
	ui.BoxHeader("GitHub Authentication")

	switch github.TokenSource() {
	case "environment":
		ui.Success("GitHub token available from environment variable")
		return true
	case "file":
		ui.Success("GitHub token already configured")
		return true
	}

	if !ui.Confirm("Configure GitHub token now?") {
		return true
	}

	token := strings.TrimSpace(ui.SecretInput("Paste GitHub token"))
	if token == "" {
		return false
	}

	if user, err := github.ValidateToken(token); err == nil {
		_ = github.SaveAuth(token, user)
		ui.Success("Authenticated as: " + user)
		return true
	}
	return false
}

func ensureGitHubRepo(cfg *config.Config) {
	if cfg.Owner == "" || cfg.Repo == "" {
		return
	}

	exists, _ := github.RepoExists(cfg.Owner, cfg.Repo)
	if exists {
		ui.Success("GitHub repository confirmed")
		return
	}

	if ui.Confirm("Create repository on GitHub if it does not exist?") {
		private := ui.Confirm("Make repository private?")
		if err := github.CreateRepo(cfg.Owner, cfg.Repo, private); err != nil {
			ui.Warn("Failed to create GitHub repository")
		} else {
			ui.Success("GitHub repository created")
		}
	}
}

func configureRemote(cfg *config.Config) error {
	if cfg.Remote == "" || cfg.Owner == "" || cfg.Repo == "" {
		return nil
	}

	url := fmt.Sprintf("https://github.com/%s/%s.git", cfg.Owner, cfg.Repo)
	return system.EnsureRemote(cfg.Remote, url)
}
