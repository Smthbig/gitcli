package menu

import (
	"git-genius/internal/gitops"
	"git-genius/internal/ui"
)

func dailyMenu() {
	for {
		ui.Clear()
		ui.BoxHeader("Daily Git Operations")

		ui.BoxMenu("Operations", []string{
			"1) Push changes (commit + push)",
			"2) Pull changes (standard merge)",
			"3) Sync with Rebase (clean history)",
			"4) Smart Pull (auto-stash + pull)",
			"5) Force Push (with lease/safe)",
			"6) Check for Conflicts",
			"7) Fetch all remotes",
			"8) Git status",
			"9) Back",
			"",
			"Tip: h = help",
		})

		switch ui.MenuChoice() {
		case "1":
			track("daily", "push", func() bool { return gitops.Push("") })
		case "2":
			track("daily", "pull", gitops.Pull)
		case "3":
			track("daily", "pull_rebase", gitops.PullRebase)
		case "4":
			track("daily", "smart_pull", gitops.SmartPull)
		case "5":
			track("daily", "force_push", gitops.ForcePush)
		case "6":
			gitops.ShowConflicts()
		case "7":
			track("daily", "fetch", gitops.Fetch)
		case "8":
			track("daily", "status", gitops.Status)
		case "9", "b", "q":
			return
		case "h":
			sectionHelp("Daily Git Operations", ui.HelpDaily)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}
