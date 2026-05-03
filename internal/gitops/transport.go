package gitops

import (
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/system"
)

func hasAnyCommit() bool {
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

func CurrentBranch() string {
	branch, err := system.GitOutput("branch", "--show-current")
	if err != nil || branch == "" {
		return "-"
	}
	return branch
}

func Status() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	return system.RunGit("status") == nil
}

func Fetch() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	remotes, err := system.RemoteNames()
	if err != nil {
		return false
	}

	for _, remote := range remotes {
		_ = system.RunGitWithRemote(remote, "fetch", remote)
	}
	return true
}
