package gitops

import (
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func Push(msg string) bool {
	return pushCore(msg, false)
}

func ForcePush() bool {
	ui.Warn("Force pushing can overwrite remote history!")
	if !ui.Confirm("Are you sure you want to force push (with lease)?") {
		return false
	}
	return pushCore("", true)
}

func pushCore(msg string, force bool) bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ensureSafeDirectory()
	cfg := config.Load()
	stateBefore := InspectRepoState()
	createdCommit := false

	if !hasAnyCommit() {
		if cfg.Branch != "" && system.CurrentGitBranch() == "" {
			_ = system.PrepareBranch(cfg.Branch)
		}
		if msg == "" {
			msg = ComposeConventionalMessage("Initial commit")
		}
		ui.Info("Creating initial commit")
		_ = system.RunGit("add", ".")
		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Initial commit failed")
			return false
		}
		createdCommit = true
	} else if isWorkingTreeDirty() {
		if msg == "" {
			msg = ComposeConventionalMessage("Update")
		}
		_ = system.RunGit("add", ".")
		if err := system.RunGit("commit", "-m", msg); err != nil {
			ui.Error("Commit failed")
			return false
		}
		createdCommit = true
	}

	if cfg.Remote == "" || !system.HasRemote(cfg.Remote) {
		ui.Warn("Remote not configured correctly")
		return false
	}

	branch := CurrentBranch()
	args := []string{"push", "-u", cfg.Remote, branch}
	if force {
		args = append(args, "--force-with-lease")
	}

	stderr, err := system.RunGitWithRemoteBuffered(cfg.Remote, args...)
	if err != nil {
		ui.Error("Push failed")
		printPushFailureSummary(cfg.Remote, branch, stderr)
		return false
	}

	cfg.FirstPushDone = true
	config.Save(cfg)
	printPushSuccessSummary(cfg.Remote, branch, stateBefore, createdCommit)
	return true
}

func printPushSuccessSummary(remote, branch string, before RepoState, createdCommit bool) {
	after := InspectRepoState()
	ui.BoxHeader("Push Summary")
	ui.Success("Remote : " + remote)
	ui.Success("Branch : " + branch)

	if after.HasAheadBehind {
		ui.Info("Ahead/Behind : " + after.AheadBehindSummary())
	}
}

func printPushFailureSummary(remote, branch, stderr string) {
	ui.BoxHeader("Push Summary")
	ui.Error("Remote : " + remote)
	ui.Error("Branch : " + branch)
	if strings.Contains(strings.ToLower(stderr), "rejected") {
		ui.Warn("Push rejected: remote has changes you don't have.")
		ui.Info("Try pulling changes first, or use Force Push if you know what you're doing.")
	}
}
