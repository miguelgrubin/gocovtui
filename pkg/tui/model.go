package tui

import (
	"fmt"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/miguelgrubin/gocovtui/pkg/coverage"
)

// Model is the Bubble Tea model for the coverage TUI.
type Model struct {
	list    list.Model
	summary coverage.SummaryStats
	width   int
	height  int
	ready   bool
}

// NewModel creates a new TUI Model from coverage stats.
// Files are grouped by folder: folders sorted by coverage ascending,
// with each folder's files also sorted by coverage ascending.
func NewModel(stats *coverage.Stats) Model {
	items := []list.Item{}
	if stats != nil {
		// Build a map of dir → files for grouping
		filesByDir := make(map[string][]*coverage.FileStats)
		for _, f := range stats.FilesSortedByCoverage(true) {
			dir := dirOf(f.Filename)
			filesByDir[dir] = append(filesByDir[dir], f)
		}

		// Emit folder row then its files, folders ordered by coverage ascending
		for _, folder := range stats.FoldersSortedByCoverage(true) {
			items = append(items, folderItem{
				dir:       folder.Dir,
				coverPct:  folder.CoveragePercent,
				fileCount: folder.FileCount,
				total:     folder.TotalStatements,
				covered:   folder.CoveredStatements,
			})
			for _, f := range filesByDir[folder.Dir] {
				items = append(items, fileItem{
					filename: f.Filename,
					total:    f.TotalStatements,
					covered:  f.CoveredStatements,
					coverPct: f.CoveragePercent,
				})
			}
		}
	}

	l := list.New(items, fileDelegate{}, 0, 0)
	l.Title = "gocovtui"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.SetShowHelp(true)
	l.Styles.Title = titleStyle
	l.Styles.TitleBar = lipgloss.NewStyle().Background(colorBackground).Padding(0, 1)
	l.Styles.NoItems = lipgloss.NewStyle().Foreground(colorGray).Padding(0, 1)

	var summary coverage.SummaryStats
	if stats != nil {
		summary = stats.GetSummary()
	}

	return Model{
		list:    l,
		summary: summary,
	}
}

// Init satisfies tea.Model; no initial commands needed.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles Bubble Tea messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerH := 3 // header lines
		m.list.SetSize(msg.Width, msg.Height-headerH)
		m.ready = true

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the full TUI screen.
func (m Model) View() string {
	if !m.ready {
		return "\n  Loading…"
	}
	return renderHeader(m.summary, m.width) + "\n" + m.list.View()
}

// dirOf returns the parent directory of a file path using forward-slash semantics.
func dirOf(filename string) string {
	return path.Dir(filename)
}

// renderHeader builds the aggregate stats banner.
func renderHeader(s coverage.SummaryStats, width int) string {
	label := headerLabelStyle.Render
	value := headerValueStyle.Render

	coverage := fmt.Sprintf("%s %s  %s %s/%s  %s %s",
		label("Coverage:"),
		value(fmt.Sprintf("%.1f%%", s.CoveragePercent)),
		label("Statements:"),
		value(fmt.Sprintf("%d", s.CoveredStatements)),
		value(fmt.Sprintf("%d", s.TotalStatements)),
		label("Files:"),
		value(fmt.Sprintf("%d", s.FileCount)),
	)

	bar := headerStyle.Width(width).Render(coverage)
	return bar
}
