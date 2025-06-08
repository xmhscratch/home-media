package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dlclark/regexp2"
	"github.com/samber/lo"
)

func (ctx *TuiManager) NewInstallerModel() InstallerModel {
	return newInstallerModel()
}

func newInstallerModel() InstallerModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	var _index int = 0
	m := InstallerModel{
		ViewModel:      InstallerViewModel{p, s},
		packages:       map[int]string{},
		notes:          map[int]string{},
		total:          0,
		line:           0,
		index:          _index,
		width:          0,
		height:         0,
		done:           0,
		statusInfoText: "",
		statusPkgText:  "",
		_cursor:        &(_index),
	}
	return m
}

func (m *InstallerModel) UpdateInstaller(pipeData T_PipeData) tea.Cmd {
	packages, notes, total := insertPackagesData(pipeData, m._cursor)

	if m.total != 0 && m.index == m.total && (total != 0 && m.total != total || *m._cursor < m.index) {
		m.done = REFRESH_RATE_IN_SECONDS * 3
		return m.TickCmd()
	}

	if m.total == 0 {
		m.total = total
	}
	m.packages = lo.Assign(m.packages, packages)

	mergedNotes := lo.MapValues(notes, func(v string, k int) string {
		if _, ok := m.notes[k]; !ok {
			return v
		} else {
			return strings.Join([]string{"", v}, "\n")
		}
	})
	m.notes = lo.Assign(m.notes, mergedNotes)

	return processPkgMsg(m)
}

func (m *InstallerModel) Reset() tea.Cmd {
	*m = newInstallerModel()

	return tea.Batch(
		tea.ClearScreen,
		tea.WindowSize(),
	)
}

func (m *InstallerModel) SetSize(w int, h int) {
	m.width = w
	m.height = h
}

func (m *InstallerModel) TickCmd() tea.Cmd {
	if m.done > 0 {
		m.done -= 1
		if m.done <= 0 {
			return m.Reset()
		}
		return tickCmd(REFRESH_RATE)
	}
	return tea.Batch(
		processPkgMsg(m),
		m.ViewModel.Spinner.Tick,
	)
}

func (m *InstallerModel) BindExtraKeyCommands(mgr TuiManager, msg tea.KeyMsg) tea.Cmd {
	if msg.String() == "enter" {
		return m.Reset()
	}
	return nil
}

func (m *InstallerModel) BindExtraCustomCommands(mgr TuiManager, msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.ViewModel.Spinner, cmd = m.ViewModel.Spinner.Update(msg)
		return cmd

	case progress.FrameMsg:
		newModel, cmd := m.ViewModel.Progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.ViewModel.Progress = newModel
		}
		return cmd

	case installPackageMsg:
		if m.total == 0 {
			return nil
		}

		if m.index > len(m.packages)-1 {
			m.line += len(m.notes)
			m.statusInfoText += fmt.Sprintf("%s\n", Styles.Subtle.Render(m.notes[*m._cursor]))
			return nil
		}

		pkgName := m.packages[m.index]
		m.index = msg.int
		m.line += 1

		m.statusInfoText += lipgloss.NewStyle().Render(fmt.Sprintf("%s %s\n", Styles.CheckMark, pkgName))
		m.statusPkgText = fmt.Sprintf("Installing %s", Styles.CurrentPkgName.Render(pkgName))

		return m.ViewModel.Progress.SetPercent(float64(m.index) / float64(m.total))
	}
	return nil
}

func (m *InstallerModel) RenderView() string {
	if m.total == 0 {
		return Styles.Main.Render("")
	}

	if m.line == m.total+len(m.notes) {
		return Styles.Main.Render(m.statusInfoText) + Styles.Main.Render(
			Styles.Done.Render(fmt.Sprintf("Done! Installed %d packages. Closing in %ds...\n", m.total, (m.done/REFRESH_RATE_IN_SECONDS)+1)),
		)
	}

	w := lipgloss.Width(fmt.Sprintf("%d", m.total))
	spin := m.ViewModel.Spinner.View() + " "
	prog := m.ViewModel.Progress.View()

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, m.total)
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(m.statusPkgText)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return Styles.Main.Render(m.statusInfoText) + Styles.Main.Render(spin+info+gap+prog+pkgCount)
}

func processPkgMsg(m *InstallerModel) tea.Cmd {
	return tea.Tick(time.Duration(REFRESH_RATE)*time.Millisecond, func(t time.Time) tea.Msg {
		if len(m.packages) <= m.total || m.index <= m.total {
			return installPackageMsg{m.index + 1}
		}
		return nil
	})
}

func insertPackagesData(pipeData T_PipeData, _cursor *int) (map[int]string, map[int]string, int) {
	var (
		packages map[int]string = map[int]string{}
		notes    map[int]string = map[int]string{}
		total    int            = 0
	)

	re := regexp2.MustCompile(RGXP_INSTALL_PKGINFO, regexp2.RE2|regexp2.Multiline)

	for _, v := range pipeData {
		matches, err := re.FindStringMatch(v[0])
		if err != nil {
			continue
		}

		if matches == nil {
			if _, ok := notes[*_cursor]; !ok {
				notes[*_cursor] = v[0]
			} else {
				notes[*_cursor] = strings.Join([]string{notes[*_cursor], v[0]}, "\n")
			}
			continue
		}

		pkgIndex, err := strconv.Atoi(matches.GroupByNumber(1).String())
		if err != nil {
			pkgIndex = 0
		} else {
			pkgIndex = pkgIndex - 1
		}
		*_cursor = pkgIndex

		pkgTotal, err := strconv.Atoi(matches.GroupByNumber(2).String())
		if err != nil {
			pkgTotal = 0
		}
		total = max(total, pkgTotal)
		pkgName := matches.GroupByNumber(3).String()
		packages[*_cursor] = pkgName
	}

	return packages, notes, total
}
