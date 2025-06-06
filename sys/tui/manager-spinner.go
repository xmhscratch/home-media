package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (ctx *TuiManager) NewSpinnerModel() SpinnerModel {
	return SpinnerModel{
		ViewModel: spinner.New(
			spinner.WithSpinner(spinner.Meter),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
	}
}

func (m *SpinnerModel) UpdateSpinner(pipeData T_PipeData) tea.Cmd {
	m.loadingText = parseTextData(pipeData)
	return nil
}
