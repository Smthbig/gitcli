package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#00FFFF") // Cyan
	AccentColor    = lipgloss.Color("#FF00FF") // Magenta
	SuccessColor   = lipgloss.Color("#00FF00") // Green
	WarningColor   = lipgloss.Color("#FFFF00") // Yellow
	ErrorColor     = lipgloss.Color("#FF0000") // Red
	HighlightColor = lipgloss.Color("#7D56F4") // Purple/Blue for selection
	DimColor       = lipgloss.Color("#777777") // Gray

	// Styles
	DimStyle = lipgloss.NewStyle().Foreground(DimColor)
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(AccentColor).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(AccentColor).
			Padding(0, 1).
			MarginBottom(1).
			Align(lipgloss.Center)

	MainBorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1)

	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(HighlightColor).
			Bold(true)

	TitleStyle = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true).
			Underline(true)

	InfoStyle    = lipgloss.NewStyle().Foreground(PrimaryColor)
	SuccessStyle = lipgloss.NewStyle().Foreground(SuccessColor)
	WarningStyle = lipgloss.NewStyle().Foreground(WarningColor)
	ErrorStyle   = lipgloss.NewStyle().Foreground(ErrorColor)
	DividerStyle = lipgloss.NewStyle().Foreground(AccentColor)

	DocStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)
)
