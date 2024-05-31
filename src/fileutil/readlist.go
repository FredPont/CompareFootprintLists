package fileutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func ReadLists() {
	// Create a log file
	logFile, err := os.Create("output.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	la, lb := GetLists()
	// Channels to receive results from tasks
	ch1 := make(chan map[string]string)
	ch2 := make(chan map[string]string)
	cha := make(chan int)
	chb := make(chan int)

	// Launch tasks as goroutines
	go processOneList(la, "list_A/", ch1, cha)
	go processOneList(lb, "list_B/", ch2, chb)

	mapA := <-ch1
	mapB := <-ch2
	nbFilaA := <-cha
	nbFilaB := <-chb
	//fmt.Println(nbFilaA, " files in ", la)
	//fmt.Println(nbFilaB, " files in ", lb)
	writeLogSTDout(strconv.Itoa(nbFilaA)+" files in "+la, logFile)
	writeLogSTDout(strconv.Itoa(nbFilaB)+" files in "+lb, logFile)

	diffCount, diff := compareMaps(mapA, mapB)

	switch diffCount {
	case 0:
		if nbFilaA == nbFilaB {
			//fmt.Println("The files are identical !")
			writeLogSTDout("The files are identical !", logFile)
		} else {
			//fmt.Println("The number of files is not the same !")
			writeLogSTDout("The number of files is not the same !", logFile)
			//fmt.Println("no differences found !")
			writeLogSTDout("no differences found !", logFile)
		}

	default:
		if nbFilaA != nbFilaB {
			//fmt.Println("The number of files is not the same !")
			writeLogSTDout("The number of files is not the same !", logFile)
		}
		//fmt.Println(diffCount, "differences found !")
		writeLogSTDout(strconv.Itoa(diffCount)+" differences found !", logFile)
		printCSV(diff, "diff.csv")
	}

}

func processOneList(list string, Listdir string, ch chan map[string]string, ch_ct chan int) {
	data := ReadOneList(Listdir + list)
	dataMap := strSliceToMap(data)
	ch <- dataMap
	ch_ct <- len(dataMap)
}

func ReadOneList(path string) [][]string {
	var rows [][]string
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
		rows = append(rows, []string{line[0], line[1]})
	}

	return rows
}

func strSliceToMap(slice [][]string) map[string]string {
	// filename => footprint map
	fpMap := make(map[string]string, len(slice))

	for _, row := range slice {
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

func writeLogSTDout(message string, logFile *os.File) {
	// Combine stdout and log file writers
	writer := io.MultiWriter(os.Stdout, logFile)
	// Write to the combined writer
	log.SetOutput(writer)
	log.Println(message)
}
