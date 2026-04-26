package tui

import (
	"fmt"
	"path"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/miguelgrubin/gocovtui/pkg/coverage"
)

// Model is the Bubble Tea model for the coverage TUI.
type Model struct {
	viewport viewport.Model
	items    []rowData
	summary  coverage.SummaryStats
	cursor   int
	width    int
	height   int
	ready    bool
}

// NewModel creates a new TUI Model from coverage stats.
// Files are grouped by folder: folders sorted by coverage ascending,
// with each folder's files also sorted by coverage ascending.
func NewModel(stats *coverage.Stats) Model {
	var items []rowData
	if stats != nil {
		filesByDir := make(map[string][]*coverage.FileStats)
		for _, f := range stats.FilesSortedByCoverage(true) {
			dir := dirOf(f.Filename)
			filesByDir[dir] = append(filesByDir[dir], f)
		}

		for _, folder := range stats.FoldersSortedByCoverage(true) {
			items = append(items, rowData{
				kind:     kindFolder,
				name:     folder.Dir,
				total:    folder.TotalStatements,
				covered:  folder.CoveredStatements,
				coverPct: folder.CoveragePercent,
			})
			for _, f := range filesByDir[folder.Dir] {
				items = append(items, rowData{
					kind:     kindFile,
					name:     f.Filename,
					total:    f.TotalStatements,
					covered:  f.CoveredStatements,
					coverPct: f.CoveragePercent,
				})
			}
		}
	}

	var summary coverage.SummaryStats
	if stats != nil {
		summary = stats.GetSummary()
	}

	return Model{
		items:   items,
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
		headerH := 3
		m.viewport = viewport.New(msg.Width, msg.Height-headerH)
		m.viewport.Style = lipgloss.NewStyle().Background(colorBackground)
		m.ready = true
		m.refreshTable()

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.refreshTable()
				m.scrollToCursor()
			}
			return m, nil
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
				m.refreshTable()
				m.scrollToCursor()
			}
			return m, nil
		case "home", "g":
			m.cursor = 0
			m.refreshTable()
			m.viewport.GotoTop()
			return m, nil
		case "end", "G":
			m.cursor = len(m.items) - 1
			m.refreshTable()
			m.viewport.GotoBottom()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the full TUI screen.
func (m Model) View() string {
	if !m.ready {
		return "\n  Loading…"
	}
	title := titleStyle.Render("gocovtui")
	header := renderHeader(m.summary, m.width)
	help := lipgloss.NewStyle().Foreground(colorDimGray).Padding(0, 1).
		Render("↑/↓ navigate • q quit")
	return title + "\n" + header + "\n" + m.viewport.View() + "\n" + help
}

// refreshTable rebuilds the lipgloss table and sets it as the viewport content.
func (m *Model) refreshTable() {
	if m.width == 0 {
		return
	}
	items := m.items
	cursor := m.cursor

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(colorPurple)).
		Headers("NAME", "STMTS", "COVERED", "COVERAGE").
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return tableHeaderStyle
			}
			idx := row // 0-based data row index
			if idx < 0 || idx >= len(items) {
				return tableFileStyle
			}
			item := items[idx]
			selected := idx == cursor

			if col == 3 {
				// Coverage column: always color by percentage
				base := coverageStyle(item.coverPct)
				if selected {
					if item.kind == kindFolder {
						return tableSelectedFolderStyle
					}
					return tableSelectedFileStyle
				}
				return base
			}

			if selected {
				if item.kind == kindFolder {
					return tableSelectedFolderStyle
				}
				return tableSelectedFileStyle
			}
			if item.kind == kindFolder {
				return tableFolderStyle
			}
			return tableFileStyle
		}).
		Width(m.width)

	for _, item := range items {
		name := item.name
		if item.kind == kindFolder {
			name = "▶ " + name
		} else {
			name = "  " + name
		}
		t.Row(
			name,
			fmt.Sprintf("%d", item.total),
			fmt.Sprintf("%d", item.covered),
			fmt.Sprintf("%.1f%%", item.coverPct),
		)
	}

	m.viewport.SetContent(t.String())
}

// scrollToCursor ensures the viewport shows the currently selected row.
// Each data row is 1 line tall inside the table (plus 1-line header and borders).
func (m *Model) scrollToCursor() {
	// Table has: top border (1) + header (1) + header-bottom border (1) = 3 lines before data
	lineOffset := 3 + m.cursor
	if lineOffset < m.viewport.YOffset {
		m.viewport.SetYOffset(lineOffset)
	} else if lineOffset >= m.viewport.YOffset+m.viewport.Height {
		m.viewport.SetYOffset(lineOffset - m.viewport.Height + 1)
	}
}

// dirOf returns the parent directory of a file path using forward-slash semantics.
func dirOf(filename string) string {
	return path.Dir(filename)
}

// renderHeader builds the aggregate stats banner.
func renderHeader(s coverage.SummaryStats, width int) string {
	label := headerLabelStyle.Render
	value := headerValueStyle.Render

	covText := fmt.Sprintf("%s %s  %s %s/%s  %s %s",
		label("Coverage:"),
		value(fmt.Sprintf("%.1f%%", s.CoveragePercent)),
		label("Statements:"),
		value(fmt.Sprintf("%d", s.CoveredStatements)),
		value(fmt.Sprintf("%d", s.TotalStatements)),
		label("Files:"),
		value(fmt.Sprintf("%d", s.FileCount)),
	)

	bar := headerStyle.Width(width).Render(covText)
	return bar
}
