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

package fileutil

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// ReadLists read 2 lists of footprints and compare them
func ReadLists() {

	args := Config.Arg
	// Create a log file and close it
	logFile, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	logFile.Close()

	la, lb := GetLists()
	// Channels to receive results from tasks
	ch1 := make(chan map[string]string)
	ch2 := make(chan map[string]string)
	cha := make(chan int)
	chb := make(chan int)
	chDupA := make(chan int)
	chDupB := make(chan int)

	// Launch tasks as goroutines
	go processOneList(la, "list_A/", ch1, cha, chDupA, args)
	go processOneList(lb, "list_B/", ch2, chb, chDupB, args)

	mapA := <-ch1
	mapB := <-ch2
	nbFilaA := <-cha
	nbFilaB := <-chb
	nbDuplicateA := <-chDupA
	nbDuplicateB := <-chDupB

	logFile, err = os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	writeLogSTDout(strconv.Itoa(nbFilaA)+" files in "+la, logFile)
	writeLogSTDout(strconv.Itoa(nbFilaB)+" files in "+lb, logFile)

	// fmt.Println(mapA)
	// fmt.Println(mapB)
	fmt.Println(Config)

	diffCount, diff := compareMaps(mapA, mapB)

	switch diffCount {
	case 0:
		if nbFilaA == nbFilaB {
			writeLogSTDout("The number of files is the same !", logFile)
			writeLogSTDout("no differences found !", logFile)
			duplicateAlert(nbDuplicateA, nbDuplicateB)
		} else {
			writeLogSTDout("The number of files is not the same !", logFile)
			writeLogSTDout("no differences found in common files !", logFile)
		}

	default:
		if nbFilaA != nbFilaB {
			writeLogSTDout("The number of files is not the same !", logFile)
		}

		writeLogSTDout(strconv.Itoa(diffCount)+" differences found !", logFile)
		printCSV(diff, "diff.csv")
	}

}

// duplicateAlert print message if there are duplicate files
func duplicateAlert(nbDuplicateA, nbDuplicateB int) {
	if nbDuplicateA != 0 || nbDuplicateB != 0 {
		fmt.Println("\033[31m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
		fmt.Println("There are", nbDuplicateA, "duplicate files in list_A")
		fmt.Println("There are", nbDuplicateB, "duplicate files in list_B")
		fmt.Println("It is recommended to do a comparison by path instead \nof files using the -p option to compare footprints.")
		fmt.Println("\033[31m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	}
}

// processOneList read one list of footprints and build a map of file name to footprint
func processOneList(list string, Listdir string, ch chan map[string]string, ch_ct, chDup chan int, args Args) {
	var data [][]string
	if Config.TrimPath {
		data = ReadOneListAndTrimPath(Listdir+list, args)
	} else {
		data = ReadOneList(Listdir+list, args)
	}

	dataMap, duplicates := strSliceToMap(data) // dataMap is a map of file name to footprint

	// Ensure to send a default value if data is empty
	if len(data) == 0 {
		ch <- make(map[string]string) // Send an empty map if no data
		ch_ct <- 0
	} else {
		ch <- dataMap
		ch_ct <- len(data)
	}
	chDup <- duplicates
}

// ReadOneList read one list of footprints and return a slice of string
func ReadOneList(path string, args Args) [][]string {
	var rows [][]string

	// comparison criteria = map key for signature comparison, can be filename (row 1) or path (row 2)
	rowIndex := 1
	if args.ComparisonCriteria == "path" {
		rowIndex = 2
	}

	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new csv.Reader
	reader := csv.NewReader(file)
	// Set the delimiter to TAB
	reader.Comma = '\t'
	// Set the comment character to '#'
	reader.Comment = '#'
	// Set the number of fields per record to 2, ie footprint in the first column and file name in the second column
	//reader.FieldsPerRecord = 2
	// Loop through the remaining lines
	for {
		// Read a line
		line, err := reader.Read()
		// Check the error value
		if err != nil {
			// Break the loop when the end of the file is reached
			if err == io.EOF {
				break
			}
			// Print the error otherwise
			fmt.Println(err)
			return [][]string{}
		}

		// Append the value to allPath
		// line[0]=footprint, line[1]=filename

		rows = append(rows, []string{line[0], line[rowIndex]})
	}

	return rows
}

// ReadOneListAndTrimPath read one list of footprints and return a 2 columns slice of string with the signature and the trimed path
func ReadOneListAndTrimPath(path string, args Args) [][]string {
	fmt.Println("\033[34m━━━━━━━━━━━ reconstructPathByIndex ━━━━━━━━━━━\033[0m")

	trimIndex := 0
	if strings.Contains(path, "list_A") {
		trimIndex = Config.TrimIndexPathA
	} else {
		trimIndex = Config.TrimIndexPathB
	}
	fmt.Println("trimIndex", trimIndex)
	var rows [][]string

	// comparison criteria = map key for signature comparison, can be filename (row 1) or path (row 2)
	rowIndex := 1
	if args.ComparisonCriteria == "path" {
		rowIndex = 2
	}

	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new csv.Reader
	reader := csv.NewReader(file)
	// Set the delimiter to TAB
	reader.Comma = '\t'
	// Set the comment character to '#'
	reader.Comment = '#'
	// Set the number of fields per record to 2, ie footprint in the first column and file name in the second column
	//reader.FieldsPerRecord = 2
	// Loop through the remaining lines
	for {
		// Read a line
		line, err := reader.Read()
		// Check the error value
		if err != nil {
			// Break the loop when the end of the file is reached
			if err == io.EOF {
				break
			}
			// Print the error otherwise
			fmt.Println(err)
			return [][]string{}
		}

		// Append the value to allPath
		// line[0]=footprint, line[1]=filename

		rows = append(rows, []string{line[0], ReconstructPathByIndex(removeLeadingSlash(line[rowIndex]), trimIndex, Config.CommonDirSep)})
	}

	return rows
}

// strSliceToMap convert a slice of string to a map of file name to footprint
func strSliceToMap(slice [][]string) (map[string]string, int) {
	duplicates := 0
	// Open the log file for appending (create if it doesn't exist)
	logFile, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	// filename/path => footprint map
	fpMap := make(map[string]string, len(slice))

	for _, row := range slice {
		if Haskey(fpMap, row[1]) {
			writeLogSTDout("Duplicate file ! "+row[1], logFile)
			duplicates++
		}
		fpMap[row[1]] = row[0]
	}
	//fmt.Println("duplicates", duplicates)
	return fpMap, duplicates
}

// compareMaps compare 2 maps of file name to footprint
func compareMaps(map1 map[string]string, map2 map[string]string) (int, [][]string) {
	//var differences [][]string
	differences := [][]string{{"file", "footprint A", "footprint B"}}
	diffCount := 0
	for key, value := range map1 {
		if val, ok := map2[key]; !ok || val != value {
			differences = append(differences, []string{key, value, val})
			diffCount++
		}
	}

	return diffCount, differences
}

// writeLogSTDout write message to  Combine stdout and log file writers
func writeLogSTDout(message string, logFile *os.File) {
	fmt.Println(message)
	log.SetOutput(logFile)
	log.Println(message)
}

// Haskey test if item is in map
func Haskey(myMap map[string]string, key string) bool {
	// Check if key exists
	_, ok := myMap[key]

	return ok
}

// ReadTsvHead read the header of a TSV file
func ReadTsvHead(path string) [][]string {
	fmt.Println("\033[34m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println("ReadT head", path)
	records := [][]string{}
	// Open the CSV file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		records = append(records, fields)
		fmt.Println("\033[33m" + strings.Join(fields, "\t") + "\033[0m")
		lineCount++
		if lineCount == 3 {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return records
}
