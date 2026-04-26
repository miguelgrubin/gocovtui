package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// fileItem represents a single file's coverage data as a list item.
type fileItem struct {
	filename string
	total    int
	covered  int
	coverPct float64
}

func (f fileItem) FilterValue() string { return f.filename }
func (f fileItem) Title() string       { return f.filename }
func (f fileItem) Description() string {
	return fmt.Sprintf("stmts: %d  covered: %d", f.total, f.covered)
}

// folderItem represents a directory/package summary row in the list.
type folderItem struct {
	dir       string
	coverPct  float64
	fileCount int
	total     int
	covered   int
}

func (f folderItem) FilterValue() string { return f.dir }
func (f folderItem) Title() string       { return f.dir }
func (f folderItem) Description() string {
	return fmt.Sprintf("%d files  %d/%d stmts", f.fileCount, f.covered, f.total)
}

// fileDelegate is a custom list delegate that renders file and folder items with Synthwave styles.
type fileDelegate struct{}

func (d fileDelegate) Height() int                             { return 1 }
func (d fileDelegate) Spacing() int                            { return 0 }
func (d fileDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }

func (d fileDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	selected := index == m.Index()

	switch item := listItem.(type) {
	case folderItem:
		d.renderFolder(w, m, item, selected)
	case fileItem:
		d.renderFile(w, m, item, selected)
	}
}

func (d fileDelegate) renderFolder(w io.Writer, m list.Model, fi folderItem, selected bool) {
	pctStr := fmt.Sprintf("%5.1f%%", fi.coverPct)
	coloredPct := coverageStyle(fi.coverPct).Render(pctStr)

	dir := fi.dir
	maxNameLen := m.Width() - 22
	if maxNameLen < 10 {
		maxNameLen = 10
	}
	if len(dir) > maxNameLen {
		dir = "…" + dir[len(dir)-maxNameLen+1:]
	}

	padding := maxNameLen - len([]rune(dir))
	if padding < 1 {
		padding = 1
	}

	stmtStr := fmt.Sprintf("%4d stmts", fi.total)
	line := fmt.Sprintf("▶ %s%s%s  %s", dir, strings.Repeat(" ", padding), stmtStr, coloredPct)

	if selected {
		fmt.Fprint(w, folderSelectedStyle.Width(m.Width()).Render(line)) //nolint
	} else {
		fmt.Fprint(w, folderRowStyle.Width(m.Width()).Render(line)) //nolint
	}
}

func (d fileDelegate) renderFile(w io.Writer, m list.Model, fi fileItem, selected bool) {
	pctStr := fmt.Sprintf("%5.1f%%", fi.coverPct)
	coloredPct := coverageStyle(fi.coverPct).Render(pctStr)

	stmtStr := fmt.Sprintf("%4d stmts", fi.total)

	name := fi.filename
	maxNameLen := m.Width() - 20
	if maxNameLen < 10 {
		maxNameLen = 10
	}
	if len(name) > maxNameLen {
		name = "…" + name[len(name)-maxNameLen+1:]
	}

	padding := maxNameLen - len([]rune(name))
	if padding < 1 {
		padding = 1
	}

	line := fmt.Sprintf("  %s%s%s  %s", name, strings.Repeat(" ", padding), stmtStr, coloredPct)

	if selected {
		fmt.Fprint(w, itemSelectedStyle.Render(line)) //nolint
	} else {
		fmt.Fprint(w, itemNormalStyle.Render(line)) //nolint
	}
}
