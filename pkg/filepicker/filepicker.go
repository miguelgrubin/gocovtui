// Package filepicker provides an interactive TUI file picker that filters
// results to *.out files, used when gocovtui is launched without a file argument.
package filepicker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

// SelectedMsg is emitted when the user picks a file.
type SelectedMsg struct {
	Path string
}

// CancelMsg is emitted when the user quits without selecting a file.
type CancelMsg struct{}

// Model is a tea.Model that wraps bubbles/filepicker and restricts
// visible files to those with a ".out" extension.
type Model struct {
	fp       filepicker.Model
	selected string
	quitting bool
}

// NewModel returns an initialised Model. It checks whether any *.out files
// exist in the current working directory and returns an error if none are found.
func NewModel() (Model, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return Model{}, fmt.Errorf("cannot determine working directory: %w", err)
	}

	matches, err := filepath.Glob(filepath.Join(cwd, "*.out"))
	if err != nil {
		return Model{}, fmt.Errorf("glob error: %w", err)
	}
	if len(matches) == 0 {
		return Model{}, fmt.Errorf("no *.out files found in %s", cwd)
	}

	fp := filepicker.New()
	fp.AllowedTypes = []string{".out"}
	fp.CurrentDirectory = cwd
	fp.ShowHidden = false
	fp.DirAllowed = false
	fp.FileAllowed = true

	return Model{fp: fp}, nil
}

// Init initialises the underlying filepicker.
func (m Model) Init() tea.Cmd {
	return m.fp.Init()
}

// Update handles messages and delegates to the underlying filepicker.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.fp, cmd = m.fp.Update(msg)

	if ok, path := m.fp.DidSelectFile(msg); ok {
		m.selected = path
		return m, tea.Quit
	}

	return m, cmd
}

// View renders the file picker.
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	return "\n  Select a coverage file (*.out):\n\n" + m.fp.View()
}

// Result returns the chosen path and whether a file was selected.
func (m Model) Result() (path string, ok bool) {
	return m.selected, m.selected != ""
}
