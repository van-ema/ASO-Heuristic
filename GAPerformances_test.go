package main

import (
	"time"
	"testing"
	"fmt"
	"strconv"
)

const (
	RUNS = 100
)

func BenchmarkGreedySolver(b *testing.B) {

	TEST_MOVERS := [4]int{2, 3, 4, 5}

	for _, N := range TEST_MOVERS {
		CANC_TOT := 0
		COST_TOT := 0
		var TIME_TOT time.Duration = 0

		for i := 0; i < RUNS; i++ {
			nMover = N
			results, elapsed := execute()
			CANC_TOT += results.nCancelled
			COST_TOT += results.totalCost
			TIME_TOT += elapsed
		}
		printTimes(N, nOrder, TIME_TOT, COST_TOT, CANC_TOT)
	}
}

func printTimes(n int, o int, t time.Duration, cost int, canc int) {

	header := ""
	for i := 0; i < 20; i++ {
		header += "-"
	}

	header += strconv.Itoa(n) + " " + "MOVERS"
	for i := 0; i < 20; i++ {
		header += "-"
	}
	header+= "\n\t\t"
	fmt.Printf(header)

	fmt.Printf(
		"RESULTS AVERAGED AMONG %d RUNS"+
			"\n\t\t%d ORDERS ASSIGNED IN %v"+
			"\n\t\t%d CANCELED ORDERS"+
			"\n\t\tTOTAL COST = %d \n", RUNS, o, t/RUNS, canc/RUNS, cost/RUNS)

	fmt.Printf("\n\n")
}
