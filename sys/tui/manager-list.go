package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type ListItem struct {
	title, desc string
	value       string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }
func (i ListItem) GetValue() string    { return i.value }

func (ctx *TuiManager) NewListModel() ListModel {
	m := ListModel{
		ViewModel: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
	m.UpdateList(ctx.PipeData)
	return m
}

func (m *ListModel) UpdateList(pipeData T_PipeData) tea.Cmd {
	var items map[int]ListItem = parseListData(pipeData)
	data := []list.Item{}
	for i := range len(items) {
		data = append(data, items[i])
	}
	m.Items = items
	return m.ViewModel.SetItems(data)
}

func (m *ListModel) RenderView() string {
	return Styles.Main.Render(m.ViewModel.View())
}

func (m *ListModel) BindExtraKeyCommands(mgr TuiManager, msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "enter" {
		// selIndex := m.ViewModel.GlobalIndex()
		// selItem := m.Items[selIndex]
		// m.Program.Send(pipeResMsg{OUTPUT_VIEW_TEXT, "", selItem.GetValue()})
		// m.CommandExec selItem.GetValue()
	}
	return nil
}

func parseListData(pipeData T_PipeData) (data map[int]ListItem) {
	data = map[int]ListItem{}

	for i, v := range pipeData {
		var (
			title string = "(empty)"
			desc  string = "(empty)"
			value string = ""
		)
		if _, ok := v[0]; ok {
			title = v[0]
		}
		if _, ok := v[1]; ok {
			desc = strings.TrimSpace(v[1])
		}
		vs := lo.FilterMapToSlice(v, func(ik int, iv string) (string, bool) {
			if ik <= 1 {
				return "", false
			}
			return iv, true
		})
		value = strings.Join(vs, " ")
		// log.Println(desc)
		data[i] = ListItem{title, desc, value}
	}

	return data
}
