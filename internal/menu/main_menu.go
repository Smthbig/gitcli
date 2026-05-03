package menu

import (
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func mainMenu() bool {
	ui.BoxMenu("Main Menu", []string{
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
		"Tip: press 'h' for help",
	})

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
	default:
		ui.Error("Invalid option")
		ui.Pause()
	}

	return true
}

func limitedMenu() bool {
	ui.BoxMenu("Limited Mode", []string{
		"1) Setup / Reconfigure",
		"2) Switch Project",
		"3) Doctor (health check)",
		"4) Help / About",
		"5) Exit",
		"",
		"Tip: install Git to unlock more",
	})

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
