package main

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
)

func saveToFile(file *os.File, orders[][]string) int {

	w := csv.NewWriter(file)

	defer w.Flush()

	for _, v := range orders {

		err := w.Write(v)
		checkError("Error in write.", err)
	}

	return 0
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
	}
	return read
}

func openFileToWrite(filename string) *os.File {
	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_CREATE | os.O_WRONLY, 0755)
	checkError("Error in opening file.", err)

	return file
}

func openFileToRead(filename string) *os.File {
	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_RDONLY, 0755)
	checkError("Error in opening file", err)

	return file
}

func closeFile(file *os.File) {
	err := file.Close()
	checkError("Error in closing file.", err)
}

/**
 * Write in a csv file a generated OrderMatrix.
 */
func CreateInputFile(filename string) {

	data := CreateOrderMatrix(ORDER_N, MOVER_N)
	orders := ConvertIntToStringMatrix(data, ORDER_N+MOVER_N, ORDER_N)

	file := openFileToWrite(filename)

	saveToFile(file, orders)

	closeFile(file)
}

func ReadInputFile(filename string) [][]string {

	file := openFileToRead(filename)

	orders := loadFromFile(file)

	closeFile(file)

	fmt.Println(orders)

	return orders
}

/**
 * Write in a csv file the output OrderMatrix.
 */
//func CreateOutputFile(filename string) {
//
//	data :=
//	times := ConvertIntToStringMatrix(data, ..., ...)
//
//	file := openFileToWrite(filename)
//
//	saveToFile(file, times)
//
//	closeFile(file)
//}

func main() {

	//CreateInputFile("input.csv")

	//ReadInputFile("input.csv")

	fmt.Println("Yo\n")
}