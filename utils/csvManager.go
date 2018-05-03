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
	"os"
	"encoding/csv"
	"io"
	"strconv"
)

const (
	DELIVERY_TIME_FILENAME = "deliveryTime_ist2.csv"
	DISTANCE_MAT_FILENAME  = "distanceMatrix_ist2.csv"
)

var (
	DeliveryTimeFilename   = DELIVERY_TIME_FILENAME
	DistanceMatrixFilename = DISTANCE_MAT_FILENAME
)

func saveToFile(file *os.File, orders [][]string, orderIndexToName, moverIndexToName map[int]string) int {

	w := csv.NewWriter(file)

	nCol := len(orders[0])
	defer w.Flush()

	for index, v := range orders {

		if index < nCol {
			v = append(v, orderIndexToName[index])
		} else {
			v = append(v, moverIndexToName[index-nCol])
		}

		err := w.Write(append(v[len(v)-1:], v[0:len(v)-1]...))
		checkError("Error in write.", err)
	}

	return 0
}

func saveHeaderToFile(file *os.File, orders []string) int {

	w := csv.NewWriter(file)

	defer w.Flush()

	err := w.Write(orders)
	checkError("Error in write.", err)

	return 1
}

func loadFromFile(file *os.File) [][]string {

	r := csv.NewReader(file)

	var read [][]string

	i := 0
	for {
		res, err := r.Read()
		if err == io.EOF {
			break
		}
		checkError("Error reading file.", err)
		read = append(read, res)
		i++
		//fmt.Printf("%s\n", res[0])
	}
	return read
}

func openFileToWrite(filename string) *os.File {
	var file *os.File
	var err error

	if Exist(filename) {
		deleteFile(filename)
	}

	file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0755)
	checkError("Error in opening file.", err)

	return file
}

func openFileToRead(filename string) *os.File {
	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_RDONLY, 0755)
	checkError("Error in opening file ", err)

	return file
}

func closeFile(file *os.File) {
	err := file.Close()
	checkError("Error in closing file.", err)
}

func readInputFile(filename string) [][]string {

	file := openFileToRead(filename)
	orders := loadFromFile(file)
	closeFile(file)
	return orders
}

/* helps maintaining (index, ID) of order */
type S struct {
	I int
	N string
}

/**
 * map with key = (order_index, order_name)
 * 			value = slice with distance from current order
 * 					to every other order
 */
type DistMatrix map[S][]int
type TargetTimeVector map[string]int

/**
 * n = # orders
 */
func ReadDistanceMatrix(n int) (ordOrd DistMatrix, movOrd DistMatrix) {

	ordOrd = make(DistMatrix)
	movOrd = make(DistMatrix)
	mat := readInputFile(DistanceMatrixFilename)

	for i, v := range mat {
		if i == 0 {
			continue
		}
		if i <= n { /* order */
			ordOrd[S{i - 1, v[0]}] = strArrToIntArr(v[1:])
		} else {
			movOrd[S{i - n - 1, v[0]}] = strArrToIntArr(v[1:])
		}
	}

	return ordOrd, movOrd
}

func ReadOrdersTargetTime() TargetTimeVector {

	data := make(TargetTimeVector)
	delTime := readInputFile(DeliveryTimeFilename)

	for i, v := range delTime {
		if i == 0 {
			continue
		}
		data[v[0]] = strToInt(v[1])
	}
	return data
}

/**
 * Write on files the output
 */
func WriteAdjMatOnFile(filename string, y [][]uint8, orderIndexToName, moverIndexToName map[int]string) {

	nRow := len(y)
	nCol := len(y[0])

	colHeader := make([]string, nCol+1)
	colHeader[0] = "Y"
	for i := 1; i < nCol+1; i++ {
		colHeader[i] = orderIndexToName[i-1]
	}

	adjMat := ConvertUint8ToStringMatrix(y, nRow, nCol)
	file := openFileToWrite(filename)
	saveHeaderToFile(file, colHeader)
	saveToFile(file, adjMat, orderIndexToName, moverIndexToName)

	closeFile(file)
}

func WriteOrderVectorInt(filename string, x []int, orderIndexToName, moverIndexToName map[int]string, nOrder int, header []string) {

	file := openFileToWrite(filename)
	w := csv.NewWriter(file)

	err := w.Write(header)
	checkError("Error in write.", err)

	for index, value := range x {

		if index < nOrder {
			err := w.Write([]string{orderIndexToName[index], strconv.Itoa(value)})
			checkError("Error in write.", err)

		} else {
			err := w.Write([]string{moverIndexToName[index-nOrder], strconv.Itoa(value)})
			checkError("Error in write.", err)

		}
	}

	w.Flush()
	closeFile(file)
}

func WriteOrderVectorUint8(filename string, x []uint8, orderIndexToName map[int]string, header []string) {

	file := openFileToWrite(filename)
	w := csv.NewWriter(file)

	err := w.Write(header)
	checkError("Error in write.", err)
	for index, value := range x {
		err := w.Write([]string{orderIndexToName[index], strconv.Itoa(int(value))})
		checkError("Error in write.", err)

	}

	w.Flush()

	closeFile(file)
}

func WriteResultsTable(filename string, res [][]int, header []string) {
	file := openFileToWrite(filename)
	w := csv.NewWriter(file)

	err := w.Write(header)
	checkError("Error in write.", err)
	for _, val := range res {
		err := w.Write([]string{
			strconv.Itoa(val[0]),
			strconv.Itoa(val[1]),
			strconv.Itoa(val[2]),
			strconv.Itoa(val[3]),
			strconv.Itoa(val[4]),
			strconv.Itoa(val[5]),
			strconv.Itoa(val[6]),
			strconv.Itoa(val[7]),
			strconv.Itoa(val[8]),
			})
		checkError("Error in write.", err)

	}

	w.Flush()

	closeFile(file)
}

