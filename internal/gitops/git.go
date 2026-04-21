package gitops

import (
	"fmt"
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
	stateBefore := InspectRepoState()
	createdCommit := false

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

		stagedNames, err := system.GitOutput("diff", "--cached", "--name-only")
		if err != nil {
			ui.Error("Failed to inspect staged files")
			return false
		}
		if stagedNames == "" {
			ui.Warn("This repository has no files to commit yet")
			ui.Info("Create at least one file, then run Push changes again")
			if cfg.Remote != "" {
				ui.Info("Your remote can stay configured now; the first push just needs a real commit")
			}
			return false
		}

		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Initial commit failed")
			return false
		}
		createdCommit = true
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
		createdCommit = true
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

	if stateBefore.NeedsFirstPush || (!createdCommit && stateBefore.HasCommits && !stateBefore.RemoteTrackingSeen) {
		ui.Info("This looks like the first push for " + cfg.Remote + "/" + branch)
	}

	// ---------- PUSH ----------
	stderr, err := system.RunGitWithRemoteBuffered(cfg.Remote, "push", "-u", cfg.Remote, branch)
	if err != nil {
		ui.Error("Push failed")
		printPushFailureSummary(cfg.Remote, branch, stderr)
		ui.Info("Run Doctor if the problem persists")
		return false
	}

	cfg.FirstPushDone = true
	config.Save(cfg)
	printPushSuccessSummary(cfg.Remote, branch, stateBefore, createdCommit)
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

	stateBefore := InspectRepoState()
	headBefore, _ := system.GitOutput("rev-parse", "HEAD")

	if err := system.RunGitWithRemote(cfg.Remote, "pull", cfg.Remote, branch); err != nil {
		ui.Error("Pull failed")
		ui.Info("Try Smart Pull or run Doctor for more guidance")
		return false
	}

	headAfter, _ := system.GitOutput("rev-parse", "HEAD")
	printPullSuccessSummary(cfg.Remote, branch, stateBefore, headBefore != "" && headBefore != headAfter)
	return true
}

func Fetch() bool {

	if !system.EnsureGitRepo() {
		return false
	}

	remotes, err := system.RemoteNames()
	if err != nil {
		ui.Error("Fetch failed")
		return false
	}

	for _, remote := range remotes {
		if err := system.RunGitWithRemote(remote, "fetch", remote); err != nil {
			ui.Error("Fetch failed")
			return false
		}
	}

	ui.Success("Fetched all remotes")
	return true
}

func printPushSuccessSummary(remote, branch string, before RepoState, createdCommit bool) {
	after := InspectRepoState()

	ui.Header("Push Summary")
	ui.Success("Remote : " + remote)
	ui.Success("Branch : " + branch)

	switch {
	case !before.HasCommits && createdCommit:
		ui.Success("Created the first commit and published the branch")
	case before.NeedsFirstPush:
		ui.Success("Published this branch to the remote for the first time")
	case before.HasAheadBehind && before.Ahead > 0:
		ui.Success(fmt.Sprintf("Published %d local commit(s)", before.Ahead))
	case createdCommit:
		ui.Success("Created a new commit and pushed it")
	default:
		ui.Info("Remote branch was already up to date before this push")
	}

	if after.HasAheadBehind {
		ui.Info("Ahead/Behind : " + after.AheadBehindSummary())
	}
}

func printPushFailureSummary(remote, branch, stderr string) {
	ui.Header("Push Summary")
	ui.Error("Remote : " + remote)
	ui.Error("Branch : " + branch)

	if system.RemoteUsesHTTPS(remote) {
		ui.Info("HTTPS remote detected")
		ui.Info("Run Tools -> Git Auth / Credential Helper and preload the current GitHub token into Git")
		ui.Info("Check that the token still has repo access to this repository")
		if !system.HasGitCredentialHelper() {
			ui.Info("No Git credential helper is configured yet")
		}
	}

	if strings.Contains(strings.ToLower(stderr), "repository not found") {
		ui.Info("Verify the remote URL and GitHub owner/repository in Setup or Create / Link GitHub Repository")
	}
}

func printPullSuccessSummary(remote, branch string, before RepoState, headChanged bool) {
	after := InspectRepoState()

	ui.Header("Pull Summary")
	ui.Success("Remote : " + remote)
	ui.Success("Branch : " + branch)

	switch {
	case before.HasAheadBehind && before.Behind > 0:
		ui.Success(fmt.Sprintf("Integrated %d remote commit(s)", before.Behind))
	case headChanged:
		ui.Success("Local branch moved forward")
	default:
		ui.Info("Local branch was already up to date")
	}

	if after.HasAheadBehind {
		ui.Info("Ahead/Behind : " + after.AheadBehindSummary())
	}
}
