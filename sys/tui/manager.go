package tui

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dlclark/regexp2"
)

func NewTuiManager() (*TuiManager, error) {
	var err error
	m := &TuiManager{
		CurrentOutputMode: OUTPUT_VIEW_TEXT,
		RefreshRate:       5,
	}

	m.Header = ""
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

		case OUTPUT_VIEW_INSTALLER:
			return m, tea.Batch(
				downloadAndInstall(m.Output.Installer.packages[m.Output.Installer.index]),
				m.Output.Installer.SpinnerModel.Tick,
			)
		}

	case tea.WindowSizeMsg:
		h, v := Styles.Main.GetFrameSize()
		m.Output.List.ViewModel.SetSize(msg.Width-h, msg.Height-v)
		m.Output.Installer.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, tea.Quit
		}
		go func() {
			switch m.CurrentOutputMode {
			case OUTPUT_VIEW_SPINNER:
			case OUTPUT_VIEW_LIST:
				m.Output.List.BindExtraKeyCommands(m, msg)
			}
		}()

	case spinner.TickMsg:
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_INSTALLER:
			var cmd tea.Cmd
			m.Output.Installer.SpinnerModel, cmd = m.Output.Installer.SpinnerModel.Update(msg)
			return m, cmd
		}

	case progress.FrameMsg:
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_INSTALLER:
			newModel, cmd := m.Output.Installer.ProgressModel.Update(msg)
			if newModel, ok := newModel.(progress.Model); ok {
				m.Output.Installer.ProgressModel = newModel
			}
			return m, cmd
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
		case OUTPUT_VIEW_LIST:
			m.Output.List.CommandExec = msg.opts
			return m, m.UpdateOutputModels()
		case OUTPUT_VIEW_INSTALLER:
			return m, m.UpdateOutputModels()
		}
		return m, m.UpdateOutputModels()

	case installedPkgMsg:
		pkg := m.Output.Installer.packages[m.Output.Installer.index]
		if m.Output.Installer.index >= len(m.Output.Installer.packages)-1 {
			// Everything's been installed. We're done!
			m.Output.Installer.done = true
			return m, tea.Sequence(
				tea.Printf("%s %s", Styles.CheckMark, pkg), // print the last success message
				tea.ClearScreen,
			)
		}

		// Update progress bar
		m.Output.Installer.index++
		progressCmd := m.Output.Installer.ProgressModel.SetPercent(float64(m.Output.Installer.index) / float64(len(m.Output.Installer.packages)))

		return m, tea.Batch(
			progressCmd,
			tea.Printf("%s %s", Styles.CheckMark, pkg),                                // print success message above our program
			downloadAndInstall(m.Output.Installer.packages[m.Output.Installer.index]), // download the next package
		)
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

	// m.Output.Installer.ViewModel, cmd = m.Output.Installer.ViewModel.Update(msg)
	// cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m TuiManager) View() string {
	switch m.CurrentOutputMode {
	case OUTPUT_VIEW_ERROR:
		return m.Output.Error.RenderErrorView()

	case OUTPUT_VIEW_SPINNER:
		return m.Output.Spinner.RenderView()

	case OUTPUT_VIEW_TEXT:
		return m.Output.Text.RenderView()

	case OUTPUT_VIEW_LIST:
		return m.Output.List.RenderView()

	case OUTPUT_VIEW_INSTALLER:
		return m.Output.Installer.RenderView()
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
		ctx.Output.Installer.UpdateInstaller(ctx.PipeData),
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
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, "", err.Error()})
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
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, "", res.err.Error()})
			close(chRes)
			continue
		}

		re := regexp2.MustCompile(`^((\d+(?=\|))((?=\|)..[^\|\n]*|)((?=\|).*)|.*)$`, regexp2.RE2|regexp2.Singleline)
		matches, err := re.FindStringMatch(res.msg)
		if err != nil {
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, "", res.err.Error()})
			close(chRes)
			continue
		}

		var (
			messagePayload string = ""
			modePayload    string = ""
			optsPayload    string = ""
		)

		if matches.GroupByNumber(4).String() != "" {
			messagePayload = strings.Trim(matches.GroupByNumber(4).String(), "|")
		} else {
			messagePayload = matches.GroupByNumber(1).String()
		}

		if matches.GroupByNumber(2).String() != "" {
			modePayload = matches.GroupByNumber(2).String()
		} else {
			modePayload = OUTPUT_VIEW_TEXT.String()
		}

		if matches.GroupByNumber(3).String() != "" {
			optsPayload = strings.Trim(matches.GroupByNumber(3).String(), "|")
		}

		outputMode, err := strconv.Atoi(modePayload)
		if err != nil {
			m.Program.Send(pipeResMsg{OUTPUT_VIEW_ERROR, "", err.Error()})
			close(chRes)
			return
		}
		m.Program.Send(pipeResMsg{T_OutputMode(outputMode), optsPayload, messagePayload})
		close(chRes)
	}
}

func tickCmd(speed int) tea.Cmd {
	return tea.Tick(time.Duration(speed)*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg{t}
	})
}

func downloadAndInstall(pkgName string) tea.Cmd {
	return tea.Tick(time.Duration(5)*time.Second, func(t time.Time) tea.Msg {
		return installedPkgMsg{pkgName}
	})
}
