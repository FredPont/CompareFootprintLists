package fileutil

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	dirListA   = "list_A/"
	dirListB   = "list_B/"
	dirResults = "results/"
	logFile    = "output.log"
)

// FileRecord holds the three columns for one entry in a footprint list.
type FileRecord struct {
	Footprint string // column 0: hash / checksum
	Key       string // column 1 or 2: comparison key (filename or path)
	Path      string // column 2: original filesystem path
}

// ── Public entry point ────────────────────────────────────────────────────────

// ReadLists reads both footprint lists, compares them, and writes result files.
func ReadLists() {
	args := Config.Arg

	// Initialise a clean log file for this run
	logger, closeLog := newLogger()
	defer closeLog()

	la, lb := GetLists()

	// Load both lists concurrently
	type result struct {
		data      map[string]FileRecord
		count     int
		dupCount  int
	}
	chA := make(chan result, 1)
	chB := make(chan result, 1)

	go func() {
		data := loadList(dirListA+la, args)
		m, dups := buildFileMap(data, logger)
		chA <- result{m, len(data), dups}
	}()
	go func() {
		data := loadList(dirListB+lb, args)
		m, dups := buildFileMap(data, logger)
		chB <- result{m, len(data), dups}
	}()

	resA := <-chA
	resB := <-chB

	logInfo(logger, fmt.Sprintf("%d files in %s", resA.count, la))
	logInfo(logger, fmt.Sprintf("%d files in %s", resB.count, lb))

	diffCount, diffs, commonCount, commons := compareMaps(resA.data, resB.data)

	// Write common files table (always, when there are matches)
	if commonCount > 0 {
		writeCSV(commons, dirResults+"common.csv")
		logInfo(logger, fmt.Sprintf("%d common file(s) written to results/common.csv", commonCount))
	}

	// Write differences table
	if diffCount > 0 {
		writeCSV(diffs, dirResults+"diff.csv")
		logInfo(logger, fmt.Sprintf("%d difference(s) written to results/diff.csv", diffCount))
	}

	// Print summary to terminal
	PrintSummary(resA.count, resB.count, diffCount, commonCount, la, lb)
	PrintDuplicateAlert(resA.dupCount, resB.dupCount)
}

// ── List loading ──────────────────────────────────────────────────────────────

// loadList reads one TSV footprint file and returns raw 3-column rows.
// It delegates to the trim variant when path trimming is active.
func loadList(path string, args Args) [][]string {
	if Config.TrimPath {
		return loadListTrimPath(path, args)
	}
	return loadListRaw(path, args)
}

// loadListRaw reads a TSV file and returns [footprint, key, path] rows.
func loadListRaw(path string, args Args) [][]string {
	keyCol := keyColumnIndex(args)
	return parseTSV(path, func(line []string) []string {
		return []string{line[0], line[keyCol], pathColumn(line)}
	})
}

// loadListTrimPath reads a TSV file and trims the leading directories from the
// key column before returning [footprint, trimmedKey, originalPath] rows.
func loadListTrimPath(path string, args Args) [][]string {
	keyCol := keyColumnIndex(args)
	trimIdx := trimIndexFor(path)

	PrintSection(fmt.Sprintf("Trimming path — index %d — %s", trimIdx, path))

	return parseTSV(path, func(line []string) []string {
		trimmed := ReconstructPathByIndex(
			removeLeadingSlash(line[keyCol]),
			trimIdx,
			Config.CommonDirSep,
		)
		return []string{line[0], trimmed, pathColumn(line)}
	})
}

// parseTSV opens a tab-separated file, skips comment lines (#), and applies
// transform to each row. Returns nil on open error.
func parseTSV(path string, transform func([]string) []string) [][]string {
	file, err := os.Open(path)
	if err != nil {
		PrintError(fmt.Sprintf("Cannot open %s: %v", path, err))
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.Comment = '#'
	reader.LazyQuotes = true

	var rows [][]string
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			PrintWarning(fmt.Sprintf("Skipping malformed line in %s: %v", path, err))
			continue
		}
		rows = append(rows, transform(line))
	}
	return rows
}

// ── Map building ──────────────────────────────────────────────────────────────

// buildFileMap converts a 3-column row slice into a key→FileRecord map.
// Duplicate keys are counted and logged.
func buildFileMap(rows [][]string, logger *log.Logger) (map[string]FileRecord, int) {
	m := make(map[string]FileRecord, len(rows))
	duplicates := 0

	for _, row := range rows {
		if len(row) < 2 {
			continue
		}
		key := row[1]
		path := ""
		if len(row) > 2 {
			path = row[2]
		}
		if _, exists := m[key]; exists {
			logInfo(logger, "Duplicate file: "+key)
			duplicates++
		}
		m[key] = FileRecord{Footprint: row[0], Key: key, Path: path}
	}
	return m, duplicates
}

// ── Comparison ────────────────────────────────────────────────────────────────

// compareMaps compares two FileRecord maps and returns:
//   - diffCount  : number of differing or missing entries
//   - diffs      : TSV rows for results/diff.csv
//   - commonCount: number of identical entries
//   - commons    : TSV rows for results/common.csv
func compareMaps(mapA, mapB map[string]FileRecord) (int, [][]string, int, [][]string) {
	diffs := [][]string{{"file", "footprint_A", "footprint_B"}}
	commons := [][]string{{"file", "footprint", "path_in_A", "path_in_B"}}
	diffCount := 0
	commonCount := 0

	// Iterate A: find differences and common files
	for key, recA := range mapA {
		recB, inB := mapB[key]
		switch {
		case !inB:
			// Present in A, absent from B
			diffs = append(diffs, []string{key, recA.Footprint, ""})
			diffCount++
		case recA.Footprint != recB.Footprint:
			// Present in both but signatures differ
			diffs = append(diffs, []string{key, recA.Footprint, recB.Footprint})
			diffCount++
		default:
			// Identical — record with paths from both lists
			commons = append(commons, []string{key, recA.Footprint, recA.Path, recB.Path})
			commonCount++
		}
	}

	// Iterate B: find entries absent from A (not yet counted)
	for key, recB := range mapB {
		if _, inA := mapA[key]; !inA {
			diffs = append(diffs, []string{key, "", recB.Footprint})
			diffCount++
		}
	}

	return diffCount, diffs, commonCount, commons
}

// ── TSV output ────────────────────────────────────────────────────────────────

// writeCSV writes a slice of rows as a tab-separated file.
func writeCSV(data [][]string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		PrintError(fmt.Sprintf("Cannot create %s: %v", filePath, err))
		return
	}
	defer file.Close()

	w := csv.NewWriter(file)
	w.Comma = '\t'
	defer w.Flush()

	for _, row := range data {
		if err := w.Write(row); err != nil {
			PrintWarning(fmt.Sprintf("Write error in %s: %v", filePath, err))
		}
	}
	PrintSuccess("Written: " + filePath)
}

// ── TSV header reader ─────────────────────────────────────────────────────────

// ReadTSVHeader reads and returns the first three lines of a TSV file.
func ReadTSVHeader(path string) [][]string {
	PrintSection("Reading header: " + path)
	var records [][]string

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for i := 0; scanner.Scan() && i < 3; i++ {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		records = append(records, fields)
		PrintInfo(strings.Join(fields, "\t"))
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return records
}

// ── Logging ───────────────────────────────────────────────────────────────────

// newLogger creates a file-only logger for output.log (no ANSI codes in file).
// Returns the logger and a close function to defer.
func newLogger() (*log.Logger, func()) {
	f, err := os.Create(logFile)
	if err != nil {
		log.Fatal("Cannot create log file:", err)
	}
	logger := log.New(f, "", log.LstdFlags)
	return logger, func() { f.Close() }
}

// logInfo writes a plain-text message to the log file AND prints it to terminal.
func logInfo(logger *log.Logger, msg string) {
	fmt.Println(msg)
	logger.Println(msg)
}

// ── Small helpers ─────────────────────────────────────────────────────────────

// keyColumnIndex returns 1 (filename) or 2 (path) depending on config.
func keyColumnIndex(args Args) int {
	if args.ComparisonCriteria == "path" {
		return 2
	}
	return 1
}

// pathColumn safely returns column 2, or "" if the row is too short.
func pathColumn(line []string) string {
	if len(line) > 2 {
		return line[2]
	}
	return ""
}

// trimIndexFor returns the correct trim index depending on which list is processed.
func trimIndexFor(path string) int {
	if strings.Contains(path, "list_A") {
		return Config.TrimIndexPathA
	}
	return Config.TrimIndexPathB
}
