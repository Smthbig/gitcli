package menu

import (
	"git-genius/internal/doctor"
	"git-genius/internal/github"
	"git-genius/internal/menu/tui"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func toolsMenu(gitAvailable bool) {
	var items []tui.MenuItem
	items = append(items, tui.MenuItem{Label: "Setup / Reconfigure", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "setup_reconfigure", setup.Run); ui.Pause(); return true }} }})
	items = append(items, tui.MenuItem{Label: "Switch Project", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "switch_project", setup.SwitchProject); ui.Pause(); return true }} }})

	if gitAvailable {
		items = append(items,
			tui.MenuItem{Label: "Create / Link GitHub Repository", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "create_or_link_repo", setup.CreateOrLinkRepo); ui.Pause(); return true }} }},
			tui.MenuItem{Label: "GitHub Project Links", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { githubLinksMenu(); return true }} }},
			tui.MenuItem{Label: "Git Auth / Credential Helper", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "configure_git_auth", setup.ConfigureGitAuth); ui.Pause(); return true }} }},
			tui.MenuItem{Label: "Doctor (health check)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "doctor", doctor.Run); ui.Pause(); return true }} }},
		)
	} else {
		items = append(items, tui.MenuItem{Label: "Doctor (health check)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("tools", "doctor", doctor.Run); ui.Pause(); return true }} }})
	}
	items = append(items, tui.MenuItem{Label: "Back", Action: func() tea.Msg { return tea.Quit() }})

	RunSubMenu("TOOLS", items, gitAvailable)
}

func githubLinksMenu() {
	items := []tui.MenuItem{
		{Label: "Open Repository", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { github.OpenRepo(); ui.Pause(); return true }} }},
		{Label: "Open Pull Requests", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { github.OpenPRs(); ui.Pause(); return true }} }},
		{Label: "Open Issues", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { github.OpenIssues(); ui.Pause(); return true }} }},
		{Label: "Back", Action: func() tea.Msg { return tea.Quit() }},
	}

	RunSubMenu("GITHUB LINKS", items, true)
}
