package coverage

import (
	"path"
	"sort"
)

// FolderStats holds aggregate coverage statistics for a directory/package.
type FolderStats struct {
	Dir               string
	TotalStatements   int
	CoveredStatements int
	CoveragePercent   float64
	FileCount         int
}

// FileStats holds coverage statistics for a single source file.
type FileStats struct {
	Filename          string
	TotalStatements   int
	CoveredStatements int
	CoveragePercent   float64
}

// SummaryStats holds aggregate coverage statistics across all files.
type SummaryStats struct {
	TotalStatements   int
	CoveredStatements int
	CoveragePercent   float64
	FileCount         int
}

// Stats holds per-file and aggregate coverage statistics.
type Stats struct {
	files map[string]*FileStats
}

// NewStats creates a new Stats instance.
func NewStats() *Stats {
	return &Stats{files: make(map[string]*FileStats)}
}

// CalculateStats builds a Stats object from a CoverageResult.
func CalculateStats(result *CoverageResult) *Stats {
	s := NewStats()
	for _, file := range result.Files {
		s.AddFile(file)
	}
	return s
}

// AddFile adds coverage data from a FileCoverage to the Stats.
func (s *Stats) AddFile(file *FileCoverage) {
	total := file.TotalStatements()
	covered := file.CoveredStatements()
	percent := 0.0
	if total > 0 {
		percent = float64(covered) / float64(total) * 100
	}

	if existing, ok := s.files[file.Filename]; ok {
		// Merge: add statements
		existing.TotalStatements += total
		existing.CoveredStatements += covered
		if existing.TotalStatements > 0 {
			existing.CoveragePercent = float64(existing.CoveredStatements) / float64(existing.TotalStatements) * 100
		}
	} else {
		s.files[file.Filename] = &FileStats{
			Filename:          file.Filename,
			TotalStatements:   total,
			CoveredStatements: covered,
			CoveragePercent:   percent,
		}
	}
}

// Merge merges another Stats into this one, combining file statistics.
func (s *Stats) Merge(other *Stats) {
	for _, fs := range other.files {
		if existing, ok := s.files[fs.Filename]; ok {
			existing.TotalStatements += fs.TotalStatements
			existing.CoveredStatements += fs.CoveredStatements
			if existing.TotalStatements > 0 {
				existing.CoveragePercent = float64(existing.CoveredStatements) / float64(existing.TotalStatements) * 100
			}
		} else {
			cp := *fs
			s.files[fs.Filename] = &cp
		}
	}
}

// GetSummary returns aggregate project-wide coverage statistics.
func (s *Stats) GetSummary() SummaryStats {
	var total, covered int
	for _, fs := range s.files {
		total += fs.TotalStatements
		covered += fs.CoveredStatements
	}
	percent := 0.0
	if total > 0 {
		percent = float64(covered) / float64(total) * 100
	}
	return SummaryStats{
		TotalStatements:   total,
		CoveredStatements: covered,
		CoveragePercent:   percent,
		FileCount:         len(s.files),
	}
}

// Files returns all FileStats, sorted by filename by default.
func (s *Stats) Files() []*FileStats {
	result := make([]*FileStats, 0, len(s.files))
	for _, fs := range s.files {
		result = append(result, fs)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Filename < result[j].Filename
	})
	return result
}

// FilesSortedByCoverage returns all FileStats sorted by coverage percentage.
// ascending=true returns lowest coverage first.
func (s *Stats) FilesSortedByCoverage(ascending bool) []*FileStats {
	result := s.Files()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].CoveragePercent < result[j].CoveragePercent
		}
		return result[i].CoveragePercent > result[j].CoveragePercent
	})
	return result
}

// FolderStats returns aggregate coverage statistics grouped by directory.
// The directory key is derived from path.Dir(filename) for each file.
func (s *Stats) FolderStats() []*FolderStats {
	folders := make(map[string]*FolderStats)
	for _, fs := range s.files {
		dir := path.Dir(fs.Filename)
		if f, ok := folders[dir]; ok {
			f.TotalStatements += fs.TotalStatements
			f.CoveredStatements += fs.CoveredStatements
			f.FileCount++
			if f.TotalStatements > 0 {
				f.CoveragePercent = float64(f.CoveredStatements) / float64(f.TotalStatements) * 100
			}
		} else {
			pct := 0.0
			if fs.TotalStatements > 0 {
				pct = float64(fs.CoveredStatements) / float64(fs.TotalStatements) * 100
			}
			folders[dir] = &FolderStats{
				Dir:               dir,
				TotalStatements:   fs.TotalStatements,
				CoveredStatements: fs.CoveredStatements,
				CoveragePercent:   pct,
				FileCount:         1,
			}
		}
	}
	result := make([]*FolderStats, 0, len(folders))
	for _, f := range folders {
		result = append(result, f)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Dir < result[j].Dir
	})
	return result
}

// FoldersSortedByCoverage returns folder stats sorted by coverage percentage.
// ascending=true returns lowest coverage first.
func (s *Stats) FoldersSortedByCoverage(ascending bool) []*FolderStats {
	result := s.FolderStats()
	sort.Slice(result, func(i, j int) bool {
		if ascending {
			return result[i].CoveragePercent < result[j].CoveragePercent
		}
		return result[i].CoveragePercent > result[j].CoveragePercent
	})
	return result
}

// CoverageInRange returns a SummaryStats for statements within the given line range
// of the specified file. startLine and endLine are inclusive.
func (s *Stats) CoverageInRange(filename string, startLine, endLine int, result *CoverageResult) SummaryStats {
	file, ok := result.Files[filename]
	if !ok {
		return SummaryStats{}
	}

	var total, covered int
	for _, stmt := range file.Statements {
		if stmt.Start.Line >= startLine && stmt.End.Line <= endLine {
			total += stmt.NumStmt
			if stmt.Covered {
				covered += stmt.NumStmt
			}
		}
	}

	percent := 0.0
	if total > 0 {
		percent = float64(covered) / float64(total) * 100
	}

	return SummaryStats{
		TotalStatements:   total,
		CoveredStatements: covered,
		CoveragePercent:   percent,
		FileCount:         1,
	}
}
