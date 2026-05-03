package menu

import (
	"git-genius/internal/ui"
)

func mainHelp() {
	ui.Clear()
	ui.BoxHeader("Help / About Git Genius")

	ui.PrintHelp(ui.HelpMain)
	ui.PrintHelp(ui.HelpWorkflow)
	ui.PrintHelp(ui.HelpGitHub)
	ui.PrintHelp(ui.HelpTroubleshooting)

	ui.Pause()
}

func sectionHelp(title string, help []string) {
	ui.Clear()
	ui.BoxHeader(title + " – Help")
	ui.PrintHelp(help)
	ui.Pause()
}
