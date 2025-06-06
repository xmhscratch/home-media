package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type ListItem struct {
	title, desc string
}

func (i ListItem) Title() string       { return i.title }
func (i ListItem) Description() string { return i.desc }
func (i ListItem) FilterValue() string { return i.title }

func (ctx *TuiManager) NewListModel() ListModel {
	m := ListModel{
		ViewModel: list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
	}
	m.UpdateList(ctx.PipeData)
	return m
}

func (m *ListModel) UpdateList(pipeData T_PipeData) tea.Cmd {
	var items []list.Item = parseListData(pipeData)
	return m.ViewModel.SetItems(items)
}

func parseListData(pipeData T_PipeData) (data []list.Item) {
	var unsorted = map[int]list.Item{}
	for i, v := range pipeData {
		if _, ok := v[0]; !ok {
			v[0] = "(empty)"
		}
		if _, ok := v[1]; !ok {
			v[1] = "(empty)"
		}
		var (
			title string = v[0]
			desc  string = ""
		)
		for j := range len(v) - 1 {
			desc += strings.TrimSpace(v[j+1])
		}
		// log.Println(desc)
		unsorted[i] = ListItem{title, desc}
	}
	data = []list.Item{}
	for i := range len(unsorted) {
		data = append(data, unsorted[i])
	}
	return data
}
