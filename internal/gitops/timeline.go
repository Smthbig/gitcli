package gitops

import (
	"fmt"
	"strings"
	"time"

	"git-genius/internal/system"
	"git-genius/internal/ui"
)

func GetRecentActivityData(days int) []int {
	out, err := system.GitOutput("log", fmt.Sprintf("--since='%d days ago'", days), "--format=%ad", "--date=short")
	if err != nil {
		return nil
	}

	counts := make(map[string]int)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	for _, line := range lines {
		if line != "" {
			counts[line]++
		}
	}

	// We want a sorted list for the last X days, including zeros
	res := make([]int, days)
	now := time.Now()
	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		res[days-1-i] = counts[date]
	}
	return res
}

// ShowActivityTimeline displays a simple ASCII representation of commit activity.
func ShowActivityTimeline() bool {
	if !system.EnsureGitRepo() {
		return false
	}

	ui.BoxHeader("7-Day Activity Timeline")
	data := GetRecentActivityData(7)
	if data == nil {
		ui.Error("Failed to fetch activity data")
		return false
	}

	ui.Info("Trend: " + ui.CyberSparkline(data))
	ui.Info("Commits per day (ASCII Chart):")
	fmt.Println()
	
	now := time.Now()
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, -i).Format("2006-01-02")
		count := data[7-1-i]
		if count > 0 {
			bar := strings.Repeat("█", count)
			fmt.Printf(ui.Cyan+" %10s | "+ui.Yellow+"%-20s"+ui.Reset+" (%d)\n", date, bar, count)
		}
	}
	fmt.Println()

	return true
}
