package tui

import tea "github.com/charmbracelet/bubbletea"

func (m *GlamourModel) UpdateError(PipeData T_PipeData) tea.Cmd {
	_ = m.SetGlamourContent(parseTextData(PipeData))
	return nil
}
