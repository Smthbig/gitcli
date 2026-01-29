package menu

import (
	"fmt"
	"os"
	"path/filepath"

	"git-genius/internal/config"
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/ui"
)

func Start() {
	for {
		ui.Clear()
		ui.Header("Git Genius v1.0")

		showContext()

		fmt.Println("1) Daily Git Operations")
		fmt.Println("2) Branch / Remote")
		fmt.Println("3) Stash & Undo")
		fmt.Println("4) Tools")
		fmt.Println("5) Help / About")
		fmt.Println("6) Exit")
		fmt.Println()
		fmt.Println("Tip: press 'h' for help")

		switch ui.Input("Select option") {
		case "1":
			dailyMenu()
		case "2":
			branchMenu()
		case "3":
			stashMenu()
		case "4":
			toolsMenu()
		case "5", "h", "help", "?":
			mainHelp()
		case "6":
			ui.Info("Goodbye 👋")
			os.Exit(0)
		default:
			ui.Error("Invalid option")
			ui.Pause()
		}
	}
}

/* ============================================================
   Context Panel
   ============================================================ */

func showContext() {
	cfg := config.Load()
	projectDir := cfg.GetWorkDir()

	fmt.Println("Project :", filepath.Base(projectDir))
	fmt.Println("Path    :", projectDir)
	fmt.Println("Branch  :", gitops.CurrentBranch())
	fmt.Println("Remote  :", gitops.CurrentRemote())

	if cfg.Owner != "" && cfg.Repo != "" {
		fmt.Println("Repo    :", "https://github.com/"+cfg.Owner+"/"+cfg.Repo)
	}
	fmt.Println()
}

/* ============================================================
   Daily Git Operations
   ============================================================ */

func dailyMenu() {
	for {
		ui.Clear()
		ui.Header("Daily Git Operations")

		fmt.Println("1) Push changes (commit + push)")
		fmt.Println("2) Pull changes")
		fmt.Println("3) Smart Pull (auto-stash + pull)")
		fmt.Println("4) Fetch all remotes")
		fmt.Println("5) Git status")
		fmt.Println("6) Back")
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch ui.Input("Select option") {
		case "1":
			gitops.Push(ui.Input("Commit message"))
		case "2":
			gitops.Pull()
		case "3":
			gitops.SmartPull()
		case "4":
			gitops.Fetch()
		case "5":
			gitops.Status()
		case "6":
			return
		case "h", "help", "?":
			sectionHelp("Daily Git Operations", ui.HelpDaily)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}

/* ============================================================
   Branch / Remote
   ============================================================ */

func branchMenu() {
	for {
		ui.Clear()
		ui.Header("Branch / Remote")

		fmt.Println("1) Switch branch")
		fmt.Println("2) Switch remote")
		fmt.Println("3) Back")
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch ui.Input("Select option") {
		case "1":
			gitops.SwitchBranch()
		case "2":
			gitops.SwitchRemote()
		case "3":
			return
		case "h", "help", "?":
			sectionHelp("Branch / Remote", ui.HelpBranch)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}

/* ============================================================
   Stash & Undo
   ============================================================ */

func stashMenu() {
	for {
		ui.Clear()
		ui.Header("Stash & Undo")

		fmt.Println("1) Stash changes")
		fmt.Println("2) List stashes")
		fmt.Println("3) Apply last stash (pop)")
		fmt.Println("4) Undo last commit (keep changes)")
		fmt.Println("5) Back")
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch ui.Input("Select option") {
		case "1":
			gitops.StashSave()
		case "2":
			gitops.StashList()
		case "3":
			gitops.StashPop()
		case "4":
			gitops.UndoLastCommit()
		case "5":
			return
		case "h", "help", "?":
			sectionHelp("Stash & Undo", ui.HelpStash)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}

/* ============================================================
   Tools
   ============================================================ */

func toolsMenu() {
	for {
		ui.Clear()
		ui.Header("Tools")

		fmt.Println("1) Setup / Reconfigure")
		fmt.Println("2) Create / Link GitHub Repository")
		fmt.Println("3) Change Project Directory")
		fmt.Println("4) Doctor (health check)")
		fmt.Println("5) Back")
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch ui.Input("Select option") {
		case "1":
			setup.Run()
		case "2":
			setup.CreateOrLinkRepo()
		case "3":
			setup.ChangeProjectDir()
		case "4":
			doctor.Run()
		case "5":
			return
		case "h", "help", "?":
			sectionHelp("Tools", ui.HelpTools)
		default:
			ui.Error("Invalid option")
		}
		ui.Pause()
	}
}

/* ============================================================
   Help Screens
   ============================================================ */

func mainHelp() {
	ui.Clear()
	ui.Header("Help / About Git Genius")

	ui.PrintHelp(ui.HelpMain)
	ui.PrintHelp(ui.HelpGitHub)

	ui.Pause()
}

func sectionHelp(title string, help []string) {
	ui.Clear()
	ui.Header(title + " – Help")
	ui.PrintHelp(help)
	ui.Pause()
}
