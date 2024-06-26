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
 (c) Frederic Pont 2024
*/

package main

import (
	"CompareFootprintLists/src/fileutil"
	"flag"
	"fmt"
	"time"
)

func main() {
	fileutil.Title()
	args := parseARG()

	t0 := time.Now()

	fileutil.ReadLists(args)

	fmt.Println("\ndone !")
	fmt.Println("Elapsed time : ", time.Since(t0))

	// Define the duration for the countdown
	// Set the countdown time in seconds
	countdownFrom := 3
	fileutil.Timer(countdownFrom)
	//time.Sleep(3 * time.Second) // sleep to read results before windows close
}

// parse arg of the command line and return the argument struct
func parseARG() fileutil.Args {
	args := fileutil.Args{}
	flag.StringVar(&args.ComparisonCriteria, "p", "filename",
		`		Comparison by file names or path. filename or path.
		This option is useful if there are file duplicates whith the same name. 
		Caution, if path is used, the full file path must be the same in both lists`)
	flag.Parse()
	return args
}
