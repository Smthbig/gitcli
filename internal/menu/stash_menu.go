package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/ui"
)

func stashMenu() {
	for {
		ui.Clear()
		ui.BoxHeader("Stash & Undo")

		ui.BoxMenu("Options", []string{
			"1) Stash changes",
			"2) List stashes",
			"3) Apply last stash (pop)",
			"4) Undo last commit (keep changes)",
			"5) Back",
			"",
			"Tip: h = help",
		})

		switch ui.MenuChoice() {
		case "1":
			track("stash", "stash_save", gitops.StashSave)
		case "2":
			track("stash", "stash_list", gitops.StashList)
		case "3":
			track("stash", "stash_pop", gitops.StashPop)
		case "4":
			track("stash", "undo_last_commit", gitops.UndoLastCommit)
		case "5", "b", "q":
			return
		case "h":
			sectionHelp("Stash & Undo", ui.HelpStash)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}
