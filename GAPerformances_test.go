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
	DATASET_1_ORDERS = 175
	DATASET_1_MOVERS = 30
	DATASET_2_ORDERS = 176
	DATASET_2_MOVERS = 36
	DATASET_3_ORDERS = 205
	DATASET_3_MOVERS = 38
)

func BenchmarkGreedySolver(b *testing.B) {

	utils.DeliveryTimeFilename = "input/deliveryTime_ist2.csv" //+ utils.DeliveryTimeFilename
	utils.DistanceMatrixFilename = "input/distanceMatrix_ist2.csv" //+ utils.DistanceMatrixFilename

	nOrder = DATASET_2_ORDERS
	//nMover = DATASET_1_MOVERS

	for N:= 1; N <= DATASET_2_MOVERS; N++ {


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

	switch {
	case moverPolicy == 0:
		fmt.Println("\tPolicy MINIMIZE_ACTIVE_MOVERS")
	case moverPolicy == 1:
		fmt.Println("\tPolicy MAXIMIZE_ACTIVE_MOVERS")
	}
	//fmt.Printf("\t\tPolicy %d\n", moverPolicy)

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
