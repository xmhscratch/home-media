package tui

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func NewTuiManager() (*TuiManager, error) {
	var err error
	m := &TuiManager{
		CurrentOutputMode: OUTPUT_VIEW_TEXT,
		RefreshRate:       5,
	}

	m.Header = "It's good on toast"
	m.PipeData = T_PipeData{}

	m.Output.Error, err = m.NewGlamourModel(m.PipeData)
	if err != nil {
		return m, err
	}
	m.Output.Spinner = m.NewSpinnerModel()
	m.Output.Text = m.NewTextModel()
	m.Output.List = m.NewListModel()

	m.SpinnerTick = m.Output.Spinner.ViewModel.Tick
	m.CursorBlinkTick = textinput.Blink

	return m, err
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
		log.Println(pipeData, err)
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
		return m, nil
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

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

func (m *TuiManager) Start() error {
	m.Program = tea.NewProgram(m, tea.WithAltScreen())
	if _, err := m.Program.Run(); err != nil {
		return err
	}
	return nil
}

func (m *TuiManager) ListenToSocket() {
	socketPath := "/run/tuid.sock"
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("listening on", socketPath)

	for {
		conn, err := listener.Accept()
		if err != nil {
			// log.Println(err)
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, err.Error()})
			continue
		}

		var (
			res   *T_SocketResponse
			chRes chan *T_SocketResponse = make(chan *T_SocketResponse)
		)

		go func() {
			data, err := ToRawData(conn)
			chRes <- &T_SocketResponse{msg: data, err: err}
		}()

		res = <-chRes
		if res.err != nil {
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, err.Error()})
			close(chRes)
			continue
		}
		m.Program.Send(pipeResMsg{OUTPUT_VIEW_TEXT, res.msg})
		close(chRes)
	}
}

func tickCmd(speed int) tea.Cmd {
	return tea.Tick(time.Duration(speed)*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{t}
	})
}
