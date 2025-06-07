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
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	s := spinner.New()
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return InstallerModel{
		packages:  map[int]string{},
		notes:     map[int]string{},
		total:     0,
		ViewModel: InstallerViewModel{p, s},
	}
}

func (m *InstallerModel) UpdateInstaller(pipeData T_PipeData) tea.Cmd {
	packages, notes, total := parsePkgsData(pipeData)

	if m.total == 0 {
		m.total = total
	}
	for i := range len(notes) {
		m.notes[len(m.packages)+i] = notes[i]
	}
	// fmt.Printf("%v\n", m.notes)
	m.packages = lo.Assign(m.packages, packages)

	return processPkgMsg(m)
}

func (m *InstallerModel) SetSize(w int, h int) tea.Cmd {
	m.width = w
	m.height = h

	return nil
}

func (m *InstallerModel) TickCmd() tea.Cmd {
	return tea.Batch(
		processPkgMsg(m),
		m.ViewModel.Spinner.Tick,
	)
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
			m.statusInfoText += fmt.Sprintf("%s\n", Styles.Subtle.Render(m.notes[msg.int-1]))
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
		return Styles.Main.Render()
	}

	if m.line == m.total+len(m.notes) {
		return m.statusInfoText + Styles.Main.Render(
			Styles.Done.Render(fmt.Sprintf("Done! Installed %d packages.\n%v", m.total, m.line)),
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

	return m.statusInfoText + Styles.Main.Render(spin+info+gap+prog+pkgCount)
}

func processPkgMsg(m *InstallerModel) tea.Cmd {
	return tea.Tick(time.Duration(REFRESH_RATE)*time.Millisecond, func(t time.Time) tea.Msg {
		if len(m.packages) <= m.total || m.index <= m.total {
			return installPackageMsg{m.index + 1}
		}
		return nil
	})
}

func parsePkgsData(pipeData T_PipeData) (map[int]string, map[int]string, int) {
	var (
		packages map[int]string = map[int]string{}
		notes    map[int]string = map[int]string{}
		total    int            = 0
	)

	re := regexp2.MustCompile(RGXP_INSTALL_PKGINFO, regexp2.RE2|regexp2.Multiline)

	var curIndex int = 0
	for _, v := range pipeData {
		matches, err := re.FindStringMatch(v[0])
		if err != nil {
			continue
		}

		if matches == nil {
			if notes[curIndex] == "" {
				notes[curIndex] = v[0]
			} else {
				notes[curIndex] = strings.Join([]string{notes[curIndex], v[0]}, "\n")
			}
			continue
		}

		indexPkg, err := strconv.Atoi(matches.GroupByNumber(1).String())
		if err != nil {
			indexPkg = 0
		} else {
			indexPkg = indexPkg - 1
		}
		totalPkg, err := strconv.Atoi(matches.GroupByNumber(2).String())
		if err != nil {
			totalPkg = 0
		}
		total = max(total, totalPkg)
		namePkg := matches.GroupByNumber(3).String()

		curIndex = indexPkg
		packages[indexPkg] = namePkg
	}

	return packages, notes, total
}
