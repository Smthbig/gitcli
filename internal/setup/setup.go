package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Stable Production Setup Flow
- No offline mode
- No RepoExists
- No token in remote URL
- Proper error handling
*/

func Run() {
	ui.Clear()
	ui.Header("Git Genius Setup")

	cfg := config.Load()

	// STEP 0 — Select project directory
	if !selectWorkDir(&cfg) {
		return
	}
	config.Save(cfg)

	// STEP 1 — Ensure git installed
	if err := system.EnsureGitInstalled(); err != nil {
		ui.Error("Git is required.")
		return
	}

	// STEP 2 — Ensure git repo
	if !system.EnsureGitRepo() {
		return
	}

	// STEP 3 — Safe directory
	system.EnsureSafeDirectory(cfg.WorkDir)

	// STEP 4 — Branch sync
	system.EnsureBranchSync()

	// STEP 5 — Git identity
	if !ensureGitIdentity(cfg.WorkDir) {
		return
	}

	// STEP 6 — Git basics
	setupGitBasics(&cfg)

	// STEP 7 — Repo info
	if !setupRepo(&cfg) {
		return
	}

	// STEP 8 — Token
	if !setupGitHubToken() {
		return
	}

	// STEP 9 — Optional repo creation
	ensureGitHubRepo(&cfg)

	// STEP 10 — Configure remote
	if err := configureRemote(&cfg); err != nil {
		ui.Error("Failed to configure git remote")
		return
	}

	// STEP 11 — Optional first push
	offerFirstPush(&cfg)

	config.Save(cfg)

	ui.Header("Setup Summary")
	ui.Success("Project Dir : " + cfg.GetWorkDir())
	ui.Success("Branch      : " + cfg.Branch)
	ui.Success("Remote      : " + cfg.Remote)
	ui.Success("Repository  : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	ui.Success("Setup completed successfully 🎉")
}

///////////////////////////////////////////////////////////////
//////////////////// DIRECTORY ////////////////////////////////
///////////////////////////////////////////////////////////////

func selectWorkDir(cfg *config.Config) bool {
	cwd, _ := os.Getwd()
	ui.Info("Current directory: " + cwd)

	if cfg.WorkDir == "" {
		cfg.WorkDir = cwd
	}

	if !ui.Confirm("Use a DIFFERENT project directory?") {
		return true
	}

	dir := ui.Input("Enter full path of project directory")
	if dir == "" {
		ui.Error("Directory cannot be empty")
		return false
	}

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		ui.Error("Invalid directory path")
		return false
	}

	cfg.WorkDir = dir
	ui.Success("Project directory set to: " + dir)
	return true
}

///////////////////////////////////////////////////////////////
//////////////////// GIT BASICS ////////////////////////////////
///////////////////////////////////////////////////////////////

func setupGitBasics(cfg *config.Config) {
	if b := ui.Input("Default branch [" + cfg.Branch + "]"); b != "" {
		cfg.Branch = b
	}
	if r := ui.Input("Remote name [" + cfg.Remote + "]"); r != "" {
		cfg.Remote = r
	}
}

///////////////////////////////////////////////////////////////
//////////////////// REPOSITORY INFO ///////////////////////////
///////////////////////////////////////////////////////////////

func setupRepo(cfg *config.Config) bool {
	ui.Header("GitHub Repository")

	defaultRepo := filepath.Base(cfg.WorkDir)

	if cfg.Owner == "" {
		cfg.Owner = ui.Input("GitHub username or organisation")
	}

	repoInput := ui.Input("Repository name [" + defaultRepo + "]")
	if repoInput == "" {
		cfg.Repo = defaultRepo
	} else {
		cfg.Repo = repoInput
	}

	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("Owner and repository name are required")
		return false
	}

	ui.Info("Target repository:")
	ui.Info("https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	return true
}

///////////////////////////////////////////////////////////////
//////////////////// TOKEN ////////////////////////////////////
///////////////////////////////////////////////////////////////

func setupGitHubToken() bool {
	ui.Header("GitHub Authentication")

	if github.GetToken() != "" {
		ui.Success("GitHub token already configured")
		return true
	}

	ui.Info("Create token at: https://github.com/settings/tokens")
	ui.Info("Required scope: repo")

	if !ui.Confirm("Configure GitHub token now?") {
		return true
	}

	token := ui.SecretInput("Paste GitHub token")
	if token == "" {
		ui.Error("Token cannot be empty")
		return false
	}

	if err := github.Save(token); err != nil {
		ui.Error("Failed to save token")
		return false
	}

	ui.Success("Token saved")
	return true
}

///////////////////////////////////////////////////////////////
//////////////////// REPO CREATION /////////////////////////////
///////////////////////////////////////////////////////////////

func ensureGitHubRepo(cfg *config.Config) {
	if !ui.Confirm("Create repository on GitHub if not exists?") {
		return
	}

	private := ui.Confirm("Make repository PRIVATE?")

	err := github.CreateRepo(cfg.Owner, cfg.Repo, private)
	if err != nil {
		ui.Warn("Repository creation failed:")
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

	url := fmt.Sprintf(
		"https://github.com/%s/%s.git",
		cfg.Owner,
		cfg.Repo,
	)

	// Remove existing remote (ignore error)
	_ = system.RunGit("remote", "remove", cfg.Remote)

	if err := system.RunGit("remote", "add", cfg.Remote, url); err != nil {
		return err
	}

	return nil
}

///////////////////////////////////////////////////////////////
//////////////////// FIRST PUSH ////////////////////////////////
///////////////////////////////////////////////////////////////

func offerFirstPush(cfg *config.Config) {

	if !ui.Confirm("Push current code to GitHub now?") {
		return
	}

	msg := ui.Input("Commit message")
	if msg == "" {
		msg = "Initial commit"
	}

	if err := system.RunGit("add", "."); err != nil {
		ui.Error("git add failed")
		return
	}

	// Commit may fail if nothing to commit
	_ = system.RunGit("commit", "-m", msg)

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
		return
	}

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
		val := ui.Input("Enter your name")
		if val == "" {
			return false
		}
		if err := system.RunGitAt(workDir, "config", "user.name", val); err != nil {
			ui.Error("Failed to set user.name")
			return false
		}
	}

	if email == "" {
		val := ui.Input("Enter your email")
		if val == "" {
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
