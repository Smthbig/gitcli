package doctor

import (
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// Run performs full system + git health check (ANDROID SAFE)
func Run() {
	ui.Header("Git Genius Doctor 🩺")

	system.CheckInternet()

	checkGitInstalled()
	checkWorkDir()
	checkGitRepo()
	checkGitBranch()
	checkGitIdentity()
	checkRemote()
	checkInternet()
	checkGitHubToken()
	checkGitHubRepo()
	checkErrorLog()

	ui.Success("Doctor check completed")
}

/* ============================================================
   CHECKS
   ============================================================ */

func checkGitInstalled() {
	if system.CommandExists("git") {
		ui.Success("Git installed")
		return
	}

	ui.Error("Git not found in PATH")
	ui.Info("Please install git manually for your environment")
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

	if ui.Confirm("Initialize git repository here?") {
		if err := system.RunGitAt(dir, "init"); err != nil {
			ui.Error("Failed to initialize git repository")
			return
		}
		system.EnsureSafeDirectory(dir)
		ui.Success("Git repository initialized")
	}
}

func checkGitBranch() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	current := system.CurrentGitBranchAt(dir)

	if current == "" {
		ui.Warn("No commits yet (branch not created)")
		return
	}

	ui.Success("Current git branch: " + current)

	if current != cfg.Branch {
		ui.Warn("Branch mismatch detected")
		ui.Info("Config branch : " + cfg.Branch)
		ui.Info("Git branch    : " + current)
		ui.Info("Run Setup to safely sync branch")
	}
}

/*
Git identity check (ANDROID SAFE, LOCAL REPO ONLY)
*/
func checkGitIdentity() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	name := gitConfig(dir, "user.name")
	email := gitConfig(dir, "user.email")

	if name != "" && email != "" {
		ui.Success("Git identity configured")
		ui.Info("Name : " + name)
		ui.Info("Email: " + email)
		return
	}

	ui.Warn("Git identity not configured")
	ui.Info("Commits may appear as root@localhost")

	if !ui.Confirm("Configure git identity for THIS repo now?") {
		return
	}

	if name == "" {
		val := ui.Input("Enter your name")
		if val != "" {
			_ = system.RunGitAt(dir, "config", "user.name", val)
		}
	}

	if email == "" {
		val := ui.Input("Enter your email")
		if val != "" {
			_ = system.RunGitAt(dir, "config", "user.email", val)
		}
	}

	ui.Success("Git identity configured (local repository)")
}

func checkRemote() {
	cfg := config.Load()
	dir := cfg.GetWorkDir()

	if cfg.Remote == "" {
		ui.Warn("No git remote configured")
		return
	}

	if err := system.RunGitAt(dir, "remote", "get-url", cfg.Remote); err != nil {
		ui.Warn("Remote not found: " + cfg.Remote)
		ui.Info("Run Tools → Create / Link GitHub Repository")
		return
	}

	ui.Success("Git remote configured: " + cfg.Remote)
}

func checkInternet() {
	if system.Online {
		ui.Success("Internet connection available")
	} else {
		ui.Warn("Offline mode detected")
		ui.Info("GitHub validation & push may fail")
	}
}

func checkGitHubToken() {
	token := github.GetToken()
	if token == "" {
		ui.Warn("GitHub token not configured")
		ui.Info("Run Setup to configure token")
		return
	}

	user, err := github.Validate()
	if err != nil {
		ui.Error("GitHub token invalid or expired")
		ui.Info("Run Setup to reconfigure token")
		return
	}

	if user == "offline-mode" {
		ui.Warn("GitHub token validation skipped (offline)")
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
		ui.Info("Run Tools → Create / Link GitHub Repository")
	}
}

func checkErrorLog() {
	cfg := config.Load()
	logPath := filepath.Join(cfg.GetWorkDir(), ".git", ".genius", "error.log")

	if _, err := os.Stat(logPath); err == nil {
		ui.Warn("Error log exists")
		ui.Info("Check: " + logPath)
	} else {
		ui.Success("No error log found")
	}
}

/* ============================================================
   HELPERS
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
