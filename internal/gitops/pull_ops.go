package gitops

import (
	"git-genius/internal/config"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func Pull() bool {
	return pullCore(false)
}

func PullRebase() bool {
	ui.Info("Pulling with rebase to keep history clean...")
	return pullCore(true)
}

func pullCore(rebase bool) bool {
	if !system.EnsureGitRepo() {
		return false
	}

	cfg := config.Load()
	if cfg.Remote == "" || !system.HasRemote(cfg.Remote) {
		ui.Warn("Remote not configured correctly")
		return false
	}

	branch := CurrentBranch()
	if isWorkingTreeDirty() {
		ui.Warn("Uncommitted changes detected")
		if !ui.Confirm("Continue pull anyway? (Risk of conflicts)") {
			return false
		}
	}

	args := []string{"pull"}
	if rebase {
		args = append(args, "--rebase")
	}
	args = append(args, cfg.Remote, branch)

	stateBefore := InspectRepoState()
	if err := system.RunGitWithRemote(cfg.Remote, args...); err != nil {
		ui.Error("Pull failed")
		if rebase {
			ui.Info("If there are conflicts, resolve them and run 'git rebase --continue'")
		}
		return false
	}

	printPullSuccessSummary(cfg.Remote, branch, stateBefore, true)
	return true
}

func printPullSuccessSummary(remote, branch string, before RepoState, headChanged bool) {
	after := InspectRepoState()
	ui.BoxHeader("Pull Summary")
	ui.Success("Remote : " + remote)
	ui.Success("Branch : " + branch)

	if after.HasAheadBehind {
		ui.Info("Ahead/Behind : " + after.AheadBehindSummary())
	}
}
