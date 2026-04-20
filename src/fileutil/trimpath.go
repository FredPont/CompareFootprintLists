package fileutil

import (
	"strings"
)

// RemoveLeadingDirs auto-detects the first common directory segment between
// path1 and path2, stores the trim indices in Config, and returns the trimmed
// versions of both paths.
func RemoveLeadingDirs(path1, path2 string) (string, string) {
	sep1 := GetSeparator(path1)
	sep2 := GetSeparator(path2)
	Config.CommonDirSep = sep1

	segments1 := strings.Split(path1, sep1)
	segments2 := strings.Split(path2, sep2)

	// Build a lookup set from segments2 for O(1) membership test
	set2 := make(map[string]bool, len(segments2))
	for _, s := range segments2 {
		set2[s] = true
	}

	// Find the first segment of path1 that appears in path2
	for i, seg := range segments1 {
		if set2[seg] {
			Config.TrimIndexPathA = i
			break
		}
	}

	// Build a lookup set from segments1 for O(1) membership test
	set1 := make(map[string]bool, len(segments1))
	for _, s := range segments1 {
		set1[s] = true
	}

	// Find the first segment of path2 that appears in path1
	for i, seg := range segments2 {
		if set1[seg] {
			Config.TrimIndexPathB = i
			break
		}
	}

	newPath1 := ReconstructPathByIndex(path1, Config.TrimIndexPathA, sep1)
	newPath2 := ReconstructPathByIndex(path2, Config.TrimIndexPathB, sep1)

	return newPath1, newPath2
}

// ReconstructPathByIndex rebuilds a path by dropping the first startIdx
// directory segments. Returns "" if startIdx is out of range.
func ReconstructPathByIndex(path string, startIdx int, sep string) string {
	if sep == "" {
		sep = GetSeparator(path)
	}
	segments := strings.Split(path, sep)
	if startIdx < 0 || startIdx >= len(segments) {
		return ""
	}
	return strings.Join(segments[startIdx:], sep)
}

// GetSeparator returns the path separator used in the given path string.
func GetSeparator(path string) string {
	if strings.Contains(path, "\\") {
		return "\\"
	}
	return "/"
}

// removeLeadingSlash strips a single leading "/" if present.
func removeLeadingSlash(s string) string {
	return strings.TrimPrefix(s, "/")
}
