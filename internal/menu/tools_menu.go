package menu

import (
	"git-genius/internal/doctor"
	"git-genius/internal/github"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func toolsMenu(gitAvailable bool) {
	for {
		ui.Clear()
		ui.BoxHeader("Tools")

		options := []string{
			"1) Setup / Reconfigure",
			"2) Switch Project",
		}

		if gitAvailable {
			options = append(options,
				"3) Create / Link GitHub Repository",
				"4) GitHub Project Links",
				"5) Git Auth / Credential Helper",
				"6) Doctor (health check)",
				"7) Back",
			)
		} else {
			options = append(options,
				"3) Doctor (health check)",
				"4) Back",
			)
		}
		options = append(options, "", "Tip: h = help")

		ui.BoxMenu("Toolbox", options)

		switch ui.MenuChoice() {
		case "1":
			track("tools", "setup_reconfigure", setup.Run)
		case "2":
			track("tools", "switch_project", setup.SwitchProject)
		case "3":
			if gitAvailable {
				track("tools", "create_or_link_repo", setup.CreateOrLinkRepo)
			} else {
				track("tools", "doctor", doctor.Run)
			}
		case "4":
			if gitAvailable {
				githubLinksMenu()
			} else {
				return
			}
		case "5":
			if gitAvailable {
				track("tools", "configure_git_auth", setup.ConfigureGitAuth)
			} else {
				ui.Error("Invalid option")
			}
		case "6":
			if gitAvailable {
				track("tools", "doctor", doctor.Run)
			} else {
				ui.Error("Invalid option")
			}
		case "7", "b", "q":
			if gitAvailable {
				return
			}
			ui.Error("Invalid option")
		case "h":
			sectionHelp("Tools", ui.HelpTools)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}

func githubLinksMenu() {
	for {
		ui.Clear()
		ui.BoxHeader("GitHub Project Links")

		ui.BoxMenu("Shortcuts", []string{
			"1) Open Repository",
			"2) Open Pull Requests",
			"3) Open Issues",
			"4) Back",
		})

		switch ui.MenuChoice() {
		case "1":
			github.OpenRepo()
		case "2":
			github.OpenPRs()
		case "3":
			github.OpenIssues()
		case "4", "b", "q":
			return
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}
