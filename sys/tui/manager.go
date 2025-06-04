package tui

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func tickCmd(speed int) tea.Cmd {
	return tea.Tick(time.Duration(speed)*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{t}
	})
}

func NewTuiManager() (*TuiManager, error) {
	var err error
	tm := &TuiManager{
		CurrentOutputMode: OUTPUT_VIEW_TEXT,
		RefreshRate:       5,
	}

	tm.Header = "It's good on toast"
	tm.PipeData = T_PipeData{}

	tm.Output.Error, err = tm.NewGlamourModel(tm.PipeData)
	if err != nil {
		return tm, err
	}
	tm.Output.Spinner = tm.NewSpinnerModel()
	tm.Output.Text = tm.NewTextModel()
	tm.Output.List = tm.NewListModel()

	tm.SpinnerTick = tm.Output.Spinner.ViewModel.Tick
	tm.CursorBlinkTick = textinput.Blink

	tm.Program = tea.NewProgram(tm, tea.WithAltScreen())
	if _, err := tm.Program.Run(); err != nil {
		return tm, err
	}

	return tm, err
}

func (m TuiManager) Init() tea.Cmd {
	return tickCmd(m.RefreshRate)
}

func (m TuiManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_SPINNER:
			return m, m.SpinnerTick

		case OUTPUT_VIEW_TEXT:
			m.Output.Text.current += 1
			if m.Output.Text.current > len(m.Output.Text.rawText) {
				m.Output.Text.current = len(m.Output.Text.rawText)
				return m, m.CursorBlinkTick
			}
			return m, tickCmd(m.RefreshRate)
		}

	case tea.WindowSizeMsg:
		h, v := Styles.Main.GetFrameSize()
		m.Output.List.ViewModel.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}

	case pipeResMsg:
		pipeData, err := ParseInput(msg.string)
		if err != nil {
			m.CurrentOutputMode = OUTPUT_VIEW_ERROR
			m.PipeData, _ = ParseInput(err.Error())
			return m, m.UpdateOutputModels()
		}
		m.CurrentOutputMode = msg.T_OutputMode
		m.PipeData = pipeData

		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_SPINNER:
			return m, tea.Batch(m.UpdateOutputModels(), m.SpinnerTick)
		case OUTPUT_VIEW_TEXT:
			return m, tea.Batch(m.UpdateOutputModels(), tickCmd(m.RefreshRate))
		}
		return m, m.UpdateOutputModels()
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.Output.Error.ViewModel, cmd = m.Output.Error.ViewModel.Update(msg)
	cmds = append(cmds, cmd)

	m.Output.Spinner.ViewModel, cmd = m.Output.Spinner.ViewModel.Update(msg)
	cmds = append(cmds, cmd)

	m.Output.Text.ViewModel, cmd = m.Output.Text.ViewModel.Update(msg)
	cmds = append(cmds, cmd)

	m.Output.List.ViewModel, cmd = m.Output.List.ViewModel.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TuiManager) View() string {
	switch m.CurrentOutputMode {
	case OUTPUT_VIEW_SPINNER:
		return Styles.Main.Render(
			fmt.Sprintf("\n\n    %s    %s\n\n", m.Output.Spinner.ViewModel.View(), m.Output.Spinner.loadingText),
		)

	case OUTPUT_VIEW_LIST:
		return Styles.Main.Render(m.Output.List.ViewModel.View())

	case OUTPUT_VIEW_TEXT:
		var (
			text    string
			current int    = m.Output.Text.current
			rawText string = m.Output.Text.rawText
		)
		if current <= 1 {
			text = ""
		} else {
			text = rawText[0:current]
		}
		m.Output.Text.ViewModel.SetCursor(current)
		return Styles.Main.Render(fmt.Sprintf("%s%s", text, m.Output.Text.ViewModel.View()))

	case OUTPUT_VIEW_ERROR:
		return Styles.Main.Render(
			m.Output.Error.ViewModel.View(),
			lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\n  ↑/↓: Navigate • q: Quit\n"),
		)
	}
	return ""
}

func (ctx *TuiManager) UpdateOutputModels() tea.Cmd {
	return tea.Batch(
		tea.ClearScreen,
		ctx.Output.Error.UpdateError(ctx.PipeData),
		ctx.Output.Spinner.UpdateSpinner(ctx.PipeData),
		ctx.Output.Text.UpdateText(ctx.PipeData),
		ctx.Output.List.UpdateList(ctx.PipeData),
	)
}

func (m *GlamourModel) UpdateError(PipeData T_PipeData) tea.Cmd {
	_ = m.SetGlamourContent(parseTextData(PipeData))
	return nil
}

// =======================================================================
func (ctx *TuiManager) NewSpinnerModel() SpinnerModel {
	return SpinnerModel{
		ViewModel: spinner.New(
			spinner.WithSpinner(spinner.Meter),
			spinner.WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("205"))),
		),
	}
}

func (m *SpinnerModel) UpdateSpinner(pipeData T_PipeData) tea.Cmd {
	m.loadingText = parseTextData(pipeData)
	return nil
}

// =======================================================================
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

// =======================================================================
func (ctx *TuiManager) NewTextModel() TextModel {
	m := TextModel{}
	m.ViewModel = textinput.New()
	m.UpdateText(ctx.PipeData)

	return m
}

func (m *TextModel) UpdateText(pipeData T_PipeData) tea.Cmd {
	m.ViewModel.Prompt = ""
	m.ViewModel.Cursor.Style = Styles.Cursor
	m.ViewModel.Width = 48
	m.ViewModel.CursorStart()
	m.ViewModel.SetValue("")
	m.ViewModel.CursorEnd()
	m.ViewModel.Focus()

	m.current = 0
	m.rawText = parseTextData(pipeData)

	return nil
}

// =======================================================================
const GLAMOUR_WIDTH = 100
const GLAMOUR_GUTTER = 20

func (ctx *TuiManager) NewGlamourModel(pipeData T_PipeData) (GlamourModel, error) {
	var err error
	m := GlamourModel{
		ViewModel: viewport.New(GLAMOUR_WIDTH, 20),
	}
	m.ViewModel.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("211")).
		PaddingRight(20)

	glamourRenderWidth := GLAMOUR_WIDTH - m.ViewModel.Style.GetHorizontalFrameSize() - GLAMOUR_GUTTER

	m.renderer, err = glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return m, err
	}

	m.SetGlamourContent(parseTextData(pipeData))
	return m, nil
}

func (m *GlamourModel) SetGlamourContent(input string) error {
	str, err := m.renderer.Render(input)
	if err != nil {
		return err
	}
	m.ViewModel.SetContent(str)
	return nil
}

// =======================================================================
func (m *TuiManager) ListenToSocket() {
	socketPath := "/run/tuid.sock"

	// Remove the socket file if it already exists
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// Listen on Unix domain socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("listening on", socketPath)

	var conn net.Conn
	var chRawInput chan *T_SocketResponse = make(chan *T_SocketResponse)

	defer close(chRawInput)
	// defer conn.Close()

	for {
		conn, err = listener.Accept()
		if err != nil {
			// log.Println(err)
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, err.Error()})
			continue
		}

		go ReadFromSocket(conn, chRawInput)
		res := <-chRawInput
		if res.err != nil {
			// log.Println(res.err)
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, res.err.Error()})
			continue
		}
		// log.Println(res.msg)
		// m.Program.Send(pipeResMsg{OUTPUT_VIEW_TEXT, res.msg})
		m.Program.Send(pipeResMsg{OUTPUT_VIEW_TEXT, `Raspberry Pi’s					I have ’em all over my house
Nutella							It's good on toast
Bitter melon					It cools you down
Nice socks						And by that I mean socks without holes
Eight hours of sleep			I had this once
Cats							Usually`})
	}
}

func ToRawData(rd io.Reader) (string, error) {
	var (
		err error
		b   strings.Builder
	)

	reader := bufio.NewReader(rd)

	for {
		r, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		_, err = b.WriteRune(r)
		if err != nil {
			return "", fmt.Errorf("error getting input: %s", err)
		}
	}
	return strings.TrimSpace(b.String()), err
}

func ReadFromSocket[R T_SocketResponse](c net.Conn, res chan *R) {
	for {
		data, err := ToRawData(c)
		res <- &R{msg: data, err: err}
	}
}

func ReadFromPipe() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0 {
		return "", fmt.Errorf("try piping in some text")
	}
	return ToRawData(os.Stdin)
}

func ParseInput(rawInput string) (data T_PipeData, err error) {
	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
		mu      *sync.Mutex = &sync.Mutex{}
	)

	mu.Lock()
	{
		inputScanner := bufio.NewScanner(strings.NewReader(rawInput))

		var sanitized string
		for inputScanner.Scan() {
			line := inputScanner.Text()
			cleanLine := normalizeTSVLine(line)
			sanitized += cleanLine + "\n"
		}
		// log.Println(sanitized)

		reader = strings.NewReader(sanitized)
		scanner = bufio.NewScanner(reader)

		i := 0
		data = T_PipeData{}
		for scanner.Scan() {
			line := scanner.Text()
			columns := strings.Split(line, "\t")
			data[i] = map[int]string{}
			for j, col := range columns {
				data[i][j] = col
			}
			i += 1
		}
	}
	mu.Unlock()

	if err := scanner.Err(); err != nil {
		return data, fmt.Errorf("error reading lines: %s", err)
	}

	return data, nil
}
