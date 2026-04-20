/*
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

Written by Frederic PONT.
(c) Frederic Pont 2024
*/

package main

import (
	"CompareFootprintLists/src/fileutil"
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	fileutil.PrintTitle()

	fileutil.Config.Arg = parseArgs()

	// Ask user for path-trim preferences BEFORE launching goroutines,
	// so that Config is fully populated when the workers read it.
	fileutil.AskUser()

	t0 := time.Now()

	fileutil.ReadLists()

	fileutil.PrintInfo(fmt.Sprintf("Elapsed time: %s", time.Since(t0)))

	fmt.Println()
	fileutil.PrintInfo("Press Enter to exit.")

	reader := bufio.NewReader(os.Stdin)
	if _, err := reader.ReadString('\n'); err != nil {
		fileutil.PrintError(fmt.Sprintf("Error reading input: %v", err))
	}
}

// parseArgs parses CLI flags and returns an Args struct.
func parseArgs() fileutil.Args {
	args := fileutil.Args{}
	flag.StringVar(&args.ComparisonCriteria, "p", "filename",
		`Comparison key: "filename" (default) or "path".
Use "path" when duplicate filenames exist across lists.
Caution: with "path", the full path must be identical in both lists.`)
	flag.Parse()
	return args
}
