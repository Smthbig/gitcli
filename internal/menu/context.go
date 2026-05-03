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

func showContext(gitAvailable bool) {
	cfg := config.Load()
	projectDir := cfg.GetWorkDir()

	width := 42
	fmt.Println(ui.Cyan + "┏" + strings.Repeat("━", width-2) + "┓" + ui.Reset)
	fmt.Printf(ui.Cyan+"┃ "+ui.Reset+ui.Bold+"PROJECT: "+ui.Reset+"%-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-13, filepath.Base(projectDir))
	fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"Path   : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, projectDir)

	if gitAvailable {
		state := gitops.InspectRepoState()
		fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"Branch : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, state.Branch)
		fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"Remote : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, system.CurrentRemote())
		if state.HasAheadBehind {
			fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"Sync   : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, state.AheadBehindSummary())
		}
	} else {
		fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"Mode   : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, "Limited (Git unavailable)")
	}

	if cfg.Owner != "" && cfg.Repo != "" {
		repoURL := "github.com/" + cfg.Owner + "/" + cfg.Repo
		fmt.Printf(ui.Cyan+"┃ "+ui.Reset+"GitHub : %-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-11, repoURL)
	}
	
	suggestions := config.HistorySuggestions(projectDir)
	if len(suggestions) > 0 {
		fmt.Println(ui.Cyan + "┣" + strings.Repeat("━", width-2) + "┫" + ui.Reset)
		for _, s := range suggestions {
			fmt.Printf(ui.Cyan+"┃ "+ui.Reset+ui.Yellow+"󱐌 "+ui.Reset+"%-*s"+ui.Cyan+" ┃\n"+ui.Reset, width-13, s)
		}
	}
	fmt.Println(ui.Cyan + "┗" + strings.Repeat("━", width-2) + "┛" + ui.Reset)
}

func normalizeChoice(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
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
