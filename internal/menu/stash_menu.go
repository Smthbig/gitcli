package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/menu/tui"
	"git-genius/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func stashMenu() {
	items := []tui.MenuItem{
		{Label: "Stash changes", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("stash", "stash_save", gitops.StashSave); ui.Pause(); return true }} }},
		{Label: "List stashes", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("stash", "stash_list", gitops.StashList); ui.Pause(); return true }} }},
		{Label: "Apply last stash (pop)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("stash", "stash_pop", gitops.StashPop); ui.Pause(); return true }} }},
		{Label: "Undo last commit (keep changes)", Action: func() tea.Msg { return tui.ActionMsg{Action: func() bool { track("stash", "undo_last_commit", gitops.UndoLastCommit); ui.Pause(); return true }} }},
		{Label: "Back", Action: func() tea.Msg { return tea.Quit() }},
	}

	RunSubMenu("STASH & UNDO", items, true)
}
