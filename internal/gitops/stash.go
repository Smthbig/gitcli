package gitops

import (
	"strings"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
StashSave saves current working tree changes
*/
func StashSave() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	msg := ui.Input("Stash message (optional)")
	args := []string{"stash", "push"}

	if msg != "" {
		args = append(args, "-m", msg)
	}

	if err := system.RunGit(args...); err != nil {
		ui.Error("Failed to stash changes")
		return false
	}

	ui.Success("Changes stashed successfully")
	return true
}

/*
StashList shows all stashes and allows interaction
*/
func StashList() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	out, err := system.GitOutput("stash", "list")
	if err != nil {
		ui.Error("Failed to list stashes")
		return false
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 || lines[0] == "" {
		ui.Info("No stashes found")
		return true
	}

	options := append(lines, "Back")
	choice := ui.Select("Select stash to manage", options)

	if choice == len(options) {
		return true
	}

	selected := lines[choice-1]
	stashID := strings.Split(selected, ":")[0]

	ui.BoxMenu("Manage Stash: "+stashID, []string{
		"1) Apply (keep stash)",
		"2) Pop (apply and delete)",
		"3) Drop (delete stash)",
		"4) Back",
	})

	switch ui.MenuChoice() {
	case "1":
		if err := system.RunGit("stash", "apply", stashID); err != nil {
			ui.Error("Failed to apply stash")
		} else {
			ui.Success("Stash applied")
		}
	case "2":
		if err := system.RunGit("stash", "pop", stashID); err != nil {
			ui.Error("Failed to pop stash")
		} else {
			ui.Success("Stash popped")
		}
	case "3":
		if ui.Confirm("Are you sure you want to delete " + stashID + "?") {
			if err := system.RunGit("stash", "drop", stashID); err != nil {
				ui.Error("Failed to drop stash")
			} else {
				ui.Success("Stash dropped")
			}
		}
	}

	return true
}

/*
StashPop applies and removes latest stash
*/
func StashPop() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	if err := system.RunGit("stash", "pop"); err != nil {
		ui.Error("Failed to apply stash")
		return false
	}

	ui.Success("Stash applied successfully")
	return true
}
