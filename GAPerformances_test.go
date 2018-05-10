package main

import (
	"time"
	"testing"
	"orderSchedulingAlgorithm/utils"
)

const (
	DATASET_1_ORDERS = 175
	DATASET_1_MOVERS = 30
	DATASET_2_ORDERS = 176
	DATASET_2_MOVERS = 36
	DATASET_3_ORDERS = 205
	DATASET_3_MOVERS = 38
	DATASET = 1
)

func BenchmarkGreedySolver(b *testing.B) {

	switch (DATASET) {
	case 1:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist1.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist1.csv"
	case 2:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist2.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist2.csv"
	case 3:
		utils.DeliveryTimeFilename = "input/deliveryTime_ist3.csv"
		utils.DistanceMatrixFilename = "input/distanceMatrix_ist3.csv"
	}
	utils.DeliveryTimeFilename = "datasets/deliveryTime_ist2.csv"     //+ utils.DeliveryTimeFilename
	utils.DistanceMatrixFilename = "datasets/distanceMatrix_ist2.csv" //+ utils.DistanceMatrixFilename

	nOrder = 239
	nMover = 36

	for N:= 28; N<=36; N++ {

			moverPolicy = MINIMIZE_ACTIVE_MOVERS

		test:

			CANC_TOT := 0
			COST_TOT := 0
			var TIME_TOT time.Duration = 0

			nMover = N
			results, elapsed := execute()
			CANC_TOT += results.nCancelled
			COST_TOT += results.totalCost
			TIME_TOT += elapsed

			printTimes(N, nOrder-CANC_TOT, TIME_TOT, COST_TOT, CANC_TOT, DATASET)

			if moverPolicy == MAXIMIZE_ACTIVE_MOVERS {
				continue
			}

			moverPolicy = MAXIMIZE_ACTIVE_MOVERS

			goto test
	}
}
