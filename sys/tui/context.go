package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type T_PipeData map[int]map[int]string

type TuiModel struct {
	tea.Model
	Output   any
	PipeData T_PipeData
	Render   func() (tea.Model, error)
}

type ListModel struct {
	*TuiModel
	Header string
	Output list.Model
}

type PipeModel struct {
	*TuiModel
	Output textinput.Model
}
