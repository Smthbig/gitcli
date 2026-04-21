package gitops

import (
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

///////////////////////////////////////////////////////////////
// INTERNAL HELPERS (ANDROID SAFE)
///////////////////////////////////////////////////////////////

func hasAnyCommit() bool {
	// log -1 is safer than rev-parse on Android
	cmd, err := system.GitCmd("log", "-1")
	if err != nil {
		return false
	}
	return cmd.Run() == nil
}

func isWorkingTreeDirty() bool {
	out, err := system.GitOutput("status", "--porcelain")
	if err != nil {
		return false
	}
	return strings.TrimSpace(out) != ""
}

func ensureSafeDirectory() {
	cfg := config.Load()
	if cfg.WorkDir == "" {
		return
	}
	system.EnsureSafeDirectory(cfg.WorkDir)
}

///////////////////////////////////////////////////////////////
// CONTEXT HELPERS
///////////////////////////////////////////////////////////////

func CurrentBranch() string {
	branch, err := system.GitOutput("branch", "--show-current")
	if err != nil || branch == "" {
		return "-"
	}
	return branch
}

///////////////////////////////////////////////////////////////
// CORE OPERATIONS
///////////////////////////////////////////////////////////////

func Status() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	if err := system.RunGit("status"); err != nil {
		ui.Error("Git status failed")
		return false
	}
	return true
}

func Push(msg string) bool {

	if !system.EnsureGitRepo() {
		return false
	}

	ensureSafeDirectory()
	cfg := config.Load()

	// ---------- FIRST COMMIT ----------
	if !hasAnyCommit() {
		if cfg.Branch != "" && system.CurrentGitBranch() == "" {
			if err := system.PrepareBranch(cfg.Branch); err != nil {
				ui.Error("Failed to prepare branch for first commit")
				ui.Info(err.Error())
				return false
			}
		}

		if msg == "" {
			msg = "Initial commit"
		}

		ui.Info("Creating initial commit")

		if err := system.RunGit("add", "."); err != nil {
			ui.Error("git add failed")
			return false
		}

		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Initial commit failed")
			return false
		}
	} else if isWorkingTreeDirty() {
		// ---------- NORMAL COMMIT ----------
		if msg == "" {
			ui.Error("Commit message required")
			return false
		}

		if err := system.RunGit("add", "."); err != nil {
			ui.Error("git add failed")
			return false
		}

		// Ignore nothing-to-commit safely
		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Commit failed")
			return false
		}
	} else {
		ui.Info("No local file changes detected")
		ui.Info("Attempting to push any existing local commits")
	}

	// ---------- REMOTE CHECK ----------
	if cfg.Remote == "" {
		ui.Warn("No remote configured")
		ui.Info("Run Branch / Remote -> Configure remote or Tools -> Create / Link GitHub Repository")
		return false
	}

	if !system.HasRemote(cfg.Remote) {
		ui.Warn("Configured remote not found: " + cfg.Remote)
		ui.Info("Run Branch / Remote -> Configure remote")
		return false
	}

	branch := CurrentBranch()
	if branch == "-" {
		branch = cfg.Branch
	}

	if branch == "" {
		ui.Error("No branch detected")
		return false
	}

	// ---------- PUSH ----------
	if err := system.RunGit("push", "-u", cfg.Remote, branch); err != nil {
		ui.Error("Push failed")
		if !system.HasGitCredentialHelper() {
			ui.Info("Run Tools -> Git Auth / Credential Helper to reduce repeated HTTPS auth prompts")
		}
		ui.Info("Run Doctor if the problem persists")
		return false
	}

	ui.Success("Changes pushed successfully")
	return true
}

func Pull() bool {

	if !system.EnsureGitRepo() {
		return false
	}

	cfg := config.Load()
	if cfg.Remote == "" {
		ui.Warn("No remote configured")
		ui.Info("Run Branch / Remote -> Configure remote")
		return false
	}
	if !system.HasRemote(cfg.Remote) {
		ui.Warn("Configured remote not found: " + cfg.Remote)
		ui.Info("Run Branch / Remote -> Configure remote")
		return false
	}

	branch := CurrentBranch()
	if branch == "-" {
		branch = cfg.Branch
	}
	if branch == "" {
		ui.Error("No branch detected")
		return false
	}

	if isWorkingTreeDirty() {
		ui.Warn("Uncommitted changes detected")
		ui.Info("Use Smart Pull if you want Git Genius to stash and restore changes automatically")
		if !ui.Confirm("Continue with a normal pull anyway?") {
			ui.Warn("Pull cancelled")
			return false
		}
	}

	if err := system.RunGit("pull", cfg.Remote, branch); err != nil {
		ui.Error("Pull failed")
		ui.Info("Try Smart Pull or run Doctor for more guidance")
		return false
	}

	ui.Success("Pull completed")
	return true
}

func Fetch() bool {

	if !system.EnsureGitRepo() {
		return false
	}

	if err := system.RunGit("fetch", "--all"); err != nil {
		ui.Error("Fetch failed")
		return false
	}

	ui.Success("Fetched all remotes")
	return true
}
