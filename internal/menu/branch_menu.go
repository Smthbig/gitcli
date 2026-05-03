package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/menu/tui"
	"git-genius/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func branchMenu() {
	items := []tui.MenuItem{
		{Label: "Switch to existing branch", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("branch", "switch_branch", gitops.SwitchBranch); ui.Pause(); return true }} }},
		{Label: "Create new branch", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("branch", "create_branch", gitops.CreateBranch); ui.Pause(); return true }} }},
		{Label: "Cleanup merged branches", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("branch", "cleanup_branches", gitops.CleanupMergedBranches); ui.Pause(); return true }} }},
		{Label: "Configure remote", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("branch", "switch_remote", gitops.SwitchRemote); ui.Pause(); return true }} }},
		{Label: "Back", Action: func() tea.Msg { return tea.Quit() }},
	}

	RunSubMenu("BRANCH / REMOTE", items, true)
}
