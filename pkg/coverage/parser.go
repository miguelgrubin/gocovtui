// Package coverage provides functionality for parsing Go test coverprofile files
// and extracting structured coverage statistics.
package coverage

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Position represents a line and column position in a source file.
type Position struct {
	Line int
	Col  int
}

// StatementCoverage stores individual statement coverage data for a range in a source file.
type StatementCoverage struct {
	Start   Position
	End     Position
	NumStmt int
	Count   int
	Covered bool
}

// FileCoverage stores file-level coverage data including all statements in the file.
type FileCoverage struct {
	Filename   string
	Statements []StatementCoverage
}

// TotalStatements returns the total number of statements (weighted by NumStmt) in the file.
func (f *FileCoverage) TotalStatements() int {
	total := 0
	for _, s := range f.Statements {
		total += s.NumStmt
	}
	return total
}

// CoveredStatements returns the number of covered statements (weighted by NumStmt) in the file.
func (f *FileCoverage) CoveredStatements() int {
	covered := 0
	for _, s := range f.Statements {
		if s.Covered {
			covered += s.NumStmt
		}
	}
	return covered
}

// CoveragePercent returns the coverage percentage for this file.
// Returns 0 if there are no statements.
func (f *FileCoverage) CoveragePercent() float64 {
	total := f.TotalStatements()
	if total == 0 {
		return 0
	}
	return float64(f.CoveredStatements()) / float64(total) * 100
}

// FunctionCoverage represents function-level coverage data derived from statement ranges.
type FunctionCoverage struct {
	Filename string
	Start    Position
	End      Position
	NumStmt  int
	Covered  int
}

// CoveragePercent returns the coverage percentage for this function.
func (fn *FunctionCoverage) CoveragePercent() float64 {
	if fn.NumStmt == 0 {
		return 0
	}
	return float64(fn.Covered) / float64(fn.NumStmt) * 100
}

// CoverageResult holds the complete parsed coverage data.
type CoverageResult struct {
	Mode  string
	Files map[string]*FileCoverage
}

// ExtractFunctions groups statements by proximity into logical function boundaries.
// Statements are considered part of the same function if their line ranges overlap or
// are immediately adjacent (within 1 line of each other).
func (r *CoverageResult) ExtractFunctions() []*FunctionCoverage {
	var functions []*FunctionCoverage

	for filename, file := range r.Files {
		stmts := make([]StatementCoverage, len(file.Statements))
		copy(stmts, file.Statements)

		// Sort statements by start line
		sort.Slice(stmts, func(i, j int) bool {
			if stmts[i].Start.Line != stmts[j].Start.Line {
				return stmts[i].Start.Line < stmts[j].Start.Line
			}
			return stmts[i].Start.Col < stmts[j].Start.Col
		})

		if len(stmts) == 0 {
			continue
		}

		// Group statements into functions by proximity
		groups := groupStatements(stmts)

		for _, group := range groups {
			fn := &FunctionCoverage{
				Filename: filename,
				Start:    group[0].Start,
				End:      group[len(group)-1].End,
			}
			for _, s := range group {
				fn.NumStmt += s.NumStmt
				if s.Covered {
					fn.Covered += s.NumStmt
				}
			}
			functions = append(functions, fn)
		}
	}

	return functions
}

// groupStatements groups adjacent/overlapping statements into logical function groups.
// Statements within the same contiguous block (no line gap > 1) are grouped together.
func groupStatements(stmts []StatementCoverage) [][]StatementCoverage {
	if len(stmts) == 0 {
		return nil
	}

	var groups [][]StatementCoverage
	current := []StatementCoverage{stmts[0]}

	for i := 1; i < len(stmts); i++ {
		prev := current[len(current)-1]
		curr := stmts[i]

		// If this statement starts after a gap from the previous end, start a new group
		if curr.Start.Line > prev.End.Line+1 {
			groups = append(groups, current)
			current = []StatementCoverage{curr}
		} else {
			current = append(current, curr)
		}
	}
	groups = append(groups, current)

	return groups
}

// Parse reads and parses coverprofile data from an io.Reader.
// Returns a CoverageResult or an error if parsing fails.
func Parse(reader io.Reader) (*CoverageResult, error) {
	result := &CoverageResult{
		Files: make(map[string]*FileCoverage),
	}

	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNum++

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Handle mode line (first line)
		if strings.HasPrefix(line, "mode:") {
			result.Mode = strings.TrimSpace(strings.TrimPrefix(line, "mode:"))
			continue
		}

		// Parse coverage line
		stmt, filename, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNum, err)
		}

		if _, ok := result.Files[filename]; !ok {
			result.Files[filename] = &FileCoverage{Filename: filename}
		}
		result.Files[filename].Statements = append(result.Files[filename].Statements, stmt)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading coverprofile: %w", err)
	}

	return result, nil
}

// ParseFile reads and parses a coverprofile from the given file path.
// Returns a CoverageResult or an error if the file cannot be read or parsed.
func ParseFile(path string) (*CoverageResult, error) {
	f, err := os.Open(path) // #nosec G304 -- path is provided by the caller for coverprofile reading
	if err != nil {
		return nil, fmt.Errorf("opening coverprofile %q: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

// ParseCoverprofile parses coverprofile data from a string.
// Returns a CoverageResult or an error if parsing fails.
func ParseCoverprofile(data string) (*CoverageResult, error) {
	return Parse(strings.NewReader(data))
}

// parseLine parses a single coverprofile line of the form:
// <filename>:<startLine>.<startCol>,<endLine>.<endCol> <numStmt> <count>
func parseLine(line string) (StatementCoverage, string, error) {
	// Split into location part and counts
	spaceIdx := strings.Index(line, " ")
	if spaceIdx < 0 {
		return StatementCoverage{}, "", fmt.Errorf("malformed coverprofile line: %q", line)
	}

	location := line[:spaceIdx]
	rest := strings.TrimSpace(line[spaceIdx+1:])

	// Parse numStmt and count from rest
	parts := strings.Fields(rest)
	if len(parts) != 2 {
		return StatementCoverage{}, "", fmt.Errorf("malformed coverprofile line (expected 2 fields after location): %q", line)
	}

	numStmt, err := strconv.Atoi(parts[0])
	if err != nil {
		return StatementCoverage{}, "", fmt.Errorf("malformed numStmt in line %q: %w", line, err)
	}

	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return StatementCoverage{}, "", fmt.Errorf("malformed count in line %q: %w", line, err)
	}

	// Parse location: filename:startLine.startCol,endLine.endCol
	// The filename can contain colons (Windows paths), so find the last colon before the range
	colonIdx := strings.LastIndex(location, ":")
	if colonIdx < 0 {
		return StatementCoverage{}, "", fmt.Errorf("malformed location (no colon) in line %q", line)
	}

	filename := location[:colonIdx]
	rangePart := location[colonIdx+1:]

	commaIdx := strings.Index(rangePart, ",")
	if commaIdx < 0 {
		return StatementCoverage{}, "", fmt.Errorf("malformed range (no comma) in line %q", line)
	}

	start, err := parsePosition(rangePart[:commaIdx])
	if err != nil {
		return StatementCoverage{}, "", fmt.Errorf("malformed start position in line %q: %w", line, err)
	}

	end, err := parsePosition(rangePart[commaIdx+1:])
	if err != nil {
		return StatementCoverage{}, "", fmt.Errorf("malformed end position in line %q: %w", line, err)
	}

	stmt := StatementCoverage{
		Start:   start,
		End:     end,
		NumStmt: numStmt,
		Count:   count,
		Covered: count > 0,
	}

	return stmt, filename, nil
}

// parsePosition parses a position string of the form "line.col".
func parsePosition(s string) (Position, error) {
	dotIdx := strings.Index(s, ".")
	if dotIdx < 0 {
		return Position{}, fmt.Errorf("malformed position (no dot): %q", s)
	}

	line, err := strconv.Atoi(s[:dotIdx])
	if err != nil {
		return Position{}, fmt.Errorf("malformed line number: %w", err)
	}

	col, err := strconv.Atoi(s[dotIdx+1:])
	if err != nil {
		return Position{}, fmt.Errorf("malformed column number: %w", err)
	}

	return Position{Line: line, Col: col}, nil
}
