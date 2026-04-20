package fileutil

import (
	"fmt"
	"os"
)

// GetLists returns the first filename found in list_A/ and list_B/.
func GetLists() (string, string) {
	fa := readDir(dirListA)
	fb := readDir(dirListB)
	PrintKeyVal("List A", fa[0])
	PrintKeyVal("List B", fb[0])
	return fa[0], fb[0]
}

// readDir returns the names of all files in dirPath.
func readDir(dirPath string) []string {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		PrintError(fmt.Sprintf("Cannot read directory %s: %v", dirPath, err))
		return []string{}
	}

	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names
}
