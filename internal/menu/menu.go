package menu

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/system"
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

		switch normalizeChoice(ui.Input("Select option")) {
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

func normalizeChoice(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func track(section, action string, fn func() bool) {
	cfg := config.Load()
	ok := false
	note := ""
	defer func() {
		if r := recover(); r != nil {
			ok = false
			note = fmt.Sprintf("panic: %v", r)
			ui.Error("Unexpected error occurred")
		}
		config.RecordHistory(cfg.GetWorkDir(), section, action, ok, note)
	}()
	ok = fn()
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
	fmt.Println("Remote  :", system.CurrentRemote())

	if cfg.Owner != "" && cfg.Repo != "" {
		fmt.Println("Repo    :", "https://github.com/"+cfg.Owner+"/"+cfg.Repo)
	}
	for _, s := range config.HistorySuggestions(projectDir) {
		fmt.Println("Suggest :", s)
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

		switch normalizeChoice(ui.Input("Select option")) {
		case "1":
			track("daily", "push", func() bool { return gitops.Push(ui.Input("Commit message")) })
		case "2":
			track("daily", "pull", gitops.Pull)
		case "3":
			track("daily", "smart_pull", gitops.SmartPull)
		case "4":
			track("daily", "fetch", gitops.Fetch)
		case "5":
			track("daily", "status", gitops.Status)
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

		switch normalizeChoice(ui.Input("Select option")) {
		case "1":
			track("branch", "switch_branch", gitops.SwitchBranch)
		case "2":
			track("branch", "switch_remote", gitops.SwitchRemote)
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

		switch normalizeChoice(ui.Input("Select option")) {
		case "1":
			track("stash", "stash_save", gitops.StashSave)
		case "2":
			track("stash", "stash_list", gitops.StashList)
		case "3":
			track("stash", "stash_pop", gitops.StashPop)
		case "4":
			track("stash", "undo_last_commit", gitops.UndoLastCommit)
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

		switch normalizeChoice(ui.Input("Select option")) {
		case "1":
			track("tools", "setup_reconfigure", func() bool { setup.Run(); return true })
		case "2":
			track("tools", "create_or_link_repo", func() bool { setup.CreateOrLinkRepo(); return true })
		case "3":
			track("tools", "change_project_dir", func() bool { setup.ChangeProjectDir(); return true })
		case "4":
			track("tools", "doctor", func() bool { doctor.Run(); return true })
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
