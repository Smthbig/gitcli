package tui

import (
	"fmt"
	"strings"

	"git-genius/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	// Header
	header := ui.HeaderStyle.Render(fmt.Sprintf("⚡ GIT GENIUS v%s ⚡", m.Version))

	// Menu
	var menuBuilder strings.Builder
	for i, item := range m.MenuItems {
		if i == m.Cursor {
			menuBuilder.WriteString(ui.SelectedStyle.Render("❯ " + item.Label))
		} else {
			menuBuilder.WriteString("  " + item.Label)
		}
		menuBuilder.WriteString("\n")
	}
	
	menuBox := ui.MainBorderStyle.
		Width(30).
		Render(ui.TitleStyle.Render(" MENU ") + "\n\n" + menuBuilder.String())

	// Dashboard / Status
	var statusBuilder strings.Builder
	for _, line := range m.Dashboard {
		statusBuilder.WriteString(line + "\n")
	}

	statusBox := ui.PanelStyle.
		Width(m.Width - 40). // Adaptive width
		Render(ui.TitleStyle.Render(" STATUS ") + "\n\n" + statusBuilder.String())

	// Layout side-by-side or stacked
	var mainContent string
	if m.Width > 75 {
		mainContent = lipgloss.JoinHorizontal(lipgloss.Top, menuBox, "   ", statusBox)
	} else {
		mainContent = lipgloss.JoinVertical(lipgloss.Left, menuBox, "\n", statusBox)
	}

	// Footer / Help
	footer := ui.InfoStyle.Render("\n [↑/↓/j/k] Navigate • [Enter] Select • [p]ush • [u]pdate • [q]uit")

	return ui.DocStyle.Render(lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer))
}
