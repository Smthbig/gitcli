package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

const debugLogPath = "/home/mohan/coding/gitcli/.cursor/debug-bc08dd.log"

func debugLog(runID, hypothesisID, location, message string, data map[string]interface{}) {
	entry := map[string]interface{}{
		"sessionId":    "bc08dd",
		"runId":        runID,
		"hypothesisId": hypothesisID,
		"location":     location,
		"message":      message,
		"data":         data,
		"timestamp":    time.Now().UnixMilli(),
	}

	b, err := json.Marshal(entry)
	if err != nil {
		return
	}

	_ = os.MkdirAll(filepath.Dir(debugLogPath), 0700)
	f, err := os.OpenFile(debugLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(b, '\n'))
}

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
	// #region agent log
	debugLog("pre-fix", "H1", "internal/setup/setup.go:Run", "setup started", map[string]interface{}{
		"initialWorkDir": cfg.WorkDir,
		"initialBranch":  cfg.Branch,
		"initialRemote":  cfg.Remote,
	})
	// #endregion

	// STEP 0 — Select project directory
	if !selectWorkDir(&cfg) {
		// #region agent log
		debugLog("pre-fix", "H1", "internal/setup/setup.go:Run", "selectWorkDir returned false", map[string]interface{}{})
		// #endregion
		return
	}
	config.Save(cfg)
	// #region agent log
	debugLog("pre-fix", "H1", "internal/setup/setup.go:Run", "workdir saved", map[string]interface{}{
		"savedWorkDir": cfg.WorkDir,
	})
	// #endregion

	// STEP 1 — Ensure git installed
	if err := system.EnsureGitInstalled(); err != nil {
		ui.Error("Git is required.")
		return
	}

	// STEP 2 — Ensure git repo
	if !system.EnsureGitRepo() {
		// #region agent log
		debugLog("pre-fix", "H2", "internal/setup/setup.go:Run", "EnsureGitRepo returned false", map[string]interface{}{
			"workDir": cfg.WorkDir,
		})
		// #endregion
		return
	}
	// #region agent log
	debugLog("pre-fix", "H2", "internal/setup/setup.go:Run", "EnsureGitRepo succeeded", map[string]interface{}{
		"workDir": cfg.WorkDir,
	})
	// #endregion

	// STEP 3 — Safe directory
	system.EnsureSafeDirectory(cfg.WorkDir)

	// STEP 4 — Branch sync
	system.EnsureBranchSync()

	// STEP 5 — Git identity
	if !ensureGitIdentity(cfg.WorkDir) {
		// #region agent log
		debugLog("pre-fix", "H3", "internal/setup/setup.go:Run", "ensureGitIdentity returned false", map[string]interface{}{
			"workDir": cfg.WorkDir,
		})
		// #endregion
		return
	}

	// STEP 6 — Git basics
	setupGitBasics(&cfg)

	// STEP 7 — Repo info
	if !setupRepo(&cfg) {
		// #region agent log
		debugLog("pre-fix", "H4", "internal/setup/setup.go:Run", "setupRepo returned false", map[string]interface{}{
			"owner": cfg.Owner,
			"repo":  cfg.Repo,
		})
		// #endregion
		return
	}
	// #region agent log
	debugLog("pre-fix", "H4", "internal/setup/setup.go:Run", "setupRepo succeeded", map[string]interface{}{
		"owner": cfg.Owner,
		"repo":  cfg.Repo,
	})
	// #endregion

	// STEP 8 — Token
	if !setupGitHubToken() {
		return
	}

	// STEP 9 — Optional repo creation
	ensureGitHubRepo(&cfg)

	// STEP 10 — Configure remote
	if err := configureRemote(&cfg); err != nil {
		// #region agent log
		debugLog("pre-fix", "H5", "internal/setup/setup.go:Run", "configureRemote failed", map[string]interface{}{
			"remote": cfg.Remote,
			"owner":  cfg.Owner,
			"repo":   cfg.Repo,
			"error":  err.Error(),
		})
		// #endregion
		ui.Error("Failed to configure git remote")
		return
	}
	// #region agent log
	debugLog("pre-fix", "H5", "internal/setup/setup.go:Run", "configureRemote succeeded", map[string]interface{}{
		"remote": cfg.Remote,
		"owner":  cfg.Owner,
		"repo":   cfg.Repo,
	})
	// #endregion

	// STEP 11 — Optional first push
	offerFirstPush(&cfg)
	// #region agent log
	debugLog("pre-fix", "H5", "internal/setup/setup.go:Run", "offerFirstPush completed", map[string]interface{}{
		"remote": cfg.Remote,
		"branch": cfg.Branch,
	})
	// #endregion

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

	defaultRepo := filepath.Base(cfg.GetWorkDir())

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
	// Best-effort existence check to avoid noisy 422 errors.
	exists, err := github.RepoExists(cfg.Owner, cfg.Repo)
	if err == nil && exists {
		ui.Success("GitHub repository already exists")
		return
	}

	// If verification failed, fall back to user-driven create flow.
	if !ui.Confirm("Create repository on GitHub if not exists?") {
		return
	}

	private := ui.Confirm("Make repository PRIVATE?")

	createErr := github.CreateRepo(cfg.Owner, cfg.Repo, private)
	if createErr != nil {
		ui.Warn("Repository creation failed:")
		ui.Info(createErr.Error())
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

	// If remote exists, prefer set-url to avoid leaving remote unset on partial failure.
	if _, err := system.GitOutput("remote", "get-url", cfg.Remote); err == nil {
		return system.RunGit("remote", "set-url", cfg.Remote, url)
	}
	return system.RunGit("remote", "add", cfg.Remote, url)
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

	// Avoid pushing without a real first commit.
	// Common case: empty repo => "nothing to commit" after `git add .`.
	stagedNames, err := system.GitOutput("diff", "--cached", "--name-only")
	if err != nil {
		ui.Error("Failed to check staged changes")
		return
	}
	if stagedNames == "" {
		ui.Warn("Nothing to commit yet (working tree is empty)")
		ui.Info("Add at least one file, then push again.")
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
