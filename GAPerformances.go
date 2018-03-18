package main

import (
	"fmt"
	"time"
)

const (
	RUNS = 100
	DEFAULT_POLICY = SAME_COST_LOWER_ASSIGNED_ORDER_NUMBER
)

func main() {

	TEST_MOVERS := [4]int{20,30,34,38}

	for _, N := range TEST_MOVERS {
		CANC_TOT := 0
		COST_TOT := 0
		var TIME_TOT time.Duration = 0

		for i := 0; i<RUNS; i++ {
			results, elapsed := execute(DEFAULT_POLICY, false)
			CANC_TOT += results.nCancelled
			COST_TOT += results.totalCost
			TIME_TOT += elapsed
		}
		printTimes(N, nOrder, TIME_TOT, COST_TOT, CANC_TOT)
	}
}

func printTimes(n int, o int, t time.Duration, cost int, canc int) {

	for i:=0; i<20; i++ {
		fmt.Print("-")
	}
	fmt.Printf("%d MOVERS", n)
	for i:=0; i<20; i++ {
		fmt.Print("-")
	}
	fmt.Printf("\n\t\t")

	fmt.Printf(
		"RESULTS AVERAGED AMONG %d RUNS" +
		"\n\t\t%d ORDERS ASSIGNED IN %v" +
		"\n\t\t%d CANCELED ORDERS" +
		"\n\t\tTOTAL COST = %d \n", RUNS, o, t/RUNS, canc/RUNS, cost/RUNS)

	fmt.Printf("\n\n")
}