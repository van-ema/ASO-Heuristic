package utils

import (
	"fmt"
	"os"
	"encoding/csv"
	"io"
)

const (
	deliveryTimeFilename = "input/deliveryTime_ist3.csv"
	distanceMatrixFilename = "input/distanceMatrix_ist3.csv")


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
		fmt.Printf("%s\n",res[0])
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

func readInputFile(filename string) [][]string {

	file := openFileToRead(filename)
	orders := loadFromFile(file)
	closeFile(file)
	return orders
}

/* helps maintaining (index, ID) of order */
type S struct {
	i int
	n string
}

/**
 * map with key = (order_index, order_name)
 * 			value = slice with distance from current order
 * 					to every other order
 */
type DistMatrix map[S][]int


/**
 * n = # orders
 */
func ReadDistanceMatrix(n int) (ordOrd DistMatrix, movOrd DistMatrix) {

	ordOrd = make(DistMatrix)
	movOrd = make(DistMatrix)
	mat := readInputFile(distanceMatrixFilename)

	for i, v := range mat {
		if i == 0 { continue }
		if i <= n { 	/* order */
			ordOrd[S{i - 1, v[0]}] = strArrToIntArr(v[1:])
		} else {
			movOrd[S{i - n - 1, v[0]}] = strArrToIntArr(v[1:])
		}
	}

	for k, v := range movOrd {
		fmt.Printf("Key: %d,%s Value: %d\n ", k.i, k.n, v[0])
	}

	return ordOrd, movOrd

}

func ReadOrdersTargetTime() map[string]int {

	data := make(map[string]int)
	delTime := readInputFile(deliveryTimeFilename)

	for i, v := range delTime {
		if i == 0 { continue }
		data[v[0]] = strToInt(v[1])
	}

	return data
	
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