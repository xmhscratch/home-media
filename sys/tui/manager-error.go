package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *GlamourModel) UpdateError(PipeData T_PipeData) tea.Cmd {
	_ = m.SetGlamourContent(parseTextData(PipeData))
	return nil
}

func (m *GlamourModel) RenderErrorView() string {
	return Styles.Main.Render(
		m.ViewModel.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\n  ↑/↓: Navigate • q: Quit\n"),
	)
}
