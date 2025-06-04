package tui

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type T_PipeData map[int]map[int]string
type T_OutputMode int

const (
	_                   T_OutputMode = iota
	OUTPUT_VIEW_SPINNER T_OutputMode = 1
	OUTPUT_VIEW_TEXT    T_OutputMode = 2
	OUTPUT_VIEW_ERROR   T_OutputMode = 3
	OUTPUT_VIEW_LIST    T_OutputMode = 4
)

const (
	progressBarWidth  = 71
	progressFullChar  = "█"
	progressEmptyChar = "░"
	dotChar           = " • "
)

type DefinedStyle struct {
	Main          lipgloss.Style
	Keyword       lipgloss.Style
	Subtle        lipgloss.Style
	Ticks         lipgloss.Style
	Checkbox      lipgloss.Style
	Cursor        lipgloss.Style
	ProgressEmpty string
	Dot           string
	RampColor     []lipgloss.Style
}

type TuiManager struct {
	Program           *tea.Program
	Header            string
	CurrentOutputMode T_OutputMode
	PipeData          T_PipeData
	PipeErr           string
	Output            struct {
		Error   GlamourModel
		Spinner SpinnerModel
		List    ListModel
		Text    TextModel
	}
	RefreshRate     int
	SpinnerTick     tea.Cmd
	CursorBlinkTick tea.Cmd
}

type SpinnerModel struct {
	ViewModel   spinner.Model
	loadingText string
}

type TextModel struct {
	ViewModel textinput.Model
	rawText   string
	current   int
}

type ListModel struct {
	ViewModel list.Model
}

type GlamourModel struct {
	ViewModel viewport.Model
	renderer  *glamour.TermRenderer
}

type T_SocketResponse struct {
	msg string
	err error
}

type pipeResMsg struct {
	T_OutputMode
	string
}

type tickMsg struct{ time.Time }

var Styles *DefinedStyle = &DefinedStyle{
	Main:          lipgloss.NewStyle().Margin(1, 2),
	Keyword:       lipgloss.NewStyle().Foreground(lipgloss.Color("211")),
	Subtle:        lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	Ticks:         lipgloss.NewStyle().Foreground(lipgloss.Color("79")),
	Checkbox:      lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
	Cursor:        lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
	ProgressEmpty: lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(progressEmptyChar),
	Dot:           lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(dotChar),
	RampColor:     makeRampStyles("#B14FFF", "#00FFA3", progressBarWidth),
}
