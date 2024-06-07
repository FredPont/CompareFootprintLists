package fileutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

type Args struct {
	ComparisonCriteria string // comparison by filename or path (to compare files with same file names)
}

func ReadLists(args Args) {
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

	// Launch tasks as goroutines
	go processOneList(la, "list_A/", ch1, cha, args)
	go processOneList(lb, "list_B/", ch2, chb, args)

	mapA := <-ch1
	mapB := <-ch2
	nbFilaA := <-cha
	nbFilaB := <-chb

	logFile, err = os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	writeLogSTDout(strconv.Itoa(nbFilaA)+" files in "+la, logFile)
	writeLogSTDout(strconv.Itoa(nbFilaB)+" files in "+lb, logFile)

	diffCount, diff := compareMaps(mapA, mapB)

	switch diffCount {
	case 0:
		if nbFilaA == nbFilaB {

			writeLogSTDout("The files are identical !", logFile)
		} else {

			writeLogSTDout("The number of files is not the same !", logFile)

			writeLogSTDout("no differences found !", logFile)
		}

	default:
		if nbFilaA != nbFilaB {

			writeLogSTDout("The number of files is not the same !", logFile)
		}

		writeLogSTDout(strconv.Itoa(diffCount)+" differences found !", logFile)
		printCSV(diff, "diff.csv")
	}

}

func processOneList(list string, Listdir string, ch chan map[string]string, ch_ct chan int, args Args) {
	data := ReadOneList(Listdir+list, args)
	dataMap := strSliceToMap(data)
	ch <- dataMap
	ch_ct <- len(dataMap)
}

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

func strSliceToMap(slice [][]string) map[string]string {
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
		}
		fpMap[row[1]] = row[0]
	}
	return fpMap
}

func compareMaps(map1 map[string]string, map2 map[string]string) (int, [][]string) {
	var differences [][]string
	diffCount := 0
	for key, value := range map1 {
		if val, ok := map2[key]; !ok || val != value {
			differences = append(differences, []string{key, val, value})
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
