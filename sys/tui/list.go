package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.Output.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.Output, cmd = m.Output.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return docStyle.Render(m.Output.View())
}

func (m *ListModel) Render() error {
	var (
		data []list.Item
		err  error
	)

	if data, err = m.MarshalData(); err != nil {
		return err
	}

	md := ListModel{Output: list.New(data, list.NewDefaultDelegate(), 0, 0)}
	md.Output.Title = m.Header

	p := tea.NewProgram(md, tea.WithAltScreen())
	_, err = p.Run()

	return err
}

func (m *ListModel) MarshalData() (data []list.Item, err error) {
	var unsorted = map[int]list.Item{}
	for i, v := range m.TuiModel.PipeData {
		if _, ok := v[0]; !ok {
			v[0] = "(empty)"
		}
		if _, ok := v[1]; !ok {
			v[1] = "(empty)"
		}
		unsorted[i] = item{title: v[0], desc: v[1]}
	}
	data = []list.Item{}
	for i := range len(unsorted) {
		data = append(data, unsorted[i])
	}
	return data, err
}
