package coverage_test

import (
	"strings"
	"testing"

	"github.com/miguelgrubin/gocovtui/pkg/coverage"
)

const validCoverprofile = `mode: set
github.com/miguelgrubin/gocovtui/pkg/app.go:10.5,12.20 3 1
github.com/miguelgrubin/gocovtui/pkg/app.go:15.3,18.7 2 0
github.com/miguelgrubin/gocovtui/pkg/foo.go:5.1,8.10 4 1
`

func TestParseCoverprofile_Valid(t *testing.T) {
	result, err := coverage.ParseCoverprofile(validCoverprofile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != "set" {
		t.Errorf("expected mode 'set', got %q", result.Mode)
	}
	if len(result.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(result.Files))
	}
	appFile, ok := result.Files["github.com/miguelgrubin/gocovtui/pkg/app.go"]
	if !ok {
		t.Fatal("expected app.go in result")
	}
	if len(appFile.Statements) != 2 {
		t.Errorf("expected 2 statements in app.go, got %d", len(appFile.Statements))
	}
}

func TestParseCoverprofile_ModeAtomic(t *testing.T) {
	data := "mode: atomic\ngithub.com/x/y.go:1.1,2.2 1 5\n"
	result, err := coverage.ParseCoverprofile(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Mode != "atomic" {
		t.Errorf("expected mode 'atomic', got %q", result.Mode)
	}
}

func TestParseCoverprofile_StatementCovered(t *testing.T) {
	data := "mode: set\ngithub.com/x/y.go:1.1,2.2 2 1\ngithub.com/x/y.go:3.1,4.2 1 0\n"
	result, err := coverage.ParseCoverprofile(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	file := result.Files["github.com/x/y.go"]
	if !file.Statements[0].Covered {
		t.Error("expected first statement to be covered")
	}
	if file.Statements[1].Covered {
		t.Error("expected second statement to be uncovered")
	}
}

func TestParseCoverprofile_MultipleFiles(t *testing.T) {
	result, err := coverage.ParseCoverprofile(validCoverprofile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, hasApp := result.Files["github.com/miguelgrubin/gocovtui/pkg/app.go"]
	_, hasFoo := result.Files["github.com/miguelgrubin/gocovtui/pkg/foo.go"]
	if !hasApp || !hasFoo {
		t.Error("expected both files in result")
	}
}

func TestParseCoverprofile_Malformed_MissingSpace(t *testing.T) {
	data := "mode: set\ngithub.com/x/y.go:1.1,2.2\n"
	_, err := coverage.ParseCoverprofile(data)
	if err == nil {
		t.Error("expected error for malformed line, got nil")
	}
}

func TestParseCoverprofile_Malformed_BadFields(t *testing.T) {
	data := "mode: set\ngithub.com/x/y.go:1.1,2.2 abc def\n"
	_, err := coverage.ParseCoverprofile(data)
	if err == nil {
		t.Error("expected error for non-numeric fields, got nil")
	}
}

func TestParseCoverprofile_Malformed_BadPosition(t *testing.T) {
	data := "mode: set\ngithub.com/x/y.go:bad.pos,2.2 1 0\n"
	_, err := coverage.ParseCoverprofile(data)
	if err == nil {
		t.Error("expected error for bad position, got nil")
	}
}

func TestParseFile_NotFound(t *testing.T) {
	_, err := coverage.ParseFile("/nonexistent/path/to/coverage.out")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestParse_FromReader(t *testing.T) {
	reader := strings.NewReader(validCoverprofile)
	result, err := coverage.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(result.Files))
	}
}

func TestParse_EmptyReader(t *testing.T) {
	reader := strings.NewReader("")
	result, err := coverage.Parse(reader)
	if err != nil {
		t.Fatalf("unexpected error for empty reader: %v", err)
	}
	if len(result.Files) != 0 {
		t.Errorf("expected 0 files, got %d", len(result.Files))
	}
}

func TestFileCoverage_TotalAndCoveredStatements(t *testing.T) {
	result, _ := coverage.ParseCoverprofile(validCoverprofile)
	appFile := result.Files["github.com/miguelgrubin/gocovtui/pkg/app.go"]

	if appFile.TotalStatements() != 5 { // 3 + 2
		t.Errorf("expected 5 total statements, got %d", appFile.TotalStatements())
	}
	if appFile.CoveredStatements() != 3 { // only first stmt covered
		t.Errorf("expected 3 covered statements, got %d", appFile.CoveredStatements())
	}
}

func TestExtractFunctions(t *testing.T) {
	data := `mode: set
github.com/x/y.go:1.1,3.10 2 1
github.com/x/y.go:2.5,3.9 1 1
github.com/x/y.go:10.1,12.5 3 0
`
	result, err := coverage.ParseCoverprofile(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fns := result.ExtractFunctions()
	if len(fns) != 2 {
		t.Errorf("expected 2 functions, got %d", len(fns))
	}
}
