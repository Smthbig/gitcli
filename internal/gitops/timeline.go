package gitops

import (
	"fmt"
	"strings"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

// ShowActivityTimeline displays a simple ASCII representation of commit activity.
func ShowActivityTimeline() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ui.BoxHeader("7-Day Activity Timeline")

	// Get commit counts per day for the last 7 days
	// Using --since="7 days ago" and formatting by date
	out, err := system.GitOutput("log", "--since='7 days ago'", "--format=%ad", "--date=short")
	if err != nil {
		ui.Error("Failed to fetch activity data")
		return false
	}

	counts := make(map[string]int)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if line != "" {
			counts[line]++
		}
	}

	if len(lines) == 0 || lines[0] == "" {
		ui.Info("No activity in the last 7 days.")
		return true
	}

	ui.Info("Commits per day (ASCII Chart):")
	fmt.Println()
	
	// Note: In a real app we'd iterate over actual dates, 
	// but for this CLI we'll just show days with activity for simplicity.
	for date, count := range counts {
		bar := strings.Repeat("█", count)
		fmt.Printf(ui.Cyan+" %10s | "+ui.Yellow+"%-20s"+ui.Reset+" (%d)\n", date, bar, count)
	}
	fmt.Println()

	return true
}
