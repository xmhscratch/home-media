package tui

import (
	"home-media/sys"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type T_PipeData (map[int]map[int]string)
type T_OutputMode (int)

const ASCII_RS = 0x1E
const RGXP_MESSAGE_PAYLOAD = `^((\d+(?=\x1E))((?=\x1E)..[^\x1E\n]*|)((?=\x1E).*)|.*)$`
const RGXP_INSTALL_PKGINFO = `\(([\d]+)\/([\d]+)\)[\ ]*Installing[\ ]*([\w\S]+)[\ ]*\([a-z\d.-]+\)`
const RGXP_TRIM_EXTRA_VARS = `\%\!\(EXTRA[\s]{0,}([\s\S]+(?:\=[\s\S]+[\,\s]{0,})+)*\)$`
const RGXP_TRIM_MISSING_VARS = ``
const UNIX_VW_SOCKET_PATH = "/run/tuidw.sock"
const UNIX_EX_SOCKET_PATH = "/run/tuidx.sock"
const (
	_                     T_OutputMode = iota
	OUTPUT_SOCKET         T_OutputMode = 1
	OUTPUT_VIEW_SPINNER   T_OutputMode = 2
	OUTPUT_VIEW_TEXT      T_OutputMode = 3
	OUTPUT_VIEW_ERROR     T_OutputMode = 4
	OUTPUT_VIEW_LIST      T_OutputMode = 5
	OUTPUT_VIEW_INSTALLER T_OutputMode = 6
)
const REFRESH_RATE int = 5
const REFRESH_RATE_IN_SECONDS int = 1000 / 5

type DefinedStyle struct {
	Main           lipgloss.Style
	Keyword        lipgloss.Style
	Subtle         lipgloss.Style
	Ticks          lipgloss.Style
	Checkbox       lipgloss.Style
	Cursor         lipgloss.Style
	ProgressEmpty  string
	Dot            string
	RampColor      []lipgloss.Style
	CurrentPkgName lipgloss.Style
	Done           lipgloss.Style
	CheckMark      lipgloss.Style
}

type TuiManager struct {
	Config            *sys.Config
	Program           *tea.Program
	Header            string
	CurrentOutputMode T_OutputMode
	PipeData          T_PipeData
	Output            struct {
		Error     GlamourModel
		Spinner   SpinnerModel
		List      ListModel
		Text      TextModel
		Installer InstallerModel
	}
}

type SpinnerModel struct {
	ViewModel   spinner.Model
	loadingText string
}

type TextModel struct {
	ViewModel   textinput.Model
	historyText string
	rawText     string
	_stack      []string
	_cursor     *int
}

type ListModel struct {
	ViewModel   list.Model
	Items       map[int]ListItem
	CommandExec string
	uid         string
}

type GlamourModel struct {
	ViewModel viewport.Model
	renderer  *glamour.TermRenderer
}

type InstallerViewModel struct {
	Progress progress.Model
	Spinner  spinner.Model
}

type InstallerModel struct {
	ViewModel      InstallerViewModel
	packages       map[int]string
	notes          map[int]string
	total          int
	line           int
	index          int
	width          int
	height         int
	done           int
	statusInfoText string
	statusPkgText  string
	_cursor        *int // reference to current reading line
}

type T_SocketResponse struct {
	msg string
	err error
}

type execProgram struct {
	string
	args []string
}

type pipeResMsg struct {
	T_OutputMode
	opts string
	string
}

type installPackageMsg struct{ int }
type tickMsg struct{ time.Time }

var Styles *DefinedStyle = &DefinedStyle{
	Main:           lipgloss.NewStyle().Margin(1, 2),
	Keyword:        lipgloss.NewStyle().Foreground(lipgloss.Color("211")),
	Subtle:         lipgloss.NewStyle().Foreground(lipgloss.Color("241")),
	Ticks:          lipgloss.NewStyle().Foreground(lipgloss.Color("79")),
	Checkbox:       lipgloss.NewStyle().Foreground(lipgloss.Color("212")),
	Cursor:         lipgloss.NewStyle().Foreground(lipgloss.Color("63")),
	ProgressEmpty:  lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("░"),
	Dot:            lipgloss.NewStyle().Foreground(lipgloss.Color("236")).Render(" • "),
	RampColor:      makeRampStyles("#B14FFF", "#00FFA3", 71),
	CurrentPkgName: lipgloss.NewStyle().Foreground(lipgloss.Color("211")).Bold(true),
	Done:           lipgloss.NewStyle().Margin(1, 2),
	CheckMark:      lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓"),
}
