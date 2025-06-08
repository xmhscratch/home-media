package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func (ctx *TuiManager) NewTextModel() TextModel {
	return newTextModel()
}

func newTextModel() TextModel {
	m := TextModel{}
	m.ViewModel = textinput.New()
	m.ViewModel.Prompt = ""
	m.ViewModel.Cursor.Style = Styles.Cursor
	m.ViewModel.Width = 48
	m.ViewModel.CursorStart()
	m.ViewModel.SetValue("")
	m.ViewModel.CursorEnd()
	m.ViewModel.Focus()

	var _cursor int = 0
	m._stack = []string{}
	m._cursor = &(_cursor)
	m.historyText = ""
	m.rawText = ""

	return m
}

func (m *TextModel) Reset() tea.Cmd {
	*m = newTextModel()
	return tea.ClearScreen
}

func (m *TextModel) UpdateText(pipeData T_PipeData) tea.Cmd {
	m._stack = append(m._stack, parseTextData(pipeData))
	return nil
}

func (m *TextModel) TickCmd() tea.Cmd {
	if len(m._stack) == 1 && *m._cursor == 0 {
		m.rawText = m._stack[0]
		m._stack = []string{}
	}
	if len(m._stack) > 1 && *m._cursor == 0 {
		m.rawText = m._stack[0]
		m._stack = m._stack[1:]
	}
	if len(m._stack) == 0 && *m._cursor >= len(m.rawText)-1 {
		return textinput.Blink
	}
	m.historyText += m.rawText[*m._cursor : *m._cursor+1]
	*m._cursor += 1
	if *m._cursor >= len(m.rawText) {
		*m._cursor = 0
	}
	m.ViewModel.SetCursor(len(m.historyText))
	return tickCmd(REFRESH_RATE)
}

func (m *TextModel) RenderView() string {
	return Styles.Main.Render(fmt.Sprintf("%s%s", m.historyText, m.ViewModel.View()))
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
