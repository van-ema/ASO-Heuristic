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

package main

import (
	"fmt"
)

const ES = 1000 /* end of shift */
const B = 10000 /* large number */

func Validate(result SolverResult, distances Distances, deliveryTime DeliveryTimeVector) bool {

	c1 := constraint1(result, distances, deliveryTime)
	c2 := constraint2(result, distances, deliveryTime)
	c3 := constraint3(result, distances, deliveryTime)
	c4 := constraint4(result, distances, deliveryTime)
	c5 := constraint5(result, distances, deliveryTime)
	c6 := constraint6(result, distances, deliveryTime)
	c7 := constraint7(result, distances, deliveryTime)
	c8 := constraint8(result, distances, deliveryTime)
	c9 := constraint9(result, distances, deliveryTime)
	c10 := constraint10(result, distances, deliveryTime)
	c11 := constraint11(result, distances, deliveryTime)
	c12 := constraint12(result, distances, deliveryTime)
	c13 := constraint13(result, distances, deliveryTime)
	return c1 && c2 && c3 && c4 && c5 && c6 && c7 && c8 && c9 && c10 && c11 && c12 && c13
}

func constraint1(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(1)
	}
	var sum uint8
	for i := 0; i < res.nOrder; i++ {
		sum = 0
		sum += res.w[i]
		for j := 0; j < res.nOrder+res.nMover; j++ {
			if i == j {
				continue
			}
			sum += res.y[j][i]

		}
		if sum != 1 {
			if DEBUG {
				printFailed()
			}
			return false
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint2(res SolverResult, dist Distances, del DeliveryTimeVector) bool {

	var sum uint8
	if DEBUG {
		printInitMess(2)
	}
	for i := 0; i < res.nMover+res.nOrder; i++ {
		sum = 0
		for j := 0; j < res.nOrder; j++ {
			if i == j {
				continue
			}
			sum += res.y[i][j]
			if sum > 1 {
				if DEBUG {
					printFailed()
				}
				return false
			}
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint3(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	//utils.PrintMatrix(UnfeasibleOrdersPairsMatrix)
	if DEBUG {
		printInitMess(3)
	}
	for i := 0; i < res.nOrder+res.nMover; i++ {
		for j := 0; j < res.nOrder; j++ {
			//fmt.Printf("y = %d, O = %d\n",res.y[i][j],UnfeasibleOrdersPairsMatrix[i][j])
			if res.y[i][j] == 1 && UnfeasibleOrdersPairsMatrix[i][j] == 1 {
				if DEBUG {
					printFailed()
				}
				return false
			} else {
				continue
			}
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint4(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(4)
	}
	for i := 0; i < res.nOrder; i++ {
		for j := 0; j < res.nOrder; j++ {
			if i == j {
				continue
			}
			if UnfeasibleOrdersPairsMatrix[i][j] == 0 { /* pair doesn't belong */
				if res.y[i][j] > 1-res.w[i] {
					if DEBUG {
						printFailed()
					}
					return false
				}
			}
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint5(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(5)
		fmt.Println("       Useless constraint    ")
		printPassed()
	}
	return true
}

func constraint6(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(6)
	}
	var otherDel int
	for i := 0; i < res.nOrder+res.nMover; i++ {
		for j := 0; j < res.nOrder; j++ {
			if i >= res.nOrder { /* mover */
				otherDel = 0
			} else { /* order */
				otherDel = res.x[i]
			}
			if res.x[j] < otherDel+distances[i][j]-(1-int(res.y[i][j]) )*B {
				if DEBUG {
					printFailed()
				}
				return false
			}
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint7(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(7)
	}
	for i := 0; i < res.nOrder; i++ {
		val := res.x[i]
		//fmt.Printf("xi = %d, ti = %d\n",val , del[i] )
		if val < del[i]-15 && res.w[i] == 0 {
			if DEBUG {
				printFailed()
			}
			return false
		}

	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint8(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(8)
	}
	for i := 0; i < res.nOrder; i++ {
		val := res.x[i]
		//fmt.Printf("xi = %d, ti = %d\n",val , del[i] )
		if val > del[i]+60 {
			if DEBUG {
				printFailed()
			}
			return false
		}
	}
	if DEBUG {
		printPassed()
	}
	return true
}

func constraint9(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(9)
	}
	for _, val := range res.x {
		if val > ES {
			if DEBUG {
				printFailed()
			}
			return false
		}
	}

	if DEBUG {
		printPassed()
	}
	return true
}

func constraint10(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(10)
	}
	for i := 0; i < res.nOrder; i++ {
		if B*int(res.z[i]) < res.x[i]-del[i]-15 {
			if DEBUG {
				fmt.Print(res.z[i], res.x[i]-del[i])
				printFailed()
			}
			return false
		}
	}

	if DEBUG {
		printPassed()
	}
	return true
}

func constraint11(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(11)
	}
	for i := 0; i < res.nOrder; i++ {
		if B*int(res.z1[i]) < res.x[i]-del[i]-30 {
			if DEBUG {
				fmt.Print(res.z1[i], res.x[i]-del[i]-30)
				printFailed()
			}
			return false
		}
	}

	if DEBUG {
		printPassed()
	}
	return true
}

func constraint12(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(12)
	}
	for i := 0; i < res.nOrder; i++ {
		if B*int(res.z2[i]) < res.x[i]-del[i]-45 {
			if DEBUG {
				printFailed()
			}
			return false
		}
	}

	if DEBUG {
		printPassed()
	}
	return true
}

func constraint13(res SolverResult, dist Distances, del DeliveryTimeVector) bool {
	if DEBUG {
		printInitMess(13)
	}
	for _, vect := range res.y {
		for _, ele := range vect {
			if ele == 0 || ele == 1 {
				continue
			}
			if DEBUG {
				printFailed()
			}
			return false
		}
	}

	if DEBUG {
		printPassed()
	}
	return true
}

func printInitMess(c int) {
	fmt.Println()
	fmt.Printf("----Constraint %d Validation----\n", c)
}

func printFailed() {
	fmt.Println("           FAILED    ")
	fmt.Println()
}

func printPassed() {
	fmt.Println("             OK    ")
	fmt.Println()
}
