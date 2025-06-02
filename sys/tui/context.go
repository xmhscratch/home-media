package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type T_PipeData map[int]map[int]string

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

var (
	keywordStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	subtleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	ticksStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("79"))
	checkboxStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
	progressEmpty = subtleStyle.Render(progressEmptyChar)
	dotStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar)
	mainStyle     = lipgloss.NewStyle().MarginLeft(2)
	ramp          = makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth)
)

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
