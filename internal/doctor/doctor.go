package doctor

import (
	"os"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Doctor:
Full system + git + github health check
No offline mode
Compatible with new system layer
Android safe
*/
func Run() bool {

	ui.BoxHeader("Git Genius Doctor 🩺")

	checkGitInstalled()
	checkWorkDir()
	checkGitRepo()
	checkGitBranch()
	checkGitIdentity()
	checkRemote()
	checkGitHubToken()
	checkGitHubRepo()
	checkGitCredentialHelper()
	checkErrorLog()

	ui.Success("Doctor check completed")
	return true
}

///////////////////////////////////////////////////////////////
// CHECKS
///////////////////////////////////////////////////////////////

func checkGitInstalled() {
	if system.CommandExists("git") {
		ui.Success("Git installed")
		return
	}

	ui.Error("Git not found in PATH")
}

func checkWorkDir() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		ui.Error("Invalid project directory: " + dir)
		return
	}

	ui.Success("Project directory: " + dir)
}

func checkGitRepo() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	if system.IsGitRepoAt(dir) {
		ui.Success("Git repository detected")
		system.EnsureSafeDirectory(dir)
		return
	}

	ui.Warn("No git repository found")
	ui.Info("Run Setup to initialize repository")
}

func checkGitBranch() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	branch, err := system.GitOutputAt(dir, "branch", "--show-current")
	if err != nil || branch == "" {
		ui.Warn("No commits yet (branch not created)")
		return
	}

	ui.Success("Current git branch: " + branch)

	if branch != cfg.Branch {
		ui.Warn("Branch mismatch detected")
		ui.Info("Config branch : " + cfg.Branch)
		ui.Info("Git branch    : " + branch)
	}
}

func checkGitIdentity() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	name, _ := system.GitOutputAt(dir, "config", "--get", "user.name")
	email, _ := system.GitOutputAt(dir, "config", "--get", "user.email")

	if name != "" && email != "" {
		ui.Success("Git identity configured")
		ui.Info("Name : " + name)
		ui.Info("Email: " + email)
		return
	}

	ui.Warn("Git identity not configured")
	ui.Info("Run Setup to configure git identity")
}

func checkRemote() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	if cfg.Remote == "" {
		ui.Warn("No git remote configured")
		return
	}

	_, err := system.GitOutputAt(dir, "remote", "get-url", cfg.Remote)
	if err != nil {
		ui.Warn("Remote not found: " + cfg.Remote)
		ui.Info("Run Create / Link GitHub Repository")
		return
	}

	ui.Success("Git remote configured: " + cfg.Remote)
}

func checkGitHubToken() {

	token := github.GetToken()
	if token == "" {
		ui.Warn("GitHub token not configured")
		return
	}

	if github.TokenSource() == "environment" {
		ui.Info("GitHub token source: environment variable " + github.EnvTokenName)
	}

	client, err := github.NewClient()
	if err != nil {
		ui.Error("Invalid GitHub token")
		return
	}

	user, err := client.GetAuthenticatedUser()
	if err != nil {
		ui.Warn("GitHub token invalid or expired")
		return
	}

	ui.Success("GitHub authenticated as: " + user)
}

func checkGitHubRepo() {
	cfg := config.Load()

	if cfg.Owner == "" || cfg.Repo == "" {
		return
	}

	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
	if err != nil {
		ui.Warn("Unable to check GitHub repository")
		return
	}

	if exists {
		ui.Success("GitHub repository exists")
	} else {
		ui.Warn("GitHub repository does not exist")
	}
}

func checkErrorLog() {
	logPath := system.ErrorLogPath()

	if _, err := os.Stat(logPath); err == nil {
		ui.Warn("Error log exists")
		ui.Info("Check: " + logPath)
	} else {
		ui.Success("No error log found")
	}
}

func checkGitCredentialHelper() {
	helper, err := system.GitCredentialHelper()
	if err != nil || helper == "" {
		ui.Warn("No global Git credential helper configured")
		ui.Info("Run Tools -> Git Auth / Credential Helper to reduce repeated HTTPS prompts")
		return
	}

	ui.Success("Git credential helper: " + helper)

	if helper == "store" {
		ui.Warn("Credential helper stores credentials in plain text at ~/.git-credentials")
	}

	if helper == "!/usr/bin/gh auth git-credential" && !system.CommandExists("gh") {
		ui.Warn("Git credential helper uses gh, but gh is not installed")
		ui.Info("Push may still work, but auth helper warnings will appear")
	}
}
