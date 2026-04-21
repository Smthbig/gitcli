package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Run performs the full setup / reconfigure flow.
*/
func Run() bool {
	ui.Clear()
	ui.Header("Git Genius Setup")

	cfg := config.Load()

	if !selectWorkDir(&cfg) {
		return false
	}
	config.Save(cfg)

	if err := system.EnsureGitInstalled(); err != nil {
		ui.Error("Git is required for setup")
		return false
	}

	if !system.EnsureGitRepo() {
		return false
	}

	system.EnsureSafeDirectory(cfg.GetWorkDir())

	if !system.EnsureBranchSync() {
		return false
	}

	if !ensureGitIdentity(cfg.GetWorkDir()) {
		return false
	}

	setupGitBasics(&cfg)

	if !ensureConfiguredBranch(&cfg) {
		return false
	}

	if !setupRepo(&cfg) {
		return false
	}

	if !setupGitHubToken() {
		return false
	}

	ensureGitHubRepo(&cfg)

	if err := configureRemote(&cfg); err != nil {
		ui.Error("Failed to configure git remote")
		ui.Info(err.Error())
		return false
	}

	offerFirstPush(&cfg)
	config.Save(cfg)

	ui.Header("Setup Summary")
	ui.Success("Project Dir : " + cfg.GetWorkDir())
	ui.Success("Branch      : " + cfg.Branch)
	ui.Success("Remote      : " + cfg.Remote)
	if cfg.Owner != "" && cfg.Repo != "" {
		ui.Success("Repository  : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	}
	ui.Success("Setup completed successfully")
	return true
}

///////////////////////////////////////////////////////////////
//////////////////// DIRECTORY ////////////////////////////////
///////////////////////////////////////////////////////////////

func selectWorkDir(cfg *config.Config) bool {
	cwd, _ := os.Getwd()
	if cfg.WorkDir == "" {
		cfg.WorkDir = cwd
	}

	ui.Info("Current directory: " + cwd)
	recent := config.RecentWorkDirs()
	if len(recent) > 0 {
		ui.Info("Recent project directories:")
		max := len(recent)
		if max > 3 {
			max = 3
		}
		for i := 0; i < max; i++ {
			ui.Info(fmt.Sprintf("  %d) %s", i+1, recent[i]))
		}
	}

	if !ui.Confirm("Use a different project directory?") {
		return true
	}

	for {
		dir := strings.TrimSpace(ui.Input("Enter full path of project directory (or 1/2/3 for recent)"))
		if dir == "" {
			ui.Error("Directory cannot be empty")
			continue
		}

		switch dir {
		case "1", "2", "3":
			idx := int(dir[0] - '1')
			if idx >= 0 && idx < len(recent) {
				dir = recent[idx]
			}
		}

		abs, err := filepath.Abs(dir)
		if err == nil {
			dir = abs
		}

		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			ui.Error("Invalid directory path")
			if !ui.ConfirmDefault("Try another directory?", true) {
				return false
			}
			continue
		}

		cfg.WorkDir = dir
		ui.Success("Project directory set to: " + dir)
		return true
	}
}

///////////////////////////////////////////////////////////////
//////////////////// GIT BASICS ////////////////////////////////
///////////////////////////////////////////////////////////////

func setupGitBasics(cfg *config.Config) {
	cfg.Branch = strings.TrimSpace(ui.InputDefault("Default branch", cfg.Branch))
	cfg.Remote = strings.TrimSpace(ui.InputDefault("Remote name", cfg.Remote))
}

func ensureConfiguredBranch(cfg *config.Config) bool {
	desired := strings.TrimSpace(cfg.Branch)
	if desired == "" {
		return true
	}

	current := system.CurrentGitBranch()
	if current == desired {
		return true
	}

	if current == "" {
		if err := system.PrepareBranch(desired); err != nil {
			ui.Error("Failed to prepare configured branch")
			ui.Info(err.Error())
			return false
		}
		ui.Success("Prepared branch: " + desired)
		return true
	}

	ui.Warn("Current git branch differs from configured branch")
	ui.Info("Current branch    : " + current)
	ui.Info("Configured branch : " + desired)

	if system.HasLocalBranch(desired) {
		if ui.ConfirmDefault("Switch to the configured branch now?", true) {
			if err := system.SwitchToBranch(desired); err != nil {
				ui.Error("Failed to switch branch")
				ui.Info(err.Error())
				return false
			}
			ui.Success("Switched to branch: " + desired)
			return true
		}
	} else if ui.ConfirmDefault("Create and switch to the configured branch now?", true) {
		if err := system.CreateBranch(desired); err != nil {
			ui.Error("Failed to create branch")
			ui.Info(err.Error())
			return false
		}
		ui.Success("Created and switched to branch: " + desired)
		return true
	}

	cfg.Branch = current
	ui.Info("Keeping current branch and syncing config")
	return true
}

///////////////////////////////////////////////////////////////
//////////////////// REPOSITORY INFO ///////////////////////////
///////////////////////////////////////////////////////////////

func setupRepo(cfg *config.Config) bool {
	ui.Header("GitHub Repository")

	defaultRepo := filepath.Base(cfg.GetWorkDir())
	if cfg.Repo != "" {
		defaultRepo = cfg.Repo
	}

	for {
		cfg.Owner = strings.TrimSpace(ui.InputDefault("GitHub username or organisation", cfg.Owner))
		cfg.Repo = strings.TrimSpace(ui.InputDefault("Repository name", defaultRepo))

		if cfg.Owner == "" || cfg.Repo == "" {
			ui.Error("Owner and repository name are required")
			continue
		}

		ui.Info("Target repository:")
		ui.Info("https://github.com/" + cfg.Owner + "/" + cfg.Repo)
		return true
	}
}

///////////////////////////////////////////////////////////////
//////////////////// TOKEN ////////////////////////////////////
///////////////////////////////////////////////////////////////

func setupGitHubToken() bool {
	ui.Header("GitHub Authentication")

	switch github.TokenSource() {
	case "environment":
		ui.Success("GitHub token available from environment variable: " + github.EnvTokenName)
		if !github.HasStoredToken() && ui.Confirm("Save environment token to local Git Genius storage?") {
			if err := github.Save(github.GetToken()); err != nil {
				ui.Warn("Could not save environment token locally")
				ui.Info(err.Error())
			} else {
				ui.Success("Environment token saved locally")
			}
		}
		return true
	case "file":
		ui.Success("GitHub token already configured")
		return true
	}

	ui.Info("Create token at: https://github.com/settings/tokens")
	ui.Info("Required scope: repo")
	ui.Info("Automation option: export " + github.EnvTokenName + "=your_token")

	if !ui.Confirm("Configure GitHub token now?") {
		return true
	}

	for {
		token := strings.TrimSpace(ui.SecretInput("Paste GitHub token"))
		if token == "" {
			ui.Error("Token cannot be empty")
			if !ui.ConfirmDefault("Try again?", true) {
				return false
			}
			continue
		}

		if user, err := github.ValidateToken(token); err != nil {
			ui.Warn("Could not validate GitHub token")
			ui.Info(err.Error())
			if !ui.Confirm("Save token anyway?") {
				if !ui.ConfirmDefault("Try a different token?", true) {
					return false
				}
				continue
			}
		} else {
			ui.Success("GitHub authenticated as: " + user)
		}

		if err := github.Save(token); err != nil {
			ui.Error("Failed to save token")
			ui.Info(err.Error())
			return false
		}

		ui.Success("Token saved")
		return true
	}
}

///////////////////////////////////////////////////////////////
//////////////////// REPO CREATION /////////////////////////////
///////////////////////////////////////////////////////////////

func ensureGitHubRepo(cfg *config.Config) {
	if github.GetToken() == "" {
		ui.Warn("No GitHub token available, skipping repository verification")
		return
	}

	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
	if err == nil && exists {
		ui.Success("GitHub repository already exists")
		return
	}

	if err != nil {
		ui.Warn("Could not verify repository with GitHub API")
		ui.Info(err.Error())
	}

	if !ui.Confirm("Create repository on GitHub if it does not exist?") {
		return
	}

	private := ui.Confirm("Make repository private?")
	if err := github.CreateRepo(cfg.Owner, cfg.Repo, private); err != nil {
		ui.Warn("Repository creation failed")
		ui.Info(err.Error())
		return
	}

	ui.Success("Repository created successfully")
}

///////////////////////////////////////////////////////////////
//////////////////// REMOTE ///////////////////////////////////
///////////////////////////////////////////////////////////////

func configureRemote(cfg *config.Config) error {
	if cfg.Remote == "" {
		cfg.Remote = "origin"
	}

	url := fmt.Sprintf("https://github.com/%s/%s.git", cfg.Owner, cfg.Repo)
	return system.EnsureRemote(cfg.Remote, url)
}

///////////////////////////////////////////////////////////////
//////////////////// FIRST PUSH ////////////////////////////////
///////////////////////////////////////////////////////////////

func offerFirstPush(cfg *config.Config) {
	if !ui.Confirm("Push current code to GitHub now?") {
		return
	}

	if cfg.Branch != "" && system.CurrentGitBranch() == "" {
		if err := system.PrepareBranch(cfg.Branch); err != nil {
			ui.Error("Failed to prepare branch for first push")
			ui.Info(err.Error())
			return
		}
	}

	msg := strings.TrimSpace(ui.InputDefault("Commit message", "Initial commit"))

	if err := system.RunGit("add", "."); err != nil {
		ui.Error("git add failed")
		return
	}

	stagedNames, err := system.GitOutput("diff", "--cached", "--name-only")
	if err != nil {
		ui.Error("Failed to check staged changes")
		return
	}
	if stagedNames == "" {
		ui.Warn("Nothing to commit yet (working tree is empty)")
		ui.Info("Add at least one file, then push again")
		return
	}

	if err := system.RunGit("commit", "-m", msg); err != nil {
		ui.Error("Initial commit failed")
		return
	}

	branch := system.CurrentGitBranch()
	if branch == "" {
		branch = cfg.Branch
	}
	if branch == "" {
		ui.Error("No branch detected")
		return
	}

	if err := system.RunGit("push", "-u", cfg.Remote, branch); err != nil {
		ui.Error("Push failed")
		ui.Info("Run Daily Git Operations -> Push after reviewing the remote and branch setup")
		return
	}

	cfg.FirstPushDone = true
	config.Save(*cfg)
	ui.Success("Code pushed successfully")
}

///////////////////////////////////////////////////////////////
//////////////////// GIT IDENTITY //////////////////////////////
///////////////////////////////////////////////////////////////

func ensureGitIdentity(workDir string) bool {
	name, _ := system.GitOutputAt(workDir, "config", "--get", "user.name")
	email, _ := system.GitOutputAt(workDir, "config", "--get", "user.email")

	if name != "" && email != "" {
		ui.Success("Git identity already configured")
		return true
	}

	ui.Warn("Git identity not configured")

	if name == "" {
		val := strings.TrimSpace(ui.Input("Enter your name"))
		if val == "" {
			ui.Error("Name cannot be empty")
			return false
		}
		if err := system.RunGitAt(workDir, "config", "user.name", val); err != nil {
			ui.Error("Failed to set user.name")
			return false
		}
	}

	if email == "" {
		val := strings.TrimSpace(ui.Input("Enter your email"))
		if val == "" {
			ui.Error("Email cannot be empty")
			return false
		}
		if err := system.RunGitAt(workDir, "config", "user.email", val); err != nil {
			ui.Error("Failed to set user.email")
			return false
		}
	}

	ui.Success("Git identity configured")
	return true
}
