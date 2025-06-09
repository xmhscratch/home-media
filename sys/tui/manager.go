package tui

import (
	"fmt"
	"home-media/sys/command"
	"home-media/sys/runtime"
	"net"
	"os"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dlclark/regexp2"
)

func NewTuiManager() (*TuiManager, error) {
	var err error
	m := &TuiManager{
		CurrentOutputMode: OUTPUT_VIEW_LIST,
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
	m.Output.Installer = m.NewInstallerModel()

	return m, err
}

func (m TuiManager) Init() tea.Cmd {
	return tickCmd(REFRESH_RATE)
}

func (m TuiManager) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_SPINNER:
			return m, m.Output.Spinner.TickCmd()
		case OUTPUT_VIEW_TEXT:
			return m, m.Output.Text.TickCmd()
		case OUTPUT_VIEW_LIST:
			return m, m.Output.List.TickCmd()
		case OUTPUT_VIEW_INSTALLER:
			return m, m.Output.Installer.TickCmd()
		}

	case tea.WindowSizeMsg:
		h, v := Styles.Main.GetFrameSize()
		m.Output.List.SetSize(msg.Width-h, msg.Height-v)
		m.Output.Installer.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, tea.Quit
		}
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_LIST:
			var cmd tea.Cmd
			m.Output.List.ViewModel, cmd = m.Output.List.ViewModel.Update(msg)
			return m, tea.Batch(cmd, m.Output.List.BindExtraKeyCommands(m, msg))
		case OUTPUT_VIEW_INSTALLER:
			return m, m.Output.Installer.BindExtraKeyCommands(m, msg)
		}

	case pipeResMsg:
		pipeData, err := ParseInput(msg.string)
		if err != nil {
			m.CurrentOutputMode = OUTPUT_VIEW_ERROR
			m.PipeData, _ = ParseInput(err.Error())
			return m, m.UpdateOutputModels(true)
		}
		var withFreshScreen bool = m.CurrentOutputMode != msg.T_OutputMode
		m.CurrentOutputMode = msg.T_OutputMode
		m.PipeData = pipeData

		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_SPINNER:
			return m, tea.Sequence(tea.ExitAltScreen, m.UpdateOutputModels(withFreshScreen), m.Output.Spinner.TickCmd())
		case OUTPUT_VIEW_TEXT:
			return m, tea.Sequence(tea.ExitAltScreen, m.UpdateOutputModels(withFreshScreen), m.Output.Text.TickCmd())
		case OUTPUT_VIEW_LIST:
			m.Output.List.CommandExec = msg.opts
			return m, tea.Sequence(tea.ExitAltScreen, m.UpdateOutputModels(withFreshScreen))
			// return m, tea.Sequence(tea.EnterAltScreen, m.UpdateOutputModels(withFreshScreen))
		case OUTPUT_VIEW_INSTALLER:
			return m, tea.Sequence(tea.ExitAltScreen, m.UpdateOutputModels(withFreshScreen))
		}
		return m, tea.Sequence(tea.ExitAltScreen, m.UpdateOutputModels(withFreshScreen))

	case execProgram:
		var exitCode chan int = make(chan int)

		go func() {
			stdin := command.NewCommandReader()
			stdout := command.NewNullWriter()
			stderr := command.NewNullWriter()

			shell := runtime.Shell{
				PID: os.Getpid(),

				Stdin:  stdin,
				Stdout: stdout,
				Stderr: stderr,

				Args: os.Args,

				Main: CmdExecShell,
			}

			stdin.WriteVar("ExecBin", msg.string)
			stdin.WriteVar("ExecArgs", msg.args)

			exitCode <- shell.Run()
		}()
		<-exitCode
		close(exitCode)

		return m, tea.Quit

	default:
		switch m.CurrentOutputMode {
		case OUTPUT_VIEW_INSTALLER:
			return m, m.Output.Installer.BindExtraCustomCommands(m, msg)
		}
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

func (m *TuiManager) UpdateOutputModels(withFreshScreen bool) tea.Cmd {
	var cmds []tea.Cmd
	if withFreshScreen {
		cmds = append(
			cmds,
			// m.Output.Error.Reset(),
			// m.Output.Spinner.Reset(),
			m.Output.Text.Reset(),
			m.Output.List.Reset(),
			m.Output.Installer.Reset(),
		)
	}
	// fmt.Printf("%v\n", m.PipeData)
	cmds = append(
		cmds,
		m.Output.Error.UpdateError(m.PipeData),
		m.Output.Spinner.UpdateSpinner(m.PipeData),
		m.Output.Text.UpdateText(m.PipeData),
		m.Output.List.UpdateList(m.PipeData),
		m.Output.Installer.UpdateInstaller(m.PipeData),
	)
	return tea.Sequence(cmds...)
}

func (m *TuiManager) Start() error {
	m.Program = tea.NewProgram(m)
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

		re := regexp2.MustCompile(RGXP_MESSAGE_PAYLOAD, regexp2.RE2|regexp2.Singleline)
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
			messagePayload = trimRS(matches.GroupByNumber(4).String())
		} else {
			messagePayload = matches.GroupByNumber(1).String()
		}

		if matches.GroupByNumber(2).String() != "" {
			modePayload = matches.GroupByNumber(2).String()
		} else {
			modePayload = OUTPUT_VIEW_TEXT.String()
		}

		if matches.GroupByNumber(3).String() != "" {
			optsPayload = trimRS(matches.GroupByNumber(3).String())
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
