package util

import (
	"encoding/csv"
	"strings"
	"log"
	"os"
)

// SeperatorStr is unique enough...maybe..
var SeperatorStr = "|$\t$|"

// ReadCsvData returns the data read from CSV file
func ReadCsvData(file string) (data [][]string, err error){
	f, err := os.Open(file)
    if err != nil {
        return nil, err
    }
    defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return lines, nil
}

// WriteCsvToFile write the data to local file
func WriteCsvToFile(file string, data [][]string) error {
	f, err := os.Create(file)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	for _, item := range data {
		err = csvWriter.Write(item)
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		log.Println("Flush() fail")
		return err
	}

	return nil
}
// RemoveArrayItem will remove the array[index]
func RemoveArrayItem(array []string, index int) (result []string) {
	result = append(array[:index], array[index+1:]...)

	return result
}

// Index returns the first index of the target string str, or -1 if no match is found.
func Index(arr []string, str string) int {
    for i, item := range arr {
        if item == str {
            return i
        }
    }
    return -1
}

// Include returns true if the target string str is in the slice.
func Include(arr []string, str string) bool {
    return Index(arr, str) >= 0
}

// Filter returns a new slice containing all strings in the slice that satisfy the predicate f.
func Filter(arr []string, f func(string) bool) []string {
    result := make([]string, 0)
    for _, item := range arr {
        if f(item) {
            result = append(result, item)
        }
    }
    return result
}

// Map returns a new slice containing the results of applying the function f to each string in the original slice.
func Map(arr []string, f func(string) string) []string {
    result := make([]string, len(arr))
    for i, item := range arr {
        result[i] = f(item)
    }
    return result
}

// ArrayTo2DSlice returns the 2D slice transformed from 1D
func ArrayTo2DSlice(data []string, seperator string) (result [][]string){
	result = make([][]string, len(data))
    for i, row := range data {
		result[i] = strings.Split(row, seperator)
	}
    return result
}

// TwoDimensionTo1DArray returns the 1D slice transformed from 2D
func TwoDimensionTo1DArray(data [][]string, seperator string) (result []string){
	result = make([]string, len(data))
    for i, row := range data {
        result[i] = strings.Join(row, seperator)
    }
    return result
}
