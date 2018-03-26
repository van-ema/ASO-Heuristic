package main

import (
	"time"
	"testing"
	"fmt"
	"strconv"
	"orderSchedulingAlgorithm/utils"
)

const (
	RUNS = 100
	DATASET = 1
)

func BenchmarkGreedySolver(b *testing.B) {

	var MOVERS int
	switch (DATASET) {
	case 1:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist1.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist1.csv"
		MOVERS = 30
		nOrder = 175
	case 2:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist2.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist2.csv"
		MOVERS = 36
		nOrder = 179
	case 3:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist3.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist3.csv"
		MOVERS = 38
		nOrder = 205
	}

	for N:= 1; N<=MOVERS; N++ {

		moverPolicy = MINIMIZE_ACTIVE_MOVERS

		test:

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
		printTimes(N, nOrder-CANC_TOT/RUNS, TIME_TOT/RUNS, COST_TOT/RUNS, CANC_TOT/RUNS)

		if moverPolicy == MAXIMIZE_ACTIVE_MOVERS {
			continue
		}

		moverPolicy = MAXIMIZE_ACTIVE_MOVERS

		goto test
	}
}

func printTimes(n int, o int, t time.Duration, cost int, canc int) {

	fmt.Printf("\t\tPolicy %d\n", moverPolicy)

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
			"\n\t\tTOTAL COST = %d \n", RUNS, o, t, canc, cost)

	fmt.Printf("\n\n")
}
