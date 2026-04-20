package fileutil

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// ── Colour functions ──────────────────────────────────────────────────────────
// Each is a Sprint-style function: returns a coloured string without printing.
// This lets us compose them freely and keep log output free of ANSI codes.

var (
	cRed     = color.New(color.FgRed).SprintFunc()
	cGreen   = color.New(color.FgGreen).SprintFunc()
	cYellow  = color.New(color.FgYellow).SprintFunc()
	cMagenta = color.New(color.FgMagenta).SprintFunc()
	cCyan    = color.New(color.FgCyan).SprintFunc()
	cBlue    = color.New(color.FgBlue).SprintFunc()
	cBold    = color.New(color.Bold).SprintFunc()
)

const separatorChar = "━"
const separatorLen = 58

// separator returns a full-width coloured line.
func separator() string {
	return cRed(strings.Repeat(separatorChar, separatorLen))
}

// ── Public print helpers ──────────────────────────────────────────────────────

// PrintSeparator prints a red horizontal rule.
func PrintSeparator() {
	fmt.Println(separator())
}

// PrintTitle prints the application banner.
func PrintTitle() {
	fmt.Println()
	fmt.Println(cCyan("   ┌──────────────────────────────────────────┐"))
	fmt.Println(cCyan("   │") + cBold(" Compare Footprint Lists (c)Frederic PONT ") + cCyan("│"))
	fmt.Println(cCyan("   │") + "     v20250526 - Free Software GNU GPL    " + cCyan("│"))
	fmt.Println(cCyan("   └──────────────────────────────────────────┘"))
	fmt.Println()
}

// PrintSuccess prints a green success message with a check mark.
func PrintSuccess(msg string) {
	fmt.Println(cGreen("✔ " + msg))
}

// PrintWarning prints a yellow warning message with a warning sign.
func PrintWarning(msg string) {
	fmt.Println(cYellow("⚠ " + msg))
}

// PrintInfo prints a cyan informational message.
func PrintInfo(msg string) {
	fmt.Println(cCyan("ℹ " + msg))
}

// PrintError prints a red error message with a cross.
func PrintError(msg string) {
	fmt.Println(cRed("✘ " + msg))
}

// PrintBold prints a bold message.
func PrintBold(msg string) {
	fmt.Println(cBold(msg))
}

// PrintKeyVal prints a labelled value: "key: value"
func PrintKeyVal(key, val string) {
	fmt.Printf("%s %s\n", cMagenta(key+":"), val)
}

// PrintSection prints a separator then a bold section title.
func PrintSection(title string) {
	PrintSeparator()
	fmt.Println(cBold(title))
}

// PrintPathTrim prints the trimmed path pair during path-trim interactive mode.
func PrintPathTrim(labelA, pathA string, trimA int, labelB, pathB string, trimB int) {
	PrintSeparator()
	fmt.Printf("%s %s  %s %d\n",
		cRed(labelA+":"), pathA,
		cMagenta("(dirs trimmed:)"), trimA,
	)
	fmt.Printf("%s %s  %s %d\n",
		cRed(labelB+":"), pathB,
		cMagenta("(dirs trimmed:)"), trimB,
	)
}

// PrintDuplicateAlert prints the duplicate files warning block.
func PrintDuplicateAlert(nbDupA, nbDupB int) {
	if nbDupA == 0 && nbDupB == 0 {
		return
	}
	PrintSeparator()
	PrintWarning(fmt.Sprintf("%d duplicate file(s) in list_A", nbDupA))
	PrintWarning(fmt.Sprintf("%d duplicate file(s) in list_B", nbDupB))
	PrintInfo("It is recommended to use the -p path option when duplicates exist.")
	PrintSeparator()
}

// PrintSummary prints the final comparison summary block.
func PrintSummary(nbA, nbB, diffCount, commonCount int, listA, listB string) {
	PrintSeparator()
	PrintBold("Summary")
	PrintSeparator()
	PrintKeyVal("Files in "+listA, fmt.Sprintf("%d", nbA))
	PrintKeyVal("Files in "+listB, fmt.Sprintf("%d", nbB))

	if nbA == nbB {
		PrintSuccess("Same number of files in both lists")
	} else {
		PrintWarning("Different number of files between lists")
	}

	if diffCount == 0 {
		PrintSuccess("No differences found")
	} else {
		PrintWarning(fmt.Sprintf("%d difference(s) found  →  results/diff.csv", diffCount))
	}

	if commonCount > 0 {
		PrintSuccess(fmt.Sprintf("%d common file(s) found  →  results/common.csv", commonCount))
	} else {
		PrintInfo("No common files found")
	}
	PrintSeparator()
}
