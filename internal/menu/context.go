package menu

import (
	"fmt"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/gitops"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func headerTitle(appVersion string, gitAvailable bool) string {
	version := strings.TrimSpace(appVersion)
	if version == "" {
		version = "dev"
	}

	title := "Git Genius v" + version
	if !gitAvailable {
		title += " (Limited Mode)"
	}
	return title
}

func GetStatusLines(gitAvailable bool) []string {
	cfg := config.Load()
	projectDir := cfg.GetWorkDir()

	var info []string
	branch := "-"
	remote := system.CurrentRemote()
	sync := "0 ahead / 0 behind"
	
	if gitAvailable {
		state := gitops.InspectRepoState()
		branch = state.Branch
		if state.HasAheadBehind {
			sync = state.AheadBehindSummary()
		}
	}

	info = append(info, fmt.Sprintf("Project : %s", filepath.Base(projectDir)))
	info = append(info, fmt.Sprintf("Branch  : %s", branch))
	info = append(info, fmt.Sprintf("Remote  : %s", remote))
	info = append(info, fmt.Sprintf("Sync    : %s", sync))
	
	// Dynamic Health Check & Assistant
	statusText := "HEALTHY"
	nextAction := "Ready for work"
	
	if !gitAvailable {
		statusText = "LIMITED"
		nextAction = "Install Git"
	} else {
		state := gitops.InspectRepoState()
		if state.HasConflicts {
			statusText = "CONFLICT"
			nextAction = "Fix Merge Conflicts"
		} else if state.WorkingTreeDirty {
			nextAction = "Commit your changes"
		} else if state.Ahead > 0 {
			nextAction = fmt.Sprintf("Push %d commits", state.Ahead)
		} else if state.Behind > 0 {
			nextAction = "Pull latest changes"
		} else if state.NeedsFirstPush {
			nextAction = "Publish repository"
		}
	}

	info = append(info, "")
	info = append(info, fmt.Sprintf("Status  : %s", statusText))
	info = append(info, fmt.Sprintf("Next    : %s", nextAction))
	
	data := gitops.GetRecentActivityData(7)
	if data != nil {
		spark := ui.CyberSparkline(data)
		info = append(info, "")
		info = append(info, fmt.Sprintf("Trend   : %s", spark))
	}

	return info
}

func track(section, action string, fn func() bool) {
	cfg := config.Load()
	ok := false
	note := ""
	defer func() {
		if r := recover(); r != nil {
			ok = false
			note = fmt.Sprintf("panic: %v", r)
			ui.Error("Unexpected error occurred")
		}
		config.RecordHistory(cfg.GetWorkDir(), section, action, ok, note)
	}()
	ok = fn()
}
