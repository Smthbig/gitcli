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

func Status() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("status"); err != nil {
		ui.Error("Git status failed")
	}
}

func Push(msg string) {

	if !system.EnsureGitRepo() {
		return
	}

	ensureSafeDirectory()
	cfg := config.Load()

	// ---------- FIRST COMMIT ----------
	if !hasAnyCommit() {

		if msg == "" {
			msg = "Initial commit"
		}

		ui.Info("Creating initial commit")

		if err := system.RunGit("add", "."); err != nil {
			ui.Error("git add failed")
			return
		}

		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Initial commit failed")
			return
		}
	} else if isWorkingTreeDirty() {
		// ---------- NORMAL COMMIT ----------
		if msg == "" {
			ui.Error("Commit message required")
			return
		}

		if err := system.RunGit("add", "."); err != nil {
			ui.Error("git add failed")
			return
		}

		// Ignore nothing-to-commit safely
		_ = system.RunGit("commit", "-m", msg)
	} else {
		ui.Warn("Nothing to commit")
	}

	// ---------- REMOTE CHECK ----------
	if cfg.Remote == "" {
		ui.Warn("No remote configured")
		ui.Info("Run setup to configure GitHub repository")
		return
	}

	branch := CurrentBranch()
	if branch == "-" {
		branch = cfg.Branch
	}

	if branch == "" {
		ui.Error("No branch detected")
		return
	}

	// ---------- PUSH ----------
	if err := system.RunGit("push", "-u", cfg.Remote, branch); err != nil {
		ui.Error("Push failed")
		return
	}

	ui.Success("Changes pushed successfully")
}

func Pull() {

	if !system.EnsureGitRepo() {
		return
	}

	cfg := config.Load()
	if cfg.Remote == "" {
		ui.Warn("No remote configured")
		return
	}

	branch := CurrentBranch()
	if branch == "-" {
		branch = cfg.Branch
	}

	if err := system.RunGit("pull", cfg.Remote, branch); err != nil {
		ui.Error("Pull failed")
		return
	}

	ui.Success("Pull completed")
}

func Fetch() {

	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("fetch", "--all"); err != nil {
		ui.Error("Fetch failed")
		return
	}

	ui.Success("Fetched all remotes")
}
