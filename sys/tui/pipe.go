package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m *PipeModel) Render() error {
	ppStr, err := ReadPipe()
	if err != nil {
		return err
	}
	model := NewModel(ppStr)

	if _, err := tea.NewProgram(model).Run(); err != nil {
		return fmt.Errorf("couldn't start program: %s", err)
	}
	return nil
}

func (m PipeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m PipeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyCtrlC, tea.KeyEscape, tea.KeyEnter:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.Output, cmd = m.Output.Update(msg)
	return m, cmd
}

func (m PipeModel) View() string {
	return fmt.Sprintf(
		"\nYou piped in: %s\n\nPress ^C to exit",
		m.Output.View(),
	)
}

func NewModel(initialValue string) (m PipeModel) {
	i := textinput.New()
	i.Prompt = ""
	i.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	i.Width = 48
	i.SetValue(initialValue)
	i.CursorEnd()
	i.Focus()

	m.Output = i
	return
}
