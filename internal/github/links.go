package github

import (
	"fmt"

	"git-genius/internal/config"
	"git-genius/internal/ui"
)

// OpenRepo opens the GitHub repository in the browser.
func OpenRepo() bool {
	cfg := config.Load()
	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("GitHub repository not configured")
		return false
	}

	url := fmt.Sprintf("https://github.com/%s/%s", cfg.Owner, cfg.Repo)
	return openURL(url)
}

// OpenPRs opens the Pull Requests page.
func OpenPRs() bool {
	cfg := config.Load()
	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("GitHub repository not configured")
		return false
	}

	url := fmt.Sprintf("https://github.com/%s/%s/pulls", cfg.Owner, cfg.Repo)
	return openURL(url)
}

// OpenIssues opens the Issues page.
func OpenIssues() bool {
	cfg := config.Load()
	if cfg.Owner == "" || cfg.Repo == "" {
		ui.Error("GitHub repository not configured")
		return false
	}

	url := fmt.Sprintf("https://github.com/%s/%s/issues", cfg.Owner, cfg.Repo)
	return openURL(url)
}

func openURL(url string) bool {
	ui.Info("Opening: " + url)
	if err := ui.OpenURL(url); err != nil {
		ui.Error("Failed to open URL: " + err.Error())
		return false
	}
	return true
}
