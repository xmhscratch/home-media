package tui

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const GLAMOUR_WIDTH = 100
const GLAMOUR_GUTTER = 20

func (ctx *TuiManager) NewGlamourModel(pipeData T_PipeData) (GlamourModel, error) {
	var err error
	m := GlamourModel{
		ViewModel: viewport.New(GLAMOUR_WIDTH, 20),
	}
	m.ViewModel.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("211")).
		PaddingRight(20)

	glamourRenderWidth := GLAMOUR_WIDTH - m.ViewModel.Style.GetHorizontalFrameSize() - GLAMOUR_GUTTER

	m.renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return m, err
	}

	m.SetGlamourContent(parseTextData(pipeData))
	return m, nil
}

func (m *GlamourModel) RenderView() string {
	return ""
}

func (m *GlamourModel) SetGlamourContent(input string) error {
	str, err := m.renderer.Render(input)
	if err != nil {
		return err
	}
	m.ViewModel.SetContent(str)
	return nil
}
