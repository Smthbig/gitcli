package menu

import (
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func getMainMenuOptions() []string {
	return []string{
		"1) Daily Git Operations",
		"2) Visual History (Graph)",
		"3) Activity Timeline (Chart)",
		"4) Project Insights (Stats)",
		"5) Branch / Remote",
		"6) Stash & Undo",
		"7) Tools",
		"8) Help / About",
		"9) Exit",
		"",
		"Quick Keys: [p]ush, [u]pdate, [s]tatus",
	}
}

func mainMenu() bool {
	// Interaction logic remains here, but printing is handled by the main loop
	switch ui.MenuChoice() {
	case "1":
		dailyMenu()
	case "2":
		gitops.ShowGraph()
		ui.Pause()
	case "3":
		gitops.ShowActivityTimeline()
		ui.Pause()
	case "4":
		gitops.ShowRepoStats()
		ui.Pause()
	case "5":
		branchMenu()
	case "6":
		stashMenu()
	case "7":
		toolsMenu(true)
	case "8", "h":
		mainHelp()
	case "9", "q":
		ui.Info("Goodbye")
		return false
	case "p":
		track("daily", "push", func() bool { return gitops.Push("") })
		ui.Pause()
	case "u":
		track("daily", "pull", gitops.Pull)
		ui.Pause()
	case "s":
		track("daily", "status", gitops.Status)
		ui.Pause()
	default:
		ui.Error("Invalid option")
		ui.Pause()
	}

	return true
}
func getLimitedMenuOptions() []string {
	return []string{
		"1) Setup / Reconfigure",
		"2) Switch Project",
		"3) Doctor (health check)",
		"4) Help / About",
		"5) Exit",
		"",
		"Tip: install Git to unlock more",
	}
}

func limitedMenu() bool {
	switch ui.MenuChoice() {
	case "1":
		track("tools", "setup_reconfigure", setup.Run)
	case "2":
		track("tools", "switch_project", setup.SwitchProject)
	case "3":
		track("tools", "doctor", doctor.Run)
	case "4", "h":
		mainHelp()
	case "5", "q":
		ui.Info("Goodbye")
		return false
	default:
		ui.Error("Invalid option")
		ui.Pause()
	}

	return true
}
