package setup

import (
	"fmt"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
CreateOrLinkRepo creates or links a GitHub repository without requiring a
full setup rerun.
*/
func CreateOrLinkRepo() bool {
	ui.Clear()
	ui.Header("Create / Link GitHub Repository")

	cfg := config.Load()

	if !system.EnsureGitRepo() {
		return false
	}

	system.EnsureSafeDirectory(cfg.GetWorkDir())

	defaultRepo := cfg.Repo
	if defaultRepo == "" {
		defaultRepo = filepath.Base(cfg.GetWorkDir())
	}

	cfg.Owner = strings.TrimSpace(ui.InputDefault("GitHub username or organisation", cfg.Owner))
	cfg.Repo = strings.TrimSpace(ui.InputDefault("Repository name", defaultRepo))
	cfg.Remote = strings.TrimSpace(ui.InputDefault("Remote name", cfg.Remote))

	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("Owner and repository name are required")
		return false
	}
	if cfg.Remote == "" {
		cfg.Remote = "origin"
	}

	if github.GetToken() == "" {
		ui.Warn("No GitHub token configured")
		ui.Info("Remote will still be linked, but automatic repo verification and creation are skipped")
	} else {
		exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
		switch {
		case err != nil:
			ui.Warn("Could not verify repository with GitHub API")
			ui.Info(err.Error())
		case exists:
			ui.Success("GitHub repository exists")
		default:
			ui.Warn("Repository does not exist on GitHub")
			if ui.Confirm("Create repository on GitHub now?") {
				private := ui.Confirm("Make repository private?")
				if err := github.CreateRepo(cfg.Owner, cfg.Repo, private); err != nil {
					ui.Error("Repository creation failed")
					ui.Info(err.Error())
					return false
				}
				ui.Success("Repository created successfully")
			}
		}
	}

	remoteURL := fmt.Sprintf("https://github.com/%s/%s.git", cfg.Owner, cfg.Repo)
	if err := system.EnsureRemote(cfg.Remote, remoteURL); err != nil {
		ui.Error("Failed to configure remote")
		ui.Info(err.Error())
		return false
	}

	config.Save(cfg)

	ui.Success("Repository linked successfully")
	ui.Info("Use Daily Git Operations -> Push to publish local changes")
	return true
}
