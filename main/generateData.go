package main

import (
	"math/rand"
	"time"
)

const (
	ORDER_N      = 300 /* #Orders */
	MOVER_N      = 30  /* #Movers */
	MAX_DISTANCE = 200 /* Max distance between two orders*/
	STD          = 30  /* standard deviation for delivery times*/
	MEAN         = 90  /* Mean for delivery times */
)

func CreateOrderMatrix(nOrder, nMover int) [][]int {
	rand.Seed(time.Now().Unix())

	m1 := make([][]int, nOrder+nMover)
	for i := range m1 {
		m1[i] = make([]int, nOrder)
	}

	for i := 0; i < nOrder+nMover; i++ {
		for j := 0; j < nOrder; j++ {

			switch {
				case i==j :
					m1[i][j] = 0
				default :
					m1[i][j] = NextRandom(0, MAX_DISTANCE)
			}
		}
	}

	return m1
}

func CreateDeliveryTimeVector(nOrder int) []int {
	v := make([]int, nOrder)
	for i := 0; i < nOrder; i++ {
		v[i] = int(rand.NormFloat64()*STD + MEAN)
	}

	return v
}

func NextRandom(min, max int) int {
	return rand.Intn(max-min) + min
}

//func main() {
//	fmt.Print(CreateOrderMatrix(ORDER_N, MOVER_N))
//	fmt.Println("\n")
//	fmt.Println(CreateDeliveryTimeVector(ORDER_N))
//}
