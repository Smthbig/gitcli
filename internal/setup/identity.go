package setup

import (
	"git-genius/internal/system"
	"git-genius/internal/ui"
	"strings"
)

func ensureGitIdentity(dir string) bool {
	name, _ := system.GitOutputAt(dir, "config", "--get", "user.name")
	email, _ := system.GitOutputAt(dir, "config", "--get", "user.email")

	if name != "" && email != "" {
		ui.Success("Git identity confirmed")
		ui.Info("Name : " + name)
		ui.Info("Email: " + email)
		return true
	}

	ui.Warn("Git identity not configured in this directory")

	if name == "" {
		val := strings.TrimSpace(ui.Input("Enter your name"))
		if val != "" {
			if err := system.RunGitAt(dir, "config", "user.name", val); err != nil {
				ui.Error("Failed to set user.name")
				return false
			}
			name = val
		}
	}

	if email == "" {
		val := strings.TrimSpace(ui.Input("Enter your email"))
		if val != "" {
			if err := system.RunGitAt(dir, "config", "user.email", val); err != nil {
				ui.Error("Failed to set user.email")
				return false
			}
			email = val
		}
	}

	if name != "" && email != "" {
		ui.Success("Git identity configured")
		return true
	}

	ui.Error("Git identity is required for commits")
	return false
}
