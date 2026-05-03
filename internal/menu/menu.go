package menu

import (
	"fmt"
	"os"

	"git-genius/internal/config"
	"git-genius/internal/gitops"
	"git-genius/internal/menu/tui"
	"git-genius/internal/setup"
	"git-genius/internal/system"
	"git-genius/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

/*
Start is the main entry point for the Git Genius menu system.
It handles initialization and the main loop.
*/
func Start(appVersion string, gitAvailable bool) {
	maybeOfferSetup(gitAvailable)

	// Define menu items
	var items []tui.MenuItem
	if gitAvailable {
		items = []tui.MenuItem{
			{Label: "Daily Git Operations", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { dailyMenu(); return true }} }},
			{Label: "Visual History (Graph)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { gitops.ShowGraph(); ui.Pause(); return true }} }},
			{Label: "Activity Timeline (Chart)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { gitops.ShowActivityTimeline(); ui.Pause(); return true }} }},
			{Label: "Project Insights (Stats)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { gitops.ShowRepoStats(); ui.Pause(); return true }} }},
			{Label: "Branch / Remote", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { branchMenu(); return true }} }},
			{Label: "Stash & Undo", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { stashMenu(); return true }} }},
			{Label: "Tools", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { toolsMenu(true); return true }} }},
			{Label: "Help / About", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { mainHelp(); return true }} }},
			{Label: "Exit", Action: func() tea.Msg { return tea.Quit() }},
		}
	} else {
		items = []tui.MenuItem{
			{Label: "Setup / Reconfigure", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "setup_reconfigure", setup.Run); return true }} }},
			{Label: "Switch Project", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "switch_project", setup.SwitchProject); return true }} }},
			{Label: "Help / About", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { mainHelp(); return true }} }},
			{Label: "Exit", Action: func() tea.Msg { return tea.Quit() }},
		}
	}

	for {
		m := tui.NewModel(appVersion, gitAvailable, items)
		m.Dashboard = GetStatusLines(gitAvailable)

		// Wire quick keys
		if gitAvailable {
			m.OnPush = func() bool { track("daily", "push", func() bool { return gitops.Push("") }); ui.Pause(); return true }
			m.OnUpdate = func() bool { track("daily", "pull", gitops.Pull); ui.Pause(); return true }
			m.OnStatus = func() bool { track("daily", "status", gitops.Status); ui.Pause(); return true }
		}

		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error running UI: %v\n", err)
			os.Exit(1)
		}

		// Handle actions after TUI returns
		if m, ok := finalModel.(tui.Model); ok {
			if m.Quitting {
				ui.Info("Goodbye")
				return
			}
			if m.ActiveAction != nil {
				m.ActiveAction()
			}
		}
	}
}

func maybeOfferSetup(gitAvailable bool) {
	if !gitAvailable {
		return
	}

	cfg := config.Load()
	if config.HasProjectConfig(cfg.GetWorkDir()) || config.HasHistoryForWorkDir(cfg.GetWorkDir()) {
		return
	}

	if system.IsGitRepo() {
		if gitops.InspectRepoState().HasCommits {
			ui.Info("First run for this repository: Git Genius has not been configured here yet")
		} else {
			ui.Info("Brand-new repository detected: no commits yet and no Git Genius setup found")
		}
	} else {
		ui.Info("Brand-new project directory detected: setup can initialize Git and prepare the first push")
	}
	if setup.Run() {
		ui.Pause()
	}
}
