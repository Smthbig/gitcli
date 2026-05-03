package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/menu/tui"
	"git-genius/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func dailyMenu() {
	items := []tui.MenuItem{
		{Label: "Push changes (commit + push)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "push", func() bool { return gitops.Push("") }); ui.Pause(); return true }} }},
		{Label: "Pull changes (standard merge)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "pull", gitops.Pull); ui.Pause(); return true }} }},
		{Label: "Sync with Rebase (clean history)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "pull_rebase", gitops.PullRebase); ui.Pause(); return true }} }},
		{Label: "Smart Pull (auto-stash + pull)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "smart_pull", gitops.SmartPull); ui.Pause(); return true }} }},
		{Label: "Force Push (with lease/safe)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "force_push", gitops.ForcePush); ui.Pause(); return true }} }},
		{Label: "Check for Conflicts", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { gitops.ShowConflicts(); ui.Pause(); return true }} }},
		{Label: "Fetch all remotes", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "fetch", gitops.Fetch); ui.Pause(); return true }} }},
		{Label: "Git status", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("daily", "status", gitops.Status); ui.Pause(); return true }} }},
		{Label: "Back", Action: func() tea.Msg { return tea.Quit() }},
	}

	RunSubMenu("DAILY OPERATIONS", items, true)
}
