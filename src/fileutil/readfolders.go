package fileutil

import (
	"fmt"
	"os"
)

func GetLists() (string, string) {
	fa := Readdir("list_A")
	fb := Readdir("list_B")

	fmt.Println(fa, fb)
	// return 2 lists
	return fa[0], fb[0]
}

func Readdir(dirPath string) []string {
	var filenames []string
	// Read the directory contents
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []string{}
	}

	for _, v := range files {
		//fmt.Println(v.Name(), v.IsDir())
		filenames = append(filenames, v.Name())
	}
	return filenames

}
