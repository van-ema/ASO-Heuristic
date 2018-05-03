package main

import (
	"time"
	"strconv"
	"orderSchedulingAlgorithm/utils"
	"fmt"
	"testing"
)

const (
	DELIVERY_TIME_PATH = "datasets/deliveryTime_ist"
	DISTANCE_MATRIX_PATH = "datasets/distanceMatrix_ist"
	CSV = ".csv"
	DATASETS = 3
	RUNS = 100
)

func BenchmarkDatasetsGreedySolver(b *testing.B) {

	var res [][]int

	for i := 2; i <= DATASETS; i++ {

		utils.DeliveryTimeFilename = DELIVERY_TIME_PATH + strconv.Itoa(i)+CSV
		utils.DistanceMatrixFilename = DISTANCE_MATRIX_PATH + strconv.Itoa(i)+CSV

		execute()
		N := nMover

		for n:=25; n<N; n++ {

				moverPolicy = MINIMIZE_ACTIVE_MOVERS

			test:

				CANC_TOT := 0
				COST_TOT := 0
				Z_TOT := 0
				Z1_TOT:= 0
				Z2_TOT := 0
				var TIME_TOT time.Duration = 0
				for i := 0; i < RUNS; i++ {
					nMover = n
					results, elapsed := execute()
					CANC_TOT += results.nCancelled
					COST_TOT += results.totalCost
					TIME_TOT += elapsed
					Z_TOT += results.n1
					Z1_TOT += results.n2
					Z2_TOT += results.n3
				}
					// policy|id shift|num moover|num ordini| f.o.| sum(z) | sum(z1)| sum(z2)| sum(w)
				res = append(res, []int{moverPolicy, i, n, nOrder, COST_TOT/RUNS, Z_TOT/RUNS , Z1_TOT/RUNS, Z2_TOT/RUNS, CANC_TOT/RUNS})

				printTimes(n, nOrder-CANC_TOT/RUNS, TIME_TOT/RUNS, COST_TOT/RUNS, CANC_TOT/RUNS, i)

				if moverPolicy == MAXIMIZE_ACTIVE_MOVERS {
					continue
				}

				moverPolicy = MAXIMIZE_ACTIVE_MOVERS

				goto test
		}

	}
	utils.WriteResultsTable("results/table.csv", res, []string{"policy","id shift","num moover","num ordini","fo","sumz","sumz1", "sumz2","sumw"})
}

func printTimes(n int, o int, t time.Duration, cost int, canc int, ist int) {

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

	fmt.Printf(header)
	fmt.Printf("Istance file n. %d", ist)
	fmt.Printf(header)

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
	//+
	//		"\n\t\t%d orders in (3,6]\n"+
	//		"\n\t\t%d orders in (6,9]\n"+
	//		"\n\t\t%d orders in (9,12]\n", RUNS, o, t, canc, cost, n1, n2, n3)

	fmt.Printf("\n\n")
}