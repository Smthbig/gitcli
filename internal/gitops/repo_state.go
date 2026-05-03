package gitops

import (
	"fmt"

	"git-genius/internal/config"
	"git-genius/internal/system"
)

type RepoState struct {
	Branch             string
	Remote             string
	HasCommits         bool
	WorkingTreeDirty   bool
	RemoteConfigured   bool
	RemoteTrackingSeen bool
	Ahead              int
	Behind             int
	HasAheadBehind     bool
	FirstRun           bool
	NeedsFirstPush     bool
	HasConflicts       bool
}

func InspectRepoState() RepoState {
	cfg := config.Load()

	state := RepoState{
		Branch:           CurrentBranch(),
		Remote:           cfg.Remote,
		HasCommits:       hasAnyCommit(),
		WorkingTreeDirty: isWorkingTreeDirty(),
		FirstRun:         !config.HasProjectConfig(cfg.GetWorkDir()) && !config.HasHistoryForWorkDir(cfg.GetWorkDir()),
	}

	// Check for conflicts
	if out, err := system.GitOutput("diff", "--name-only", "--diff-filter=U"); err == nil && out != "" {
		state.HasConflicts = true
	}

	if state.Branch == "-" || state.Branch == "" {
		state.Branch = cfg.Branch
	}

	if state.Remote != "" && system.HasRemote(state.Remote) {
		state.RemoteConfigured = true
	}

	if state.Branch != "" && state.RemoteConfigured && system.HasRemoteTrackingBranch(state.Remote, state.Branch) {
		state.RemoteTrackingSeen = true
		ahead, behind, err := system.AheadBehind(state.Branch, state.Remote+"/"+state.Branch)
		if err == nil {
			state.Ahead = ahead
			state.Behind = behind
			state.HasAheadBehind = true
		}
	}

	state.NeedsFirstPush = state.HasCommits && state.RemoteConfigured && state.Branch != "" && !state.RemoteTrackingSeen
	return state
}

func (s RepoState) AheadBehindSummary() string {
	if !s.HasAheadBehind {
		return ""
	}
	return fmt.Sprintf("%d ahead / %d behind", s.Ahead, s.Behind)
}
