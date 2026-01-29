package gitops

import (
	"bytes"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/* ============================================================
   INTERNAL HELPERS (ANDROID SAFE)
   ============================================================ */

// hasAnyCommit checks whether repo has at least one commit
// Uses git log (rev-parse is unsafe on some Android kernels)
func hasAnyCommit() bool {
	cmd := system.GitCmd("log", "-1")
	return cmd.Run() == nil
}

// isWorkingTreeDirty checks for uncommitted changes
func isWorkingTreeDirty() bool {
	var out bytes.Buffer
	cmd := system.GitCmd("status", "--porcelain")
	cmd.Stdout = &out
	_ = cmd.Run()
	return out.Len() > 0
}

// ensureSafeDirectory fixes Git ≥2.35 "dubious ownership" (Android fix)
func ensureSafeDirectory() {
	cfg := config.Load()
	if cfg.WorkDir == "" {
		return
	}
	_ = system.RunGit("config", "--global", "--add", "safe.directory", cfg.WorkDir)
}

/* ============================================================
   CONTEXT HELPERS (REAL DATA)
   ============================================================ */

func CurrentBranch() string {
	var out bytes.Buffer
	cmd := system.GitCmd("branch", "--show-current")
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "-"
	}

	b := strings.TrimSpace(out.String())
	if b == "" {
		return "-"
	}
	return b
}

func CurrentRemote() string {
	cfg := config.Load()
	if cfg.Remote == "" {
		return "-"
	}

	cmd := system.GitCmd("remote", "get-url", cfg.Remote)
	if err := cmd.Run(); err != nil {
		return "-"
	}

	return cfg.Remote
}

/* ============================================================
   CORE GIT OPERATIONS
   ============================================================ */

func Status() {
	if !system.EnsureGitRepo() {
		return
	}

	if err := system.RunGit("status"); err != nil {
		ui.Error("Failed to get git status")
	}
}

func Push(msg string) {
	if !system.EnsureGitRepo() {
		return
	}

	// 🔐 Android / Git ≥2.35 safety
	ensureSafeDirectory()

	cfg := config.Load()

	// ---------- NO CHANGES ----------
	if hasAnyCommit() && !isWorkingTreeDirty() {
		ui.Warn("Nothing to commit")
		return
	}

	// ---------- FIRST COMMIT ----------
	if !hasAnyCommit() {
		if msg == "" {
			msg = "Initial commit"
		}

		ui.Info("Creating first commit")
		_ = system.RunGit("add", ".")

		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Initial commit failed")
			return
		}

		ui.Success("Initial commit created")

		if cfg.Remote == "" {
			ui.Warn("No remote configured")
			ui.Info("Run: Tools → Create / Link GitHub Repository")
			return
		}
	}

	// ---------- NORMAL COMMIT ----------
	if msg == "" {
		ui.Error("Commit message cannot be empty")
		return
	}

	_ = system.RunGit("add", ".")
	_ = system.RunGit("commit", "-m", msg) // ignore "nothing to commit"

	// ---------- PUSH ----------
	if cfg.Remote == "" {
		ui.Warn("No remote configured")
		ui.Info("Run: Tools → Create / Link GitHub Repository")
		return
	}

	branch := CurrentBranch()
	if branch == "-" {
		branch = cfg.Branch
	}

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
