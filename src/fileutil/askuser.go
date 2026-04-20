package fileutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AskUser asks whether the user wants to trim leading path segments before
// comparison. Only relevant when -p path is active.
func AskUser() {
	if Config.Arg.ComparisonCriteria != "path" {
		return
	}

	PrintSection("Path trimming")
	PrintInfo("Do you want to trim the leading directories of the paths? [y/n]")

	if readYN() {
		Config.TrimPath = true
		interactivePathTrim()
	} else {
		Config.TrimPath = false
	}
}

// interactivePathTrim reads the first lines of both lists, auto-detects the
// common directory, shows the result, and lets the user refine the trim index.
func interactivePathTrim() {
	la, lb := GetLists()
	headA := ReadTSVHeader(dirListA + la)
	headB := ReadTSVHeader(dirListB + lb)

	if len(headA) == 0 || len(headB) == 0 || len(headA[0]) < 3 || len(headB[0]) < 3 {
		PrintError("Could not read path column from list headers.")
		return
	}

	path1 := removeLeadingSlash(headA[0][2])
	path2 := removeLeadingSlash(headB[0][2])

	// Auto-detect common prefix and show suggestion
	newPath1, newPath2 := RemoveLeadingDirs(path1, path2)
	PrintPathTrim(
		"Suggested pathA", newPath1, Config.TrimIndexPathA,
		"Suggested pathB", newPath2, Config.TrimIndexPathB,
	)

	refineTrimIndex(path1, path2)
}

// refineTrimIndex asks the user if they want to adjust the auto-detected trim
// indices, and loops until satisfied.
func refineTrimIndex(path1, path2 string) {
	PrintInfo("Do you want to adjust the number of directories to trim? [y/n]")
	if !readYN() {
		return
	}

	Config.TrimIndexPathA = askNumber("Number of directories to trim in list_A: ")
	Config.TrimIndexPathB = askNumber("Number of directories to trim in list_B: ")

	newPath1 := ReconstructPathByIndex(path1, Config.TrimIndexPathA, Config.CommonDirSep)
	newPath2 := ReconstructPathByIndex(path2, Config.TrimIndexPathB, Config.CommonDirSep)

	PrintPathTrim(
		"New pathA", newPath1, Config.TrimIndexPathA,
		"New pathB", newPath2, Config.TrimIndexPathB,
	)

	// Recurse to allow further refinement
	refineTrimIndex(path1, path2)
}

// ── Input helpers ─────────────────────────────────────────────────────────────

// readYN reads a single line from stdin and returns true if the user typed "y".
func readYN() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(strings.ToLower(input)) == "y"
}

// askNumber prints a prompt and reads an integer from stdin.
func askNumber(prompt string) int {
	PrintInfo(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil {
		PrintWarning(fmt.Sprintf("Invalid number %q — using 0", input))
		return 0
	}
	return n
}
