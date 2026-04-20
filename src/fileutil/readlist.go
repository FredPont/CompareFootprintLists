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

// FileRecord holds all three columns for a file entry.
type FileRecord struct {
	Footprint string
	Key       string // filename or path, depending on comparison mode
	Path      string // always the original path column (col 2)
}

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

	// Channels to receive results from goroutines
	ch1 := make(chan map[string]FileRecord)
	ch2 := make(chan map[string]FileRecord)
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

	fmt.Println(Config)

	diffCount, diff, commonCount, common := compareMaps(mapA, mapB)

	// Always write the common files table
	if commonCount > 0 {
		writeLogSTDout(strconv.Itoa(commonCount)+" common files found !", logFile)
		printCSV(common, "common.csv")
	} else {
		writeLogSTDout("no common files found !", logFile)
	}

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

// processOneList read one list of footprints and build a map of key → FileRecord
func processOneList(list string, Listdir string, ch chan map[string]FileRecord, ch_ct, chDup chan int, args Args) {
	var data [][]string
	if Config.TrimPath {
		data = ReadOneListAndTrimPath(Listdir+list, args)
	} else {
		data = ReadOneList(Listdir+list, args)
	}

	dataMap, duplicates := strSliceToMap(data)

	if len(data) == 0 {
		ch <- make(map[string]FileRecord)
		ch_ct <- 0
	} else {
		ch <- dataMap
		ch_ct <- len(data)
	}
	chDup <- duplicates
}

// ReadOneList read one list of footprints and return a 3-column slice:
// [footprint, key (filename or path), path]
func ReadOneList(path string, args Args) [][]string {
	var rows [][]string

	rowIndex := 1
	if args.ComparisonCriteria == "path" {
		rowIndex = 2
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.Comment = '#'

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return [][]string{}
		}

		pathCol := ""
		if len(line) > 2 {
			pathCol = line[2]
		}
		// [footprint, key, path]
		rows = append(rows, []string{line[0], line[rowIndex], pathCol})
	}

	return rows
}

// ReadOneListAndTrimPath read one list and trim the path column before building the key.
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

	rowIndex := 1
	if args.ComparisonCriteria == "path" {
		rowIndex = 2
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return [][]string{}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.Comment = '#'

	for {
		line, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return [][]string{}
		}

		pathCol := ""
		if len(line) > 2 {
			pathCol = line[2]
		}
		trimmedKey := ReconstructPathByIndex(removeLeadingSlash(line[rowIndex]), trimIndex, Config.CommonDirSep)
		// [footprint, trimmed key, original path]
		rows = append(rows, []string{line[0], trimmedKey, pathCol})
	}

	return rows
}

// strSliceToMap convert a 3-column slice to a map of key → FileRecord
func strSliceToMap(slice [][]string) (map[string]FileRecord, int) {
	duplicates := 0
	logFile, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()

	fpMap := make(map[string]FileRecord, len(slice))

	for _, row := range slice {
		key := row[1]
		pathCol := ""
		if len(row) > 2 {
			pathCol = row[2]
		}
		if Haskey(fpMap, key) {
			writeLogSTDout("Duplicate file ! "+key, logFile)
			duplicates++
		}
		fpMap[key] = FileRecord{
			Footprint: row[0],
			Key:       key,
			Path:      pathCol,
		}
	}
	return fpMap, duplicates
}

// compareMaps compare 2 maps and return differences AND common files.
// Returns: diffCount, diff rows, commonCount, common rows
func compareMaps(mapA, mapB map[string]FileRecord) (int, [][]string, int, [][]string) {
	differences := [][]string{{"file", "footprint_A", "footprint_B"}}
	common := [][]string{{"file", "footprint", "path_in_A", "path_in_B"}}
	diffCount := 0
	commonCount := 0

	for key, recA := range mapA {
		if recB, ok := mapB[key]; !ok {
			// present in A, missing in B
			differences = append(differences, []string{key, recA.Footprint, ""})
			diffCount++
		} else if recA.Footprint != recB.Footprint {
			// present in both but different signature
			differences = append(differences, []string{key, recA.Footprint, recB.Footprint})
			diffCount++
		} else {
			// identical signature → common
			common = append(common, []string{key, recA.Footprint, recA.Path, recB.Path})
			commonCount++
		}
	}

	// Files present in B but missing in A
	for key, recB := range mapB {
		if _, ok := mapA[key]; !ok {
			differences = append(differences, []string{key, "", recB.Footprint})
			diffCount++
		}
	}

	return diffCount, differences, commonCount, common
}

// writeLogSTDout write message to stdout and log file
func writeLogSTDout(message string, logFile *os.File) {
	fmt.Println(message)
	log.SetOutput(logFile)
	log.Println(message)
}

// Haskey test if key is in map
func Haskey(myMap map[string]FileRecord, key string) bool {
	_, ok := myMap[key]
	return ok
}

// ReadTsvHead read the first 3 lines of a TSV file
func ReadTsvHead(path string) [][]string {
	fmt.Println("\033[34m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\033[0m")
	fmt.Println("ReadT head", path)
	records := [][]string{}
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
