package menu

import (
	"fmt"
	"path/filepath"
	"strings"

	"git-genius/internal/config"
	"git-genius/internal/doctor"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func Start(appVersion string, gitAvailable bool) {
	maybeOfferSetup(gitAvailable)

	for {
		ui.Clear()
		ui.Header(headerTitle(appVersion, gitAvailable))

		showContext(gitAvailable)

		if gitAvailable {
			if !mainMenu() {
				return
			}
			continue
		}

		if !limitedMenu() {
			return
		}
	}
}

func headerTitle(appVersion string, gitAvailable bool) string {
	version := strings.TrimSpace(appVersion)
	if version == "" {
		version = "dev"
	}

	title := "Git Genius v" + version
	if !gitAvailable {
		title += " (Limited Mode)"
	}
	return title
}

func maybeOfferSetup(gitAvailable bool) {
	if !gitAvailable {
		return
	}

	cfg := config.Load()
	if config.HasProjectConfig(cfg.GetWorkDir()) || config.HasHistoryForWorkDir(cfg.GetWorkDir()) {
		return
	}

	if system.IsGitRepo() {
		if gitops.InspectRepoState().HasCommits {
			ui.Info("First run for this repository: Git Genius has not been configured here yet")
		} else {
			ui.Info("Brand-new repository detected: no commits yet and no Git Genius setup found")
		}
	} else {
		ui.Info("Brand-new project directory detected: setup can initialize Git and prepare the first push")
	}
	if setup.Run() {
		ui.Pause()
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

func showContext(gitAvailable bool) {
	cfg := config.Load()
	projectDir := cfg.GetWorkDir()

	fmt.Println("Project :", filepath.Base(projectDir))
	fmt.Println("Path    :", projectDir)

	if gitAvailable {
		state := gitops.InspectRepoState()
		fmt.Println("Branch  :", state.Branch)
		fmt.Println("Remote  :", system.CurrentRemote())
		if state.HasAheadBehind {
			fmt.Println("Sync    :", state.AheadBehindSummary())
		}
		if state.FirstRun {
			fmt.Println("First Run:", "Setup recommended for this repo")
		}
		if !state.HasCommits {
			fmt.Println("Repo    :", "No commits yet")
		} else if state.NeedsFirstPush {
			fmt.Println("Publish :", "First push still pending")
		}
	} else {
		fmt.Println("Mode    :", "Limited (Git unavailable)")
	}

	if cfg.Owner != "" && cfg.Repo != "" {
		fmt.Println("Repo    :", "https://github.com/"+cfg.Owner+"/"+cfg.Repo)
	}
	for _, s := range config.HistorySuggestions(projectDir) {
		fmt.Println("Suggest :", s)
	}
	fmt.Println()
}

func mainMenu() bool {
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
		toolsMenu(true)
	case "5", "h", "help", "?":
		mainHelp()
	case "6":
		ui.Info("Goodbye")
		return false
	default:
		ui.Error("Invalid option")
		ui.Pause()
	}

	return true
}

func limitedMenu() bool {
	fmt.Println("1) Setup / Reconfigure")
	fmt.Println("2) Switch Project")
	fmt.Println("3) Doctor (health check)")
	fmt.Println("4) Help / About")
	fmt.Println("5) Exit")
	fmt.Println()
	fmt.Println("Tip: install Git to unlock daily operations")

	switch normalizeChoice(ui.Input("Select option")) {
	case "1":
		track("tools", "setup_reconfigure", setup.Run)
	case "2":
		track("tools", "switch_project", setup.SwitchProject)
	case "3":
		track("tools", "doctor", doctor.Run)
	case "4", "h", "help", "?":
		mainHelp()
	case "5":
		ui.Info("Goodbye")
		return false
	default:
		ui.Error("Invalid option")
		ui.Pause()
	}

	return true
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

		fmt.Println("1) Switch to existing branch")
		fmt.Println("2) Create new branch")
		fmt.Println("3) Configure remote")
		fmt.Println("4) Back")
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch normalizeChoice(ui.Input("Select option")) {
		case "1":
			track("branch", "switch_branch", gitops.SwitchBranch)
		case "2":
			track("branch", "create_branch", gitops.CreateBranch)
		case "3":
			track("branch", "switch_remote", gitops.SwitchRemote)
		case "4":
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

func toolsMenu(gitAvailable bool) {
	for {
		ui.Clear()
		ui.Header("Tools")

		fmt.Println("1) Setup / Reconfigure")
		fmt.Println("2) Switch Project")

		if gitAvailable {
			fmt.Println("3) Create / Link GitHub Repository")
			fmt.Println("4) Git Auth / Credential Helper")
			fmt.Println("5) Doctor (health check)")
			fmt.Println("6) Back")
		} else {
			fmt.Println("3) Doctor (health check)")
			fmt.Println("4) Back")
		}
		fmt.Println()
		fmt.Println("Tip: h = help")

		switch normalizeChoice(ui.Input("Select option")) {
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
				track("tools", "configure_git_auth", setup.ConfigureGitAuth)
			} else {
				return
			}
		case "5":
			if gitAvailable {
				track("tools", "doctor", doctor.Run)
			} else {
				ui.Error("Invalid option")
			}
		case "6":
			if gitAvailable {
				return
			}
			ui.Error("Invalid option")
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
	ui.PrintHelp(ui.HelpWorkflow)
	ui.PrintHelp(ui.HelpGitHub)
	ui.PrintHelp(ui.HelpTroubleshooting)

	ui.Pause()
}

func sectionHelp(title string, help []string) {
	ui.Clear()
	ui.Header(title + " – Help")
	ui.PrintHelp(help)
	ui.Pause()
}
