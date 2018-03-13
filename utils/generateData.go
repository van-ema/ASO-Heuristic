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
