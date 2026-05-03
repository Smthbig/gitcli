package setup

import (
	"git-genius/internal/config"
	"git-genius/internal/github"
	"git-genius/internal/system"
	"git-genius/internal/ui"
)

const (
	credentialHelperCache = "cache --timeout=43200"
	credentialHelperStore = "store"
)

func configureGitAuthIfNeeded() bool {
	helper, _ := system.GitCredentialHelper()
	if helper == "" {
		ui.Warn("Git credential helper is not configured")
		ui.Info("Without one, HTTPS pushes may ask for username and token each time")
		if !ui.ConfirmDefault("Configure Git authentication helper now?", true) {
			return true
		}
		return ConfigureGitAuth()
	}

	ui.Success("Git credential helper: " + helper)

	if github.GetToken() == "" {
		return true
	}

	if ui.Confirm("Refresh the Git credential helper with the current GitHub token?") {
		return preloadGitHubCredential()
	}

	return true
}

func ConfigureGitAuth() bool {
	ui.Clear()
	ui.BoxHeader("Git Authentication")

	if err := system.EnsureGitInstalled(); err != nil {
		ui.Error("Git is required to configure authentication")
		return false
	}

	currentHelper, _ := system.GitCredentialHelper()
	if currentHelper == "" {
		ui.Warn("No global Git credential helper configured")
	} else {
		ui.Success("Current Git credential helper: " + currentHelper)
	}

	ui.Info("1) Memory cache (recommended): keeps credentials for 12 hours")
	ui.Info("2) Persistent store: saves credentials to ~/.git-credentials")
	ui.Info("3) Clear helper: remove current credential helper")
	ui.Info("4) Keep current settings")

	choice := ui.Select("Select authentication helper mode", []string{
		"Memory cache (recommended)",
		"Persistent plain-text store",
		"Clear current helper",
		"Keep current settings",
	})

	switch choice {
	case 1:
		if err := system.SetGitCredentialHelper(credentialHelperCache); err != nil {
			ui.Error("Failed to configure memory cache helper")
			ui.Info(err.Error())
			return false
		}
		ui.Success("Configured Git credential helper: memory cache (12 hours)")
	case 2:
		if err := system.SetGitCredentialHelper(credentialHelperStore); err != nil {
			ui.Error("Failed to configure persistent store helper")
			ui.Info(err.Error())
			return false
		}
		ui.Warn("Git will store credentials in plain text at ~/.git-credentials")
		ui.Success("Configured Git credential helper: persistent store")
	case 3:
		if err := system.ClearGitCredentialHelper(); err != nil {
			ui.Error("Failed to clear Git credential helper")
			ui.Info(err.Error())
			return false
		}
		ui.Success("Git credential helper cleared")
		return true
	case 4:
		ui.Info("Keeping current Git authentication settings")
		if currentHelper == "" {
			return true
		}
	default:
		ui.Error("Invalid option")
		return false
	}

	return preloadGitHubCredential()
}

func preloadGitHubCredential() bool {
	token := github.GetToken()
	if token == "" {
		ui.Warn("No GitHub token available to preload into Git")
		ui.Info("Set a token in Setup or export " + github.EnvTokenName)
		return true
	}

	user := github.GetUsername()
	if user == "" {
		client, err := github.NewClient()
		if err == nil {
			if resolved, resolveErr := client.GetAuthenticatedUser(); resolveErr == nil {
				user = resolved
				_ = github.SaveAuth(token, user)
			}
		}
	}

	if user == "" {
		user = config.Load().Owner
	}

	if user == "" {
		ui.Warn("Could not resolve a GitHub username for credential preload")
		ui.Info("Re-run Setup with a validated token or export " + github.EnvUserName)
		return true
	}

	if !ui.ConfirmDefault("Preload GitHub credentials into the helper now?", true) {
		return true
	}

	if err := system.ApproveGitHubCredential(user, token); err != nil {
		ui.Warn("Could not preload GitHub credentials")
		ui.Info(err.Error())
		return false
	}

	ui.Success("GitHub credentials saved to the Git credential helper")
	ui.Info("Future HTTPS pushes should stop asking for credentials until the helper expires or is cleared")
	return true
}
