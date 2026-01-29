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
Run executes the full guided setup
SINGLE SOURCE OF TRUTH for configuration
Android-safe & restricted-kernel safe
*/
func Run() {
	ui.Clear()
	ui.Header("Git Genius Setup")

	cfg := config.Load()

	// STEP 0: Project directory
	if !selectWorkDir(&cfg) {
		return
	}

	// 🔒 Persist early so system layer uses correct dir
	config.Save(cfg)

	// STEP 1: Ensure git repo
	if !system.EnsureGitRepo() {
		return
	}

	// 🔐 Android + Git ≥ 2.35 safety
	system.EnsureSafeDirectory(cfg.WorkDir)

	// STEP 2: Sync branch safely
	system.EnsureBranchSync()

	// STEP 2.5: Git identity (CRITICAL)
	if !ensureGitIdentity(cfg.WorkDir) {
		return
	}

	// STEP 3: Git basics
	setupGitBasics(&cfg)

	// STEP 4: GitHub repo info
	if !setupRepo(&cfg) {
		return
	}

	// STEP 5: GitHub token
	if !setupGitHubToken() {
		return
	}

	// STEP 6: Ensure GitHub repo exists
	if !ensureGitHubRepo(&cfg) {
		return
	}

	// STEP 7: Configure remote
	if err := configureRemote(&cfg); err != nil {
		ui.Error("Failed to configure git remote")
		return
	}

	// STEP 8: Optional first push
	offerFirstPush(&cfg)

	config.Save(cfg)

	ui.Header("Setup Summary")
	ui.Success("Project Dir : " + cfg.GetWorkDir())
	ui.Success("Branch      : " + cfg.Branch)
	ui.Success("Remote      : " + cfg.Remote)
	ui.Success("Repository  : https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	ui.Success("Setup completed successfully 🎉")
}

/* ============================================================
   STEP 0: Project directory
   ============================================================ */

func selectWorkDir(cfg *config.Config) bool {
	cwd, _ := os.Getwd()
	ui.Info("Current directory: " + cwd)

	if cfg.WorkDir == "" {
		cfg.WorkDir = cwd
	}

	if !ui.Confirm("Do you want to use a DIFFERENT project directory?") {
		return true
	}

	dir := ui.Input("Enter full path of project directory")
	if dir == "" {
		ui.Error("Directory path cannot be empty")
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

/* ============================================================
   STEP 3: Git basics
   ============================================================ */

func setupGitBasics(cfg *config.Config) {
	if b := ui.Input("Default branch [" + cfg.Branch + "]"); b != "" {
		cfg.Branch = b
	}
	if r := ui.Input("Remote name [" + cfg.Remote + "]"); r != "" {
		cfg.Remote = r
	}
}

/* ============================================================
   STEP 4: GitHub repo info
   ============================================================ */

func setupRepo(cfg *config.Config) bool {
	ui.Header("GitHub Repository")

	if cfg.Repo == "" && cfg.WorkDir != "" {
		cfg.Repo = filepath.Base(cfg.WorkDir)
	}

	if cfg.Owner == "" {
		cfg.Owner = ui.Input("GitHub username or organisation")
	}
	if cfg.Repo == "" {
		cfg.Repo = ui.Input("Repository name")
	}

	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("Owner and repository name are required")
		return false
	}

	ui.Info("Target repository:")
	ui.Info("https://github.com/" + cfg.Owner + "/" + cfg.Repo)
	return true
}

/* ============================================================
   STEP 5: GitHub token
   ============================================================ */

func setupGitHubToken() bool {
	ui.Header("GitHub Authentication")

	if github.GetToken() != "" {
		ui.Success("GitHub token already configured")
		return true
	}

	ui.Info("Create token at: https://github.com/settings/tokens")
	ui.Info("Required scope: repo")

	if !ui.Confirm("Do you want to configure GitHub token now?") {
		ui.Warn("Skipping token setup")
		return true
	}

	token := ui.SecretInput("Paste GitHub token")
	if token == "" {
		ui.Error("Empty token")
		return false
	}

	if err := github.Save(token); err != nil {
		ui.Error("Failed to save token")
		return false
	}

	user, err := github.Validate()
	if err != nil {
		ui.Error("Invalid GitHub token")
		github.Delete()
		return false
	}

	ui.Success("Authenticated as: " + user)
	return true
}

/* ============================================================
   STEP 6: Ensure GitHub repo
   ============================================================ */

func ensureGitHubRepo(cfg *config.Config) bool {
	if !system.Online {
		ui.Warn("Offline mode detected")
		ui.Info("Cannot verify or create GitHub repository while offline")
		return true // allow setup to continue
	}

	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
	if err != nil {
		ui.Error("GitHub API error")
		ui.Info("Check token permissions or organisation rights")
		return false
	}

	if exists {
		ui.Success("GitHub repository exists")
		return true
	}

	ui.Warn("GitHub repository does not exist")

	if !ui.Confirm("Create this repository on GitHub?") {
		return true
	}

	private := ui.Confirm("Make repository PRIVATE?")

	if err := github.CreateRepo(cfg.Owner, cfg.Repo, private); err != nil {
		ui.Error("Repository creation failed")
		ui.Info("Possible reasons:")
		ui.Info("• Owner is an organisation")
		ui.Info("• Token lacks repo permission")
		ui.Info("• Org restricts repo creation")
		return false
	}

	ui.Success("GitHub repository created successfully")
	return true
}

/* ============================================================
   STEP 7: Configure remote
   ============================================================ */

func configureRemote(cfg *config.Config) error {
	token := github.GetToken()

	url := fmt.Sprintf(
		"https://%s@github.com/%s/%s.git",
		token,
		cfg.Owner,
		cfg.Repo,
	)

	_ = system.RunGit("remote", "remove", cfg.Remote)
	return system.RunGit("remote", "add", cfg.Remote, url)
}

/* ============================================================
   STEP 8: First push
   ============================================================ */

func offerFirstPush(cfg *config.Config) {
	if !ui.Confirm("Push current code to GitHub now?") {
		return
	}

	msg := ui.Input("Initial commit message")
	if msg == "" {
		msg = "Initial commit"
	}

	_ = system.RunGit("add", ".")
	_ = system.RunGit("commit", "-m", msg)
	_ = system.RunGit("push", "-u", cfg.Remote, cfg.Branch)

	ui.Success("Code pushed successfully")
}

/* ============================================================
   STEP 2.5: Git identity (ANDROID SAFE, LOCAL ONLY)
   ============================================================ */

func ensureGitIdentity(workDir string) bool {
	name := gitConfig(workDir, "user.name")
	email := gitConfig(workDir, "user.email")

	if name != "" && email != "" {
		ui.Success("Git identity already configured")
		return true
	}

	ui.Warn("Git identity not configured")
	ui.Info("Commits require user.name and user.email")

	if !ui.Confirm("Configure git identity now?") {
		return false
	}

	if name == "" {
		val := ui.Input("Enter your name")
		if val == "" {
			ui.Error("Name cannot be empty")
			return false
		}
		if err := system.GitCmdAt(workDir, "config", "user.name", val).Run(); err != nil {
			ui.Error("Failed to set git user.name")
			return false
		}
	}

	if email == "" {
		val := ui.Input("Enter your email")
		if val == "" {
			ui.Error("Email cannot be empty")
			return false
		}
		if err := system.GitCmdAt(workDir, "config", "user.email", val).Run(); err != nil {
			ui.Error("Failed to set git user.email")
			return false
		}
	}

	ui.Success("Git identity configured (local repository)")
	return true
}

/* ============================================================
   Helpers
   ============================================================ */

func gitConfig(dir, key string) string {
	var out strings.Builder
	cmd := system.GitCmdAt(dir, "config", "--get", key)
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return ""
	}
	return strings.TrimSpace(out.String())
}
