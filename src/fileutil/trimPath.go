/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.

Written by Frederic PONT.
(c) Frederic Pont 2025
*/
package fileutil

import (
	"fmt"
	"strings"
)

// ###########################################

func TrimPath(path1, path2 string) {
	// path1 := "/p1/p2/common/p3"
	// path2 := "/d1/common/d2"

	// Process both paths
	//path1 = removeLeadingSlash(path1)
	//path2 = removeLeadingSlash(path2)
	newPath1, newPath2 := RemoveLeadingDirs(path1, path2)

	// Print the results
	fmt.Println(newPath1) // Output: common/p3
	fmt.Println(newPath2) // Output: common/d2
}

// Function to remove leading directories before the common directory
func RemoveLeadingDirs(path1, path2 string) (string, string) {
	//path1 = removeLeadingSlash(path1)
	//path2 = removeLeadingSlash(path2)
	// Get the separator from the first path. The other path should have the same separator
	// if the same operating system is used. If not, the separator will be different.
	// and sep1 will be used to reconstruct the path.
	sep1 := GetSeparator(path1)
	sep2 := GetSeparator(path2)
	Config.CommonDirSep = sep1
	// Helper function to check if a slice contains a specific element
	contains := func(slice []string, item string) bool {
		for _, v := range slice {
			if v == item {
				return true
			}
		}
		return false
	}

	// Helper function to reconstruct the path from the common directory
	reconstructPath := func(segments []string, commonDir, sep string) string {
		for i, seg := range segments {
			if seg == commonDir {
				return strings.Join(segments[i:], sep)
			}
		}
		return "" // Return empty if commonDir is not found
	}
	// Split the paths into segments
	segments1 := strings.Split(path1, sep1)
	segments2 := strings.Split(path2, sep2)

	// Find the common directory in path1
	var commonDir1 string
	for i := 0; i < len(segments1); i++ {

		if contains(segments2, segments1[i]) {
			commonDir1 = segments1[i]
			//fmt.Println(segments1, i, segments1[i], "==>", commonDir1)
			Config.TrimIndexPathA = i
			break
		}

		//fmt.Println(segments1, i, segments1[i], "==>", commonDir1, Config)
	}

	// Find the common directory in path2
	var commonDir2 string
	for i := 0; i < len(segments2); i++ {
		if contains(segments1, segments2[i]) {
			commonDir2 = segments2[i]
			Config.TrimIndexPathB = i
			break
		}

	}

	// Reconstruct the paths from the common directory onward
	newPath1 := reconstructPath(segments1, commonDir1, sep1)
	newPath2 := reconstructPath(segments2, commonDir2, sep1)

	return newPath1, newPath2
}

func GetSeparator(path string) string {
	if strings.Contains(path, "\\") {
		return "\\"
	}
	return "/"
}

// reconstructPathByIndex reconstruct the path from the start index directory
func ReconstructPathByIndex(path string, startIdx int, sep string) string {
	segments := strings.Split(path, sep)
	if startIdx < 0 || startIdx >= len(segments) {
		return "" // Return empty if index is out of range
	}
	newPath := strings.Join(segments[startIdx:], sep)
	//fmt.Println(newPath)
	return newPath
}

// Remove leading slash
func removeLeadingSlash(s string) string {
	if strings.HasPrefix(s, "/") {
		s = s[1:]
	}
	return s
}
