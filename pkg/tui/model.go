package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

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

	// Find mode state
	findMode      bool
	searchTerm    string
	filteredItems []rowData
	filteredMap   []int // maps filtered index -> original items index

	// Module path prefix to strip for relative file paths
	modulePath string

	// Error/status message
	statusMsg string
}

// editorFinishedMsg signals that the external editor process has exited.
type editorFinishedMsg struct{ err error }

// NewModel creates a new TUI Model from coverage stats.
// Files are grouped by folder: folders sorted by coverage ascending,
// with each folder's files also sorted by coverage ascending.
// modulePath is the Go module path (from go.mod) used to resolve relative file paths.
func NewModel(stats *coverage.Stats, modulePath string) Model {
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
		items:      items,
		summary:    summary,
		modulePath: modulePath,
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

	case editorFinishedMsg:
		m.statusMsg = ""
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Editor error: %v", msg.err)
		}
		return m, nil

	case tea.KeyMsg:
		// Clear status message on any key press
		m.statusMsg = ""

		if m.findMode {
			return m.updateFindMode(msg)
		}

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
		case "f":
			m.findMode = true
			m.searchTerm = ""
			m.applyFilter()
			m.refreshTable()
			return m, nil
		case "e":
			return m.handleEditAction()
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// updateFindMode handles key events while in find mode.
func (m Model) updateFindMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "ctrl+c":
		m.findMode = false
		m.searchTerm = ""
		m.filteredItems = nil
		m.filteredMap = nil
		m.refreshTable()
		return m, nil
	case "enter":
		if len(m.filteredItems) > 0 && m.cursor < len(m.filteredMap) {
			m.cursor = m.filteredMap[m.cursor]
		}
		m.findMode = false
		m.searchTerm = ""
		m.filteredItems = nil
		m.filteredMap = nil
		m.refreshTable()
		m.scrollToCursor()
		return m, nil
	case "up", "shift+tab":
		if m.cursor > 0 {
			m.cursor--
			m.refreshTable()
			m.scrollToCursor()
		}
		return m, nil
	case "down", "tab":
		max := len(m.filteredItems) - 1
		if max < 0 {
			max = len(m.items) - 1
		}
		if m.cursor < max {
			m.cursor++
			m.refreshTable()
			m.scrollToCursor()
		}
		return m, nil
	case "backspace":
		if len(m.searchTerm) > 0 {
			m.searchTerm = m.searchTerm[:len(m.searchTerm)-1]
			m.applyFilter()
			m.refreshTable()
		}
		return m, nil
	default:
		// Only add printable single characters
		if len(msg.String()) == 1 {
			m.searchTerm += msg.String()
			m.applyFilter()
			m.refreshTable()
		}
		return m, nil
	}
}

// applyFilter filters items based on the current search term.
func (m *Model) applyFilter() {
	if m.searchTerm == "" {
		m.filteredItems = nil
		m.filteredMap = nil
		m.cursor = 0
		return
	}
	term := strings.ToLower(m.searchTerm)
	m.filteredItems = nil
	m.filteredMap = nil
	for i, item := range m.items {
		if strings.Contains(strings.ToLower(item.name), term) {
			m.filteredItems = append(m.filteredItems, item)
			m.filteredMap = append(m.filteredMap, i)
		}
	}
	m.cursor = 0
}

// handleEditAction opens the selected file in the configured EDITOR.
func (m Model) handleEditAction() (tea.Model, tea.Cmd) {
	if m.cursor < 0 || m.cursor >= len(m.items) {
		m.statusMsg = "No file selected"
		return m, nil
	}
	item := m.items[m.cursor]

	editor := os.Getenv("EDITOR")
	if editor == "" {
		m.statusMsg = "EDITOR environment variable not set"
		return m, nil
	}

	relPath, err := m.toRelativePath(item.name)
	if err != nil {
		m.statusMsg = fmt.Sprintf("Cannot resolve path: %v", err)
		return m, nil
	}

	c := exec.Command(editor, relPath)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

// toRelativePath converts a full module path (e.g. "github.com/user/repo/pkg/file.go")
// to a relative path from the current directory (e.g. "./pkg/file.go").
func (m Model) toRelativePath(fullPath string) (string, error) {
	if m.modulePath == "" {
		return "", fmt.Errorf("module path not configured")
	}
	rel, found := strings.CutPrefix(fullPath, m.modulePath)
	if !found {
		return "", fmt.Errorf("path %q does not belong to module %q", fullPath, m.modulePath)
	}
	// rel starts with "/" (e.g. "/pkg/file.go"), prepend "."
	return "." + rel, nil
}

// View renders the full TUI screen.
func (m Model) View() string {
	if !m.ready {
		return "\n  Loading…"
	}
	title := titleStyle.Render("gocovtui")
	header := renderHeader(m.summary, m.width)

	var help string
	if m.findMode {
		findBar := findLabelStyle.Render("Find:") + " " + findInputStyle.Render(m.searchTerm+"█")
		if m.searchTerm != "" && len(m.filteredItems) == 0 {
			findBar += " " + findNoResultsStyle.Render("no matches")
		} else if m.searchTerm != "" {
			findBar += " " + findLabelStyle.Render(fmt.Sprintf("(%d matches)", len(m.filteredItems)))
		}
		help = findBar + "\n" + lipgloss.NewStyle().Foreground(colorDimGray).Padding(0, 1).
			Render("↑/↓ navigate • enter select • esc cancel")
	} else {
		helpText := "↑/↓ navigate • f find • e edit • q quit"
		if m.statusMsg != "" {
			helpText = errorStyle.Render(m.statusMsg)
		}
		help = lipgloss.NewStyle().Foreground(colorDimGray).Padding(0, 1).
			Render(helpText)
	}

	return title + "\n" + header + "\n" + m.viewport.View() + "\n" + help
}

// refreshTable rebuilds the lipgloss table and sets it as the viewport content.
func (m *Model) refreshTable() {
	if m.width == 0 {
		return
	}
	items := m.items
	if m.findMode && m.filteredItems != nil {
		items = m.filteredItems
	}
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
