package coverage_test

import (
	"testing"

	"github.com/miguelgrubin/gocovtui/pkg/coverage"
)

const multiFileCoverprofile = `mode: set
github.com/x/a.go:1.1,5.10 10 8
github.com/x/a.go:6.1,10.5 5 0
github.com/x/b.go:1.1,3.5 4 4
github.com/x/c.go:1.1,2.2 2 1
`

func parseMulti(t *testing.T) *coverage.CoverageResult {
	t.Helper()
	result, err := coverage.ParseCoverprofile(multiFileCoverprofile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

func TestCalculateStats_FileCount(t *testing.T) {
	result := parseMulti(t)
	stats := coverage.CalculateStats(result)
	summary := stats.GetSummary()
	if summary.FileCount != 3 {
		t.Errorf("expected 3 files, got %d", summary.FileCount)
	}
}

func TestCalculateStats_TotalStatements(t *testing.T) {
	result := parseMulti(t)
	stats := coverage.CalculateStats(result)
	summary := stats.GetSummary()
	// a.go: 10+5=15, b.go: 4, c.go: 2 → total=21
	if summary.TotalStatements != 21 {
		t.Errorf("expected 21 total statements, got %d", summary.TotalStatements)
	}
}

func TestCalculateStats_CoveredStatements(t *testing.T) {
	result := parseMulti(t)
	stats := coverage.CalculateStats(result)
	summary := stats.GetSummary()
	// a.go: stmt1 covered(10), stmt2 uncovered(5); b.go: covered(4); c.go: covered(2)
	// covered = 10+4+2 = 16
	if summary.CoveredStatements != 16 {
		t.Errorf("expected 16 covered statements, got %d", summary.CoveredStatements)
	}
}

func TestCalculateStats_CoveragePercent(t *testing.T) {
	result := parseMulti(t)
	stats := coverage.CalculateStats(result)
	summary := stats.GetSummary()
	expected := float64(16) / float64(21) * 100
	if summary.CoveragePercent < expected-0.01 || summary.CoveragePercent > expected+0.01 {
		t.Errorf("expected %.2f%% coverage, got %.2f%%", expected, summary.CoveragePercent)
	}
}

func TestFileStats_ZeroStatements(t *testing.T) {
	data := "mode: set\ngithub.com/x/empty.go:1.1,1.2 0 0\n"
	result, _ := coverage.ParseCoverprofile(data)
	stats := coverage.CalculateStats(result)
	files := stats.Files()
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].CoveragePercent != 0 {
		t.Errorf("expected 0%% for zero-statement file, got %.2f%%", files[0].CoveragePercent)
	}
}

func TestFileStats_PartialCoverage(t *testing.T) {
	data := "mode: set\ngithub.com/x/f.go:1.1,2.2 4 2\ngithub.com/x/f.go:3.1,4.2 4 0\n"
	result, _ := coverage.ParseCoverprofile(data)
	stats := coverage.CalculateStats(result)
	files := stats.Files()
	// 4 covered out of 8 total = 50%
	if files[0].CoveragePercent < 49.99 || files[0].CoveragePercent > 50.01 {
		t.Errorf("expected 50%% coverage, got %.2f%%", files[0].CoveragePercent)
	}
}

func TestStats_FilesSortedByCoverage(t *testing.T) {
	result := parseMulti(t)
	stats := coverage.CalculateStats(result)

	asc := stats.FilesSortedByCoverage(true)
	for i := 1; i < len(asc); i++ {
		if asc[i-1].CoveragePercent > asc[i].CoveragePercent {
			t.Errorf("ascending sort broken at index %d", i)
		}
	}

	desc := stats.FilesSortedByCoverage(false)
	for i := 1; i < len(desc); i++ {
		if desc[i-1].CoveragePercent < desc[i].CoveragePercent {
			t.Errorf("descending sort broken at index %d", i)
		}
	}
}

func TestStats_Merge(t *testing.T) {
	data1 := "mode: set\ngithub.com/x/a.go:1.1,2.2 4 4\n"
	data2 := "mode: set\ngithub.com/x/a.go:3.1,4.2 4 0\n"

	r1, _ := coverage.ParseCoverprofile(data1)
	r2, _ := coverage.ParseCoverprofile(data2)

	s1 := coverage.CalculateStats(r1)
	s2 := coverage.CalculateStats(r2)
	s1.Merge(s2)

	summary := s1.GetSummary()
	if summary.TotalStatements != 8 {
		t.Errorf("expected 8 total after merge, got %d", summary.TotalStatements)
	}
	if summary.CoveredStatements != 4 {
		t.Errorf("expected 4 covered after merge, got %d", summary.CoveredStatements)
	}
}

func TestStats_CoverageInRange(t *testing.T) {
	data := `mode: set
github.com/x/f.go:1.1,5.10 3 3
github.com/x/f.go:10.1,15.5 5 0
`
	result, _ := coverage.ParseCoverprofile(data)
	stats := coverage.CalculateStats(result)

	// Query lines 1-5 only
	rangeStats := stats.CoverageInRange("github.com/x/f.go", 1, 5, result)
	if rangeStats.TotalStatements != 3 {
		t.Errorf("expected 3 statements in range, got %d", rangeStats.TotalStatements)
	}
	if rangeStats.CoveredStatements != 3 {
		t.Errorf("expected 3 covered in range, got %d", rangeStats.CoveredStatements)
	}
}

func TestNewStats_Empty(t *testing.T) {
	stats := coverage.NewStats()
	summary := stats.GetSummary()
	if summary.FileCount != 0 || summary.TotalStatements != 0 {
		t.Error("expected empty stats for new Stats")
	}
}

// multiDirCoverprofile has files in two different directories.
const multiDirCoverprofile = `mode: set
github.com/x/pkg/a.go:1.1,5.10 10 8
github.com/x/pkg/b.go:1.1,3.5 4 4
github.com/x/other/c.go:1.1,2.2 2 1
`

func parseMultiDir(t *testing.T) *coverage.CoverageResult {
	t.Helper()
	result, err := coverage.ParseCoverprofile(multiDirCoverprofile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return result
}

func TestFolderStats_Count(t *testing.T) {
	stats := coverage.CalculateStats(parseMultiDir(t))
	folders := stats.FolderStats()
	if len(folders) != 2 {
		t.Errorf("expected 2 folders, got %d", len(folders))
	}
}

func TestFolderStats_Aggregation(t *testing.T) {
	stats := coverage.CalculateStats(parseMultiDir(t))
	folders := stats.FolderStats()

	var pkg *coverage.FolderStats
	for _, f := range folders {
		if f.Dir == "github.com/x/pkg" {
			pkg = f
		}
	}
	if pkg == nil {
		t.Fatal("expected folder github.com/x/pkg not found")
	}
	// a.go: 10 stmts, 8 covered; b.go: 4 stmts, 4 covered → total 14/14 covered → 100%? No: 12/14
	if pkg.TotalStatements != 14 {
		t.Errorf("expected 14 total statements in pkg, got %d", pkg.TotalStatements)
	}
	if pkg.CoveredStatements != 14 {
		t.Errorf("expected 14 covered statements in pkg, got %d", pkg.CoveredStatements)
	}
	if pkg.FileCount != 2 {
		t.Errorf("expected 2 files in pkg, got %d", pkg.FileCount)
	}
	expected := float64(14) / float64(14) * 100
	if pkg.CoveragePercent < expected-0.01 || pkg.CoveragePercent > expected+0.01 {
		t.Errorf("expected %.2f%% coverage for pkg, got %.2f%%", expected, pkg.CoveragePercent)
	}
}

func TestFolderStats_Empty(t *testing.T) {
	stats := coverage.NewStats()
	folders := stats.FolderStats()
	if len(folders) != 0 {
		t.Errorf("expected empty folder stats, got %d", len(folders))
	}
}

func TestFoldersSortedByCoverage_Ascending(t *testing.T) {
	stats := coverage.CalculateStats(parseMultiDir(t))
	folders := stats.FoldersSortedByCoverage(true)
	for i := 1; i < len(folders); i++ {
		if folders[i-1].CoveragePercent > folders[i].CoveragePercent {
			t.Errorf("ascending sort broken at index %d", i)
		}
	}
}

func TestFoldersSortedByCoverage_Descending(t *testing.T) {
	stats := coverage.CalculateStats(parseMultiDir(t))
	folders := stats.FoldersSortedByCoverage(false)
	for i := 1; i < len(folders); i++ {
		if folders[i-1].CoveragePercent < folders[i].CoveragePercent {
			t.Errorf("descending sort broken at index %d", i)
		}
	}
}
