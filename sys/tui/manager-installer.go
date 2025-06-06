package tui

import (
	"fmt"
	"home-media/sys/sample"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (ctx *TuiManager) NewInstallerModel() InstallerModel {
	m := InstallerModel{}
	m.ViewModel = progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
		progress.WithoutPercentage(),
	)
	m.SpinnerModel = spinner.New()
	m.SpinnerModel.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return m
}

func (m *InstallerModel) UpdateInstaller(pipeData T_PipeData) tea.Cmd {
	m.packages = parsePkgsData(pipeData)
	return tea.Batch(downloadAndInstall(m.packages[m.index]), m.SpinnerModel.Tick)
}

func (m *InstallerModel) SetSize(w int, h int) tea.Cmd {
	m.width = w
	m.height = h

	return nil
}

func (m *InstallerModel) RenderView() string {
	n := len(m.packages)
	w := lipgloss.Width(fmt.Sprintf("%d", n))

	if m.done {
		return Styles.Done.Render(fmt.Sprintf("Done! Installed %d packages.\n", n))
	}

	pkgCount := fmt.Sprintf(" %*d/%*d", w, m.index, w, n)

	spin := m.SpinnerModel.View() + " "
	prog := m.ProgressModel.View()
	cellsAvail := max(0, m.width-lipgloss.Width(spin+prog+pkgCount))

	pkgName := Styles.CurrentPkgName.Render(m.packages[m.index])
	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render("Installing " + pkgName)

	cellsRemaining := max(0, m.width-lipgloss.Width(spin+info+prog+pkgCount))
	gap := strings.Repeat(" ", cellsRemaining)

	return Styles.Main.Render(spin + info + gap + prog + pkgCount)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func parsePkgsData(pipeData T_PipeData) []string {
	return sample.Sample_InstallPackages
}
