/*
 Copyright 2018 Ovidiu Daniel Barba, Laura Trivelloni, Emanuele Vannacci

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
*/

package utils

import (
	"strconv"
	"log"
	"fmt"
	"os"
)

const Inf = int(^uint(0) >> 1)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func strArrToIntArr(v []string) []int {
	l := len(v)
	intArr := make([]int, l)
	for i, val := range v {
		intArr[i] = strToInt(val)
	}
	return intArr
}

func strToInt(s string) int {
	i, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal(err)
		return -1
	}
	return int(i)
}

/**
 * Convert int values of generated OrderMatrix to string values
 * for csv package's methods.
 */
func ConvertIntToStringMatrix(from [][]int, nRow, nCol int) [][]string {

	var to = make([][]string, nRow)

	for i := 0; i < nRow; i++ {
		to[i] = make([]string, nCol)
		for j := 0; j < nCol; j++ {
			to[i][j] = strconv.Itoa(from[i][j])
		}
	}
	return to
}

func ConvertUint8ToStringMatrix(from [][]uint8, nRow, nCol int) [][]string {

	var to = make([][]string, nRow)

	for i := 0; i < nRow; i++ {
		to[i] = make([]string, nCol)
		for j := 0; j < nCol; j++ {
			to[i][j] = strconv.Itoa(int(from[i][j]))
		}
	}
	return to
}

/**
 * n = # of orders
 */
func PrintDistanceMatrix(mat [][]int, n int) {

	fmt.Printf("        |")

	/*for k := 0; k < n; k++ {
		fmt.Printf(" Order %d|", k)
	} */

	for i := 0; i < len(mat); i ++ {

		if i < n {
			fmt.Printf("\nOrder %d |", i)
		} else {
			fmt.Printf("\nMover %d |", i-n)
		}
		for j := 0; j < n; j++ {
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

func PrintAssigmentMatrix(mat [][]uint8, n int, ) {
	for i := 0; i < len(mat); i ++ {

		if i < n {
			fmt.Printf("\nOrder %d |", i)
		} else {
			fmt.Printf("\nMover %d |", i-n)
		}
		for j := 0; j < n; j++ {
			fmt.Printf("  Order %d:  %d    ", j, mat[i][j])

		}
	}

	fmt.Printf("\n")
	fmt.Printf("\n")
}

func PrintMatrix(mat [][]uint8) {
	for i := 0; i < len(mat); i++ {
		fmt.Print(mat[i])
		fmt.Print("\n")
	}
}

func Exist(f string) bool {
	_, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

func deleteFile(filename string) {
	// delete file
	var err = os.Remove(filename)
	if isError(err) {
		return
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}
