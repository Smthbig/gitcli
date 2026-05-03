package gitops

import (
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// ShowGraph displays a color-coded git log graph.
func ShowGraph() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ui.BoxHeader("Visual Git History")
	
	// --graph: draw the tree
	// --oneline: concise format
	// --decorate: show branch/tag names
	// --color: enable ANSI colors
	err := system.RunGit("log", "--graph", "--oneline", "--decorate", "--color", "-n", "15")
	if err != nil {
		ui.Error("Failed to render git graph")
		return false
	}
	
	return true
}
