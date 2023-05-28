package app

import tea "github.com/charmbracelet/bubbletea"

type fetchreposmodel struct {
}

func (m *fetchreposmodel) Init() tea.Cmd {
	return nil
}

func (m *fetchreposmodel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *fetchreposmodel) View() string {
	return ""
}
