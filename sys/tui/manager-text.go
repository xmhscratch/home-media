package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (ctx *TuiManager) NewTextModel() TextModel {
	m := TextModel{}
	m.ViewModel = textinput.New()
	m.UpdateText(ctx.PipeData)

	return m
}

func (m *TextModel) UpdateText(pipeData T_PipeData) tea.Cmd {
	m.ViewModel.Prompt = ""
	m.ViewModel.Cursor.Style = Styles.Cursor
	m.ViewModel.Width = 48
	m.ViewModel.CursorStart()
	m.ViewModel.SetValue("")
	m.ViewModel.CursorEnd()
	m.ViewModel.Focus()

	m.current = 0
	m.historyText += m.rawText
	m.rawText = parseTextData(pipeData)

	return nil
}

func (m *TextModel) TickCmd() tea.Cmd {
	m.current += 1
	if m.current > len(m.rawText) {
		m.current = len(m.rawText)
		return textinput.Blink
	}
	return tickCmd(REFRESH_RATE)
}

func (m *TextModel) RenderView() string {
	var (
		text        string
		current     int    = m.current
		historyText string = m.historyText
		rawText     string = m.rawText
	)
	if current <= 1 {
		text = ""
	} else {
		text = rawText[0:current]
	}
	m.ViewModel.SetCursor(current)
	return Styles.Main.Render(fmt.Sprintf("%s%s%s", historyText, text, m.ViewModel.View()))
}

func parseTextData(pipeData T_PipeData) string {
	var sb strings.Builder
	for i := range len(pipeData) {
		line := pipeData[i]
		for j := range len(line) {
			col := line[j]
			sb.WriteString(fmt.Sprintf("%s\t", col))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
