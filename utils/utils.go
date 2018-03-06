package utils

import (
	"strconv"
	"log"
	"fmt"
)


const Inf = int(^uint(0) >> 1)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func strArrToIntArr(v []string) []int {
	l := len(v)
	intArr := make([]int,l)
	for i, val := range v {
		intArr[i] = strToInt(val)
	}
	return intArr
}

func strToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return i
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

/**
 * n = # of orders
 */
func PrintDistanceMatrix(mat [][]int, n int){

	fmt.Printf("        |")

	/*for k := 0; k < n; k++ {
		fmt.Printf(" Order %d|", k)
	} */

	for i := 0; i < len(mat) ; i ++ {

		if i < n {
			fmt.Printf("\nOrder %d |", i)
		} else {
			fmt.Printf("\nMover %d |", i - n)
		}
		for j := 0; j < n ; j++ {
			fmt.Printf("  Order %d:  %d    ", j, mat[i][j])

		}
	}

	fmt.Printf("\n")
	fmt.Printf("\n")
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