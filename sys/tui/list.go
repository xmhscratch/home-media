package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

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
