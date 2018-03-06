package utils

import (
	"strconv"
	"log"
)

const Inf = int(^uint(0) >> 1)

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

// return the minimum value and its index from an array  of int
func FindMin(array []int) (int, int) {
	min := Inf
	index := -1
	for i, v := range array {
		if v < min {
			min = v
			index = i
		}
	}

	return index, min
}