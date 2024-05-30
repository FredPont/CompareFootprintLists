package fileutil

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

func ReadLists() {
	la, lb := GetLists()
	dataA := ReadOneList("list_A/" + la)
	dataB := ReadOneList("list_B/" + lb)
	mapA := strSliceToMap(dataA)
	mapB := strSliceToMap(dataB)
	//fmt.Println(mapA, mapB)
	diffCount, diff := compareMaps(mapA, mapB)

	switch diffCount {
	case 0:
		fmt.Println("The files are identical found !")
	default:
		fmt.Println(diffCount, "differences found !")
		printCSV(diff, "diff.csv")
	}

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
