package main

import "fmt"

const ES = 1000   /* end of shift */
const B = 10000   /* large number */

func Validate(result SolverResult, distances Distances, deliveryTime DeliveryTimeVector) bool{


	c7 := constraint7(result,distances,deliveryTime)
	c8 := constraint8(result,distances,deliveryTime)
	c9 := constraint9(result,distances,deliveryTime)
	return c7 && c8 && c9
}


func constraint7(res SolverResult , dist  Distances, del DeliveryTimeVector) bool {
	printInitMess(7)
	for i,val := range res.x {
		//fmt.Printf("xi = %d, ti = %d\n",val , del[i] )
		if val < del[i] - 15 {
			printFailed()
			return false
		}
	}
	printPassed()
	return true
}

func constraint8(res SolverResult , dist  Distances, del DeliveryTimeVector) bool {
	printInitMess(8)
	for i,val := range res.x {
		//fmt.Printf("xi = %d, ti = %d\n",val , del[i] )
		if val > del[i] + 60 {
			printFailed()
			return false
		}
	}
	printPassed()
	return true
}

func constraint9(res SolverResult , dist  Distances, del DeliveryTimeVector) bool {
	printInitMess(9)
	for _,val := range res.x {
		if val > ES {
			printFailed()
			return false
		}
	}

	printPassed()
	return true
}

func printInitMess(c int) {
	fmt.Println()
	fmt.Printf("----Constraint %d Validation----\n",c)
}

func printFailed(){
	fmt.Println("           FAILED    ")
	fmt.Println()
}

func printPassed(){
	fmt.Println("             OK    ")
	fmt.Println()
}