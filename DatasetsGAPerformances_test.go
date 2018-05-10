package main

import (
	"time"
	"strconv"
	"orderSchedulingAlgorithm/utils"
	"fmt"
	"testing"
)

const (
	DELIVERY_TIME_PATH = "datasets_alpha/deliveryTime_ist"
	DISTANCE_MATRIX_PATH = "datasets_alpha/distanceMatrix_ist"
	CSV = ".csv"
	DATASETS = 35
	TABLE_MAXIMIZE_FILE_PATH = "results/table_maximize_policy.csv"
	TABLE_MINIMIZE_FILE_PATH = "results/table_minimize_policy.csv"
	TIME_MINIMIZE_FILE_PATH = "results/times_minimize.csv"
	TIME_MAXIMIZE_FILE_PATH = "results/times_maximize.csv"
)

func BenchmarkDatasetsGreedySolver(b *testing.B) {

	var res_max [][]int
	var res_min [][]int
	var times_min []time.Duration
	var times_max []time.Duration

	for i := 2; i <= DATASETS; i++ {

		utils.DeliveryTimeFilename = DELIVERY_TIME_PATH + strconv.Itoa(i)+CSV
		utils.DistanceMatrixFilename = DISTANCE_MATRIX_PATH + strconv.Itoa(i)+CSV

		execute()
		N := nMover

		for n:=25; n<=N; n++ {

				moverPolicy = MINIMIZE_ACTIVE_MOVERS

			test:

				CANC_TOT := 0
				COST_TOT := 0
				Z_TOT := 0
				Z1_TOT:= 0
				Z2_TOT := 0
				var TIME_TOT time.Duration = 0

				nMover = n
				results, elapsed := execute()
				CANC_TOT += results.nCancelled
				COST_TOT += results.totalCost
				TIME_TOT += elapsed
				Z_TOT += results.n1
				Z1_TOT += results.n2
				Z2_TOT += results.n3

				switch moverPolicy {
				case MINIMIZE_ACTIVE_MOVERS:
					times_min = append(times_min, TIME_TOT)
					res_min = append(res_min, []int{i, n, nOrder, COST_TOT , Z_TOT  , Z1_TOT , Z2_TOT, CANC_TOT})
				case MAXIMIZE_ACTIVE_MOVERS:
					times_max = append(times_max, TIME_TOT)
					res_max = append(res_max, []int{i, n, nOrder, COST_TOT, Z_TOT , Z1_TOT, Z2_TOT, CANC_TOT})
				}

				printTimes(n, nOrder-CANC_TOT, TIME_TOT, COST_TOT, CANC_TOT, i)

				if moverPolicy == MAXIMIZE_ACTIVE_MOVERS {
					continue
				}

				moverPolicy = MAXIMIZE_ACTIVE_MOVERS

				goto test
		}

		nMover = -1
		nOrder = -1

	}
	utils.WriteResultsTable(TABLE_MAXIMIZE_FILE_PATH, res_max, []string{"id shift","num moover","num ordini","fo","sumz","sumz1", "sumz2","sumw"})
	utils.WriteResultsTable(TABLE_MINIMIZE_FILE_PATH, res_min, []string{"id shift","num moover","num ordini","fo","sumz","sumz1", "sumz2","sumw"})
	utils.WriteResultsTimes(TIME_MINIMIZE_FILE_PATH, DATASETS, times_min, []string{"id shift","exec time"})
	utils.WriteResultsTimes(TIME_MAXIMIZE_FILE_PATH, DATASETS, times_max, []string{"id shift","exec time"})
}

func printTimes(n int, o int, t time.Duration, cost int, canc int, ist int) {

	switch {
	case moverPolicy == 0:
		fmt.Println("\tPolicy MINIMIZE_ACTIVE_MOVERS")
	case moverPolicy == 1:
		fmt.Println("\tPolicy MAXIMIZE_ACTIVE_MOVERS")
	}

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
		"RESULTS"+
			"\n\t\t%d ORDERS ASSIGNED IN %v"+
			"\n\t\t%d CANCELED ORDERS"+
			"\n\t\tTOTAL COST = %d \n", o, t, canc, cost)

	fmt.Printf("\n\n")
}