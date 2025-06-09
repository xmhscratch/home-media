package tui

import (
	"crypto/md5"
	"fmt"
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
	return newListModel()
}

func newListModel() ListModel {
	m := ListModel{}
	m.ViewModel = list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	return m
}

func (m *ListModel) Reset() tea.Cmd {
	*m = newListModel()
	return nil
}

func (m *ListModel) UpdateList(pipeData T_PipeData) tea.Cmd {
	var (
		uid   string
		items map[int]ListItem
	)

	items, uid = parseListData(pipeData)
	// fmt.Printf("%v\n", items)
	// if m.uid != uid {
	// 	return m.Reset()
	// }
	m.uid = uid
	m.Items = items

	var cmds []tea.Cmd
	cmds = append(cmds, m.TickCmd())
	for index := range len(items) {
		cmds = append(cmds, m.ViewModel.InsertItem(index, list.Item(items[index])))
	}

	return tea.Sequence(cmds...)
}

func (m *ListModel) TickCmd() tea.Cmd {
	// fmt.Printf("%v\n", m.Items)
	return nil
}

func (m *ListModel) SetSize(w int, h int) {
	m.ViewModel.SetSize(w, h)
}

func (m *ListModel) RenderView() string {
	// return Styles.Main.Render("")
	return Styles.Main.Render(m.ViewModel.View())
}

func (m *ListModel) BindExtraKeyCommands(mgr TuiManager, msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "enter" {
		selIndex := m.ViewModel.GlobalIndex()
		selItem := m.Items[selIndex]
		go func() {
			mgr.Program.Send(execProgram{
				m.CommandExec,
				selItem.GetValue(),
			})
			// mgr.Program.Send(pipeResMsg{OUTPUT_VIEW_TEXT, "", fmt.Sprintf("%s\n%s\n%s", selItem.Title(), selItem.Description(), selItem.GetValue())})
		}()
	}
	return nil
}

func parseListData(pipeData T_PipeData) (data map[int]ListItem, uid string) {
	data = map[int]ListItem{}
	ids := []string{}

	for i := range len(pipeData) {
		var (
			title string = "(empty)"
			desc  string = "(empty)"
			bArgs []byte = make([]byte, 0)
		)
		if _, ok := pipeData[i][0]; ok {
			title = pipeData[i][0]
		}
		if _, ok := pipeData[i][1]; ok {
			desc = strings.TrimSpace(pipeData[i][1])
		}
		// log.Println(desc)

		args := lo.FilterMapToSlice(pipeData[i], func(ik int, iv string) (string, bool) {
			if ik <= 1 {
				return "", false
			}
			return iv, true
		})
		for i := range len(args) {
			fmt.Appendf(bArgs, " %s", args[i])
		}
		value := string(bArgs)

		data[i] = ListItem{title, desc, value}
		ids = append(ids, title)
	}

	return data, hashSlice(ids)
}

func hashSlice(s []string) string {
	h := md5.New()
	for _, v := range s {
		h.Write(fmt.Appendf([]byte{}, "%s", v))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
