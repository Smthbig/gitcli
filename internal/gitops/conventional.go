package gitops

import (
	"fmt"

	"git-genius/internal/ui"
)

// ComposeConventionalMessage guides the user to create a structured commit message.
func ComposeConventionalMessage(defaultMsg string) string {
	if !ui.Confirm("Use Conventional Commit format?") {
		return ui.InputDefault("Commit message", defaultMsg)
	}

	prefix := ui.SelectConventionalType()
	scope := ui.Input("Scope (optional, e.g. ui, core)")
	subject := ui.Input("Short description")

	if subject == "" {
		subject = defaultMsg
	}

	if scope != "" {
		return fmt.Sprintf("%s(%s): %s", prefix, scope, subject)
	}
	return fmt.Sprintf("%s: %s", prefix, subject)
}
