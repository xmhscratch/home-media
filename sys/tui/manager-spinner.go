package tui

import (
	"fmt"

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

func (m *SpinnerModel) RenderView() string {
	return Styles.Main.Render(
		fmt.Sprintf("\n\n    %s    %s\n\n", m.ViewModel.View(), m.loadingText),
	)
}

func (m *SpinnerModel) UpdateSpinner(pipeData T_PipeData) tea.Cmd {
	m.loadingText = parseTextData(pipeData)
	return nil
}
