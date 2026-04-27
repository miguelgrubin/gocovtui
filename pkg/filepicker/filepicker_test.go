package filepicker

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- NewModel ---

func TestNewModel_NoOutFiles(t *testing.T) {
	// Create a temp dir with no .out files
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("cannot chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	_, err := NewModel()
	if err == nil {
		t.Error("expected error when no *.out files exist")
	}
}

func TestNewModel_WithOutFiles(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a .out file
	if err := os.WriteFile(filepath.Join(tmpDir, "coverage.out"), []byte("mode: set\n"), 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	origDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("cannot chdir: %v", err)
	}
	t.Cleanup(func() { os.Chdir(origDir) })

	m, err := NewModel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.quitting {
		t.Error("expected quitting=false on new model")
	}
	if m.selected != "" {
		t.Error("expected selected to be empty on new model")
	}
}

// --- Init ---

func TestInit_ReturnsCmd(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.out"), []byte("mode: set\n"), 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	m, err := NewModel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cmd := m.Init()
	if cmd == nil {
		t.Error("expected Init() to return a command for the underlying filepicker")
	}
}

// --- Update: quit ---

func TestUpdate_QuitOnQ(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.out"), []byte("mode: set\n"), 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	m, _ := NewModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	um := updated.(Model)
	if !um.quitting {
		t.Error("expected quitting=true after 'q'")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

func TestUpdate_QuitOnCtrlC(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.out"), []byte("mode: set\n"), 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	m, _ := NewModel()
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	um := updated.(Model)
	if !um.quitting {
		t.Error("expected quitting=true after ctrl+c")
	}
	if cmd == nil {
		t.Error("expected quit command")
	}
}

// --- View ---

func TestView_Quitting(t *testing.T) {
	m := Model{quitting: true}
	view := m.View()
	if view != "" {
		t.Errorf("expected empty view when quitting, got %q", view)
	}
}

func TestView_Normal(t *testing.T) {
	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "test.out"), []byte("mode: set\n"), 0644); err != nil {
		t.Fatalf("cannot create test file: %v", err)
	}
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	t.Cleanup(func() { os.Chdir(origDir) })

	m, _ := NewModel()
	view := m.View()
	if view == "" {
		t.Error("expected non-empty view")
	}
}

// --- Result ---

func TestResult_NoSelection(t *testing.T) {
	m := Model{selected: ""}
	path, ok := m.Result()
	if ok {
		t.Error("expected ok=false when no file selected")
	}
	if path != "" {
		t.Errorf("expected empty path, got %q", path)
	}
}

func TestResult_WithSelection(t *testing.T) {
	m := Model{selected: "/tmp/coverage.out"}
	path, ok := m.Result()
	if !ok {
		t.Error("expected ok=true when file selected")
	}
	if path != "/tmp/coverage.out" {
		t.Errorf("expected '/tmp/coverage.out', got %q", path)
	}
}

// --- SelectedMsg and CancelMsg types ---

func TestSelectedMsg(t *testing.T) {
	msg := SelectedMsg{Path: "/some/path.out"}
	if msg.Path != "/some/path.out" {
		t.Errorf("expected path '/some/path.out', got %q", msg.Path)
	}
}

func TestCancelMsg(t *testing.T) {
	_ = CancelMsg{} // just ensure the type exists
}
