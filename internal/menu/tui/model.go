package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MenuItem struct {
	Label  string
	Action func() tea.Msg
	Key    string
}

type Model struct {
	Cursor       int
	MenuItems    []MenuItem
	GitAvailable bool
	Version      string
	Dashboard    []string
	Quitting     bool
	Width        int
	Height       int
	ActiveAction func() bool
	OnPush       func() bool
	OnUpdate     func() bool
	OnStatus     func() bool
}

func NewModel(version string, gitAvailable bool, items []MenuItem) Model {
	return Model{
		MenuItems:    items,
		Version:      version,
		GitAvailable: gitAvailable,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

type ActionMsg struct {
	Action func() bool
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ActionMsg:
		m.ActiveAction = msg.Action
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		case "down", "j":
			if m.Cursor < len(m.MenuItems)-1 {
				m.Cursor++
			}

		case "enter":
			if m.MenuItems[m.Cursor].Action != nil {
				action := m.MenuItems[m.Cursor].Action
				return m, func() tea.Msg { return action() }
			}

		// Quick keys - these trigger actions directly
		case "p":
			if m.OnPush != nil {
				return m, func() tea.Msg { return ActionMsg{Action: m.OnPush} }
			}
		case "u":
			if m.OnUpdate != nil {
				return m, func() tea.Msg { return ActionMsg{Action: m.OnUpdate} }
			}
		case "s":
			if m.OnStatus != nil {
				return m, func() tea.Msg { return ActionMsg{Action: m.OnStatus} }
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}
