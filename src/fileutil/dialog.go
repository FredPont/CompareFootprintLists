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
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AksUser ask user if he wants to trim the path of the footprints
func AksUser() {

	arg := Config.Arg
	fmt.Println(Config, "AksUser", arg)
	if arg.ComparisonCriteria == "path" {
		fmt.Println("Do you want to trim the path of the footprints ? y/n")
		// read user input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if input == "y\n" {
			Config.TrimPath = true
			read3lines()
		} else {
			Config.TrimPath = false
		}
	}

}

// read3lines read the first 3 lines of the 2 lists
func read3lines() {

	la, lb := GetLists()
	HeadA := ReadTsvHead("list_A/" + la)
	HeadB := ReadTsvHead("list_B/" + lb)
	path1, path2 := HeadA[0][2], HeadB[0][2]
	// newPath1, newPath2 := RemoveLeadingDirs(path1, path2)

	// fmt.Println(newPath1)
	// fmt.Println(newPath2)
	path1 = removeLeadingSlash(path1)
	path2 = removeLeadingSlash(path2)
	showPathAutoTrim(path1, path2)
	dirToTrim(path1, path2)
	// showPath(path1, path2)
	// fmt.Println("Do you want to change the number of directory to trim ? y/n")
	// // read user input
	// reader := bufio.NewReader(os.Stdin)
	// input, _ := reader.ReadString('\n')
	// if input == "y\n" {
	// 	Config.TrimIndexPathA = askUser()
	// 	Config.TrimIndexPathB = askUser()
	// }
}

// dirToTrim ask user if he wants to change the number of directory to trim
func dirToTrim(path1, path2 string) {
	//showPath(path1, path2)
	fmt.Println("")
	fmt.Println("Do you want to change the number of directory to trim ? y/n")
	// read user input
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	if input == "y\n" {
		Config.TrimIndexPathA = askUserNumber("Number of directory to trim in list A :")

		Config.TrimIndexPathB = askUserNumber("Number of directory to trim in list B :")
		fmt.Println(Config)

		ReconstructPathByIndex(path1, Config.TrimIndexPathA, Config.CommonDirSep)
		ReconstructPathByIndex(path2, Config.TrimIndexPathB, Config.CommonDirSep)
		showPath(path1, path2)
		dirToTrim(path1, path2)
	}
}

// askUserNumber ask user for a number
func askUserNumber(message string) int {
	fmt.Println(message)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input) // Remove whitespace and newline
	n, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Conversion error:", err)
	}
	return n
}

// showPath show the new path and the number of directory to trim
func showPath(path1, path2 string) {
	newPath1 := ReconstructPathByIndex(path1, Config.TrimIndexPathA, Config.CommonDirSep)
	newPath2 := ReconstructPathByIndex(path2, Config.TrimIndexPathB, Config.CommonDirSep)
	//newPath1, newPath2 := RemoveLeadingDirs(path1, path2)
	fmt.Println("\033[31m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println("\033[31mnewPathA : \033[0m", newPath1, "\033[35mNumber of directory to trim :\033[0m", Config.TrimIndexPathA)
	fmt.Println("\033[31mnewPathB : \033[0m", newPath2, "\033[35mNumber of directory to trim :\033[0m", Config.TrimIndexPathB)

}

// showPath show the new path and the number of directory to trim
func showPathAutoTrim(path1, path2 string) {
	newPath1, newPath2 := RemoveLeadingDirs(path1, path2)
	fmt.Println("\033[31m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println("\033[31mnewPathA : \033[0m", newPath1, "\033[35mNumber of directory to trim :\033[0m", Config.TrimIndexPathA)
	fmt.Println("\033[31mnewPathB : \033[0m", newPath2, "\033[35mNumber of directory to trim :\033[0m", Config.TrimIndexPathB)

}
