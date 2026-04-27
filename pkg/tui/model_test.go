package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/miguelgrubin/gocovtui/pkg/coverage"
)

// --- helpers ---

const testCoverprofile = `mode: set
github.com/example/mod/pkg/a/foo.go:1.1,5.10 10 8
github.com/example/mod/pkg/a/foo.go:6.1,10.5 5 0
github.com/example/mod/pkg/b/bar.go:1.1,3.5 4 4
github.com/example/mod/pkg/b/baz.go:1.1,2.2 2 1
`

func parseTestStats(t *testing.T) *coverage.Stats {
	t.Helper()
	result, err := coverage.ParseCoverprofile(testCoverprofile)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}
	return coverage.CalculateStats(result)
}

func newTestModel(t *testing.T) Model {
	t.Helper()
	return NewModel(parseTestStats(t), "github.com/example/mod")
}

// makeReady sends a WindowSizeMsg to initialise the viewport.
func makeReady(t *testing.T, m Model) Model {
	t.Helper()
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	return updated.(Model)
}

// --- NewModel ---

func TestNewModel_NilStats(t *testing.T) {
	m := NewModel(nil, "")
	if len(m.items) != 0 {
		t.Errorf("expected 0 items for nil stats, got %d", len(m.items))
	}
}

func TestNewModel_PopulatesItems(t *testing.T) {
	m := newTestModel(t)
	if len(m.items) == 0 {
		t.Fatal("expected items to be populated")
	}
	// Should contain both folders and files
	var folders, files int
	for _, item := range m.items {
		switch item.kind {
		case kindFolder:
			folders++
		case kindFile:
			files++
		}
	}
	if folders == 0 {
		t.Error("expected at least one folder row")
	}
	if files == 0 {
		t.Error("expected at least one file row")
	}
}

func TestNewModel_SummaryPopulated(t *testing.T) {
	m := newTestModel(t)
	if m.summary.FileCount == 0 {
		t.Error("expected summary FileCount > 0")
	}
	if m.summary.TotalStatements == 0 {
		t.Error("expected summary TotalStatements > 0")
	}
}

func TestNewModel_ModulePath(t *testing.T) {
	m := NewModel(nil, "github.com/example/mod")
	if m.modulePath != "github.com/example/mod" {
		t.Errorf("expected modulePath to be set, got %q", m.modulePath)
	}
}

// --- Init ---

func TestInit_ReturnsNil(t *testing.T) {
	m := newTestModel(t)
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}

// --- Update: navigation ---

func TestUpdate_CursorDown(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	um := updated.(Model)
	if um.cursor != 1 {
		t.Errorf("expected cursor=1 after down, got %d", um.cursor)
	}
}

func TestUpdate_CursorUp(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	um := updated.(Model)
	if um.cursor != 0 {
		t.Errorf("expected cursor=0 after up, got %d", um.cursor)
	}
}

func TestUpdate_CursorUpAtTop(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	um := updated.(Model)
	if um.cursor != 0 {
		t.Errorf("expected cursor to stay at 0, got %d", um.cursor)
	}
}

func TestUpdate_CursorDownAtBottom(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = len(m.items) - 1
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	um := updated.(Model)
	if um.cursor != len(m.items)-1 {
		t.Errorf("expected cursor to stay at bottom, got %d", um.cursor)
	}
}

func TestUpdate_Home(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 3
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	um := updated.(Model)
	if um.cursor != 0 {
		t.Errorf("expected cursor=0 after home, got %d", um.cursor)
	}
}

func TestUpdate_End(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}})
	um := updated.(Model)
	if um.cursor != len(m.items)-1 {
		t.Errorf("expected cursor at end (%d), got %d", len(m.items)-1, um.cursor)
	}
}

func TestUpdate_Quit(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	if cmd == nil {
		t.Error("expected quit command")
	}
}

// --- Update: find mode ---

func TestUpdate_EnterFindMode(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
	um := updated.(Model)
	if !um.findMode {
		t.Error("expected findMode to be true")
	}
	if um.searchTerm != "" {
		t.Error("expected searchTerm to be empty on enter")
	}
}

func TestUpdate_FindModeEsc(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.findMode = true
	m.searchTerm = "test"
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	um := updated.(Model)
	if um.findMode {
		t.Error("expected findMode to be false after esc")
	}
	if um.searchTerm != "" {
		t.Error("expected searchTerm to be cleared")
	}
}

func TestUpdate_FindModeTyping(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.findMode = true
	m.searchTerm = ""
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})
	um := updated.(Model)
	if um.searchTerm != "a" {
		t.Errorf("expected searchTerm='a', got %q", um.searchTerm)
	}
}

func TestUpdate_FindModeBackspace(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.findMode = true
	m.searchTerm = "ab"
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	um := updated.(Model)
	if um.searchTerm != "a" {
		t.Errorf("expected searchTerm='a', got %q", um.searchTerm)
	}
}

func TestUpdate_FindModeBackspaceEmpty(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.findMode = true
	m.searchTerm = ""
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	um := updated.(Model)
	if um.searchTerm != "" {
		t.Errorf("expected searchTerm to remain empty, got %q", um.searchTerm)
	}
}

// --- applyFilter ---

func TestApplyFilter_EmptyTerm(t *testing.T) {
	m := newTestModel(t)
	m.searchTerm = ""
	m.applyFilter()
	if m.filteredItems != nil {
		t.Error("expected nil filteredItems for empty search")
	}
}

func TestApplyFilter_MatchesFound(t *testing.T) {
	m := newTestModel(t)
	m.searchTerm = "foo"
	m.applyFilter()
	if len(m.filteredItems) == 0 {
		t.Error("expected matches for 'foo'")
	}
	for _, item := range m.filteredItems {
		if item.name == "" {
			t.Error("filtered item should have a name")
		}
	}
}

func TestApplyFilter_CaseInsensitive(t *testing.T) {
	m := newTestModel(t)
	m.searchTerm = "FOO"
	m.applyFilter()
	count := len(m.filteredItems)
	m.searchTerm = "foo"
	m.applyFilter()
	if len(m.filteredItems) != count {
		t.Error("expected case-insensitive matching")
	}
}

func TestApplyFilter_NoMatches(t *testing.T) {
	m := newTestModel(t)
	m.searchTerm = "zzzznonexistent"
	m.applyFilter()
	if len(m.filteredItems) != 0 {
		t.Errorf("expected 0 matches, got %d", len(m.filteredItems))
	}
}

func TestApplyFilter_ResetsCursor(t *testing.T) {
	m := newTestModel(t)
	m.cursor = 5
	m.searchTerm = "foo"
	m.applyFilter()
	if m.cursor != 0 {
		t.Errorf("expected cursor reset to 0, got %d", m.cursor)
	}
}

// --- toRelativePath ---

func TestToRelativePath_Success(t *testing.T) {
	m := Model{modulePath: "github.com/example/mod"}
	rel, err := m.toRelativePath("github.com/example/mod/pkg/file.go")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rel != "./pkg/file.go" {
		t.Errorf("expected './pkg/file.go', got %q", rel)
	}
}

func TestToRelativePath_NoModulePath(t *testing.T) {
	m := Model{modulePath: ""}
	_, err := m.toRelativePath("github.com/example/mod/pkg/file.go")
	if err == nil {
		t.Error("expected error when modulePath is empty")
	}
}

func TestToRelativePath_PathNotInModule(t *testing.T) {
	m := Model{modulePath: "github.com/example/mod"}
	_, err := m.toRelativePath("github.com/other/mod/pkg/file.go")
	if err == nil {
		t.Error("expected error when path does not belong to module")
	}
}

// --- handleEditAction ---

func TestHandleEditAction_NoItems(t *testing.T) {
	m := Model{cursor: -1, items: nil}
	updated, cmd := m.handleEditAction()
	um := updated.(Model)
	if cmd != nil {
		t.Error("expected no command for invalid cursor")
	}
	if um.statusMsg == "" {
		t.Error("expected status message for no file selected")
	}
}

func TestHandleEditAction_NoEditor(t *testing.T) {
	m := newTestModel(t)
	m.cursor = 0
	t.Setenv("EDITOR", "")
	updated, cmd := m.handleEditAction()
	um := updated.(Model)
	if cmd != nil {
		t.Error("expected no command when EDITOR is not set")
	}
	if um.statusMsg == "" {
		t.Error("expected status message about EDITOR not set")
	}
}

// --- editorFinishedMsg ---

func TestUpdate_EditorFinishedSuccess(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.statusMsg = "something"
	updated, _ := m.Update(editorFinishedMsg{err: nil})
	um := updated.(Model)
	if um.statusMsg != "" {
		t.Errorf("expected empty statusMsg, got %q", um.statusMsg)
	}
}

// --- View ---

func TestView_NotReady(t *testing.T) {
	m := newTestModel(t)
	view := m.View()
	if view != "\n  Loading…" {
		t.Errorf("expected loading message, got %q", view)
	}
}

func TestView_Ready(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	view := m.View()
	if view == "\n  Loading…" {
		t.Error("expected rendered view, got loading message")
	}
	if len(view) == 0 {
		t.Error("expected non-empty view")
	}
}

// --- dirOf ---

func TestDirOf(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/x/pkg/file.go", "github.com/x/pkg"},
		{"file.go", "."},
		{"a/b/c/d.go", "a/b/c"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := dirOf(tt.input)
			if got != tt.expected {
				t.Errorf("dirOf(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// --- WindowSizeMsg ---

func TestUpdate_WindowSizeMsg(t *testing.T) {
	m := newTestModel(t)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	um := updated.(Model)
	if !um.ready {
		t.Error("expected ready=true after WindowSizeMsg")
	}
	if um.width != 100 {
		t.Errorf("expected width=100, got %d", um.width)
	}
	if um.height != 50 {
		t.Errorf("expected height=50, got %d", um.height)
	}
}

// --- scrollToCursor ---

func TestScrollToCursor_NoScroll(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.cursor = 0
	m.scrollToCursor()
	// Should not panic and offset should be reasonable
}

// --- refreshTable with zero width ---

func TestRefreshTable_ZeroWidth(t *testing.T) {
	m := newTestModel(t)
	m.width = 0
	m.refreshTable() // should be a no-op, not panic
}

// --- Find mode enter selects ---

func TestFindMode_EnterSelectsItem(t *testing.T) {
	m := makeReady(t, newTestModel(t))
	m.findMode = true
	m.searchTerm = "bar"
	m.applyFilter()
	if len(m.filteredItems) == 0 {
		t.Skip("no matches for 'bar' in test data")
	}
	m.cursor = 0
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	um := updated.(Model)
	if um.findMode {
		t.Error("expected findMode to be false after enter")
	}
	// cursor should be mapped to original index
	if um.cursor < 0 || um.cursor >= len(um.items) {
		t.Errorf("cursor %d out of bounds", um.cursor)
	}
}
