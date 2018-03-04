package main

import (
	"strconv"
	"log"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

/**
 * Convert int values of generated OrderMatrix to string values
 * for csv package's methods.
 */
func ConvertIntToStringMatrix(from [][]int, nRow, nCol int) [][]string {

	var to = make([][]string, nRow)

	for i := 0; i < nRow;i++ {
		to[i] = make([]string, nCol)
		for j:=0; j < nCol; j++ {
			to[i][j] = strconv.Itoa(from[i][j])
		}
	}
	return to
}
