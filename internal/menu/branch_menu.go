package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/ui"
)

func branchMenu() {
	for {
		ui.Clear()
		ui.BoxHeader("Branch / Remote")

		ui.BoxMenu("Options", []string{
			"1) Switch to existing branch",
			"2) Create new branch",
			"3) Cleanup merged branches",
			"4) Configure remote",
			"5) Back",
			"",
			"Tip: h = help",
		})

		switch ui.MenuChoice() {
		case "1":
			track("branch", "switch_branch", gitops.SwitchBranch)
		case "2":
			track("branch", "create_branch", gitops.CreateBranch)
		case "3":
			track("branch", "cleanup_branches", gitops.CleanupMergedBranches)
		case "4":
			track("branch", "switch_remote", gitops.SwitchRemote)
		case "5", "b", "q":
			return
		case "h":
			sectionHelp("Branch / Remote", ui.HelpBranch)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}
