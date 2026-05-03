package menu

import (
	"git-genius/internal/config"
	"git-genius/internal/gitops"
	"git-genius/internal/setup"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

/*
Start is the main entry point for the Git Genius menu system.
It handles initialization and the main loop.
*/
func Start(appVersion string, gitAvailable bool) {
	maybeOfferSetup(gitAvailable)

	for {
		ui.Clear()
		ui.BoxHeader(headerTitle(appVersion, gitAvailable))

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
