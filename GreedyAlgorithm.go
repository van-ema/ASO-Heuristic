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
	"orderSchedulingAlgorithm/utils"
	"fmt"
	"container/list"
	"time"
	"github.com/pborman/getopt/v2"
	"os"
	"strconv"
	"strings"
)

var (
	ORDER_N = -1  // #Orders if not specified
	MOVER_N = -1  // #Movers if not specified
)

var DEBUG = true

type Order struct {
	id   int    // order's id
	name string // alphanumeric name
	x    int    // target delivery time
	t    int    // delivery time
	cost int    // The cost that the order add to the final solution
}

const (
	MINIMIZE_ACTIVE_MOVERS = 0
	MAXIMIZE_ACTIVE_MOVERS = 1
)

var moverPolicy = MAXIMIZE_ACTIVE_MOVERS

type Distances [][]int
type DeliveryTimeVector []int

var nOrder = ORDER_N
var nMover = MOVER_N

var distances [][]int
var deliveryTimes DeliveryTimeVector

var UnfeasibleOrdersPairsMatrix [][]uint8

/* orderIndexToName[order index] = alphanumeric ID (given in CSV files) */
var orderIndexToName map[int]string
/* moverIndexToName[mover index] = alphanumeric ID (given in CSV files);
 * mover index starts from 0 to |M| while in the matrix it starts at
 * ORDER_N
 */
var moverIndexToName map[int]string

type SolverResult struct {
	nOrder     int
	nMover     int
	totalCost  int
	y          [][]uint8
	x          []int
	w          []uint8
	z          []uint8
	z1         []uint8
	z2         []uint8
	n1         int
	n2         int
	n3         int
	nAssigned  int
	nCancelled int
}

func initResults(nOrder, nMover int) (res SolverResult) {

	y := make([][]uint8, nOrder+nMover)
	for i := 0; i < nOrder+nMover; i++ {
		y[i] = make([]uint8, nOrder)
	}

	x := make([]int, nOrder+nMover)
	w := make([]uint8, nOrder)

	z := make([]uint8, nOrder)
	z1 := make([]uint8, nOrder)
	z2 := make([]uint8, nOrder)

	return SolverResult{
		totalCost: 0,
		nMover:    nMover,
		nOrder:    nOrder,
		x:         x,
		y:         y,
		w:         w,
		z:         z,
		z1:        z1,
		z2:        z2,}
}

func makeOrderPartition() []*list.List {
	orderPartitions := make([]*list.List, nMover)
	for i := 0; i < nMover; i++ {
		orderPartitions[i] = list.New()
	}

	return orderPartitions
}

func GreedySolver(nOrder, nMover int) SolverResult {

	results := initResults(nOrder, nMover)                         // Keep output struct
	orders := initOrder(deliveryTimes, nOrder)                     // Get linked list of orders to schedule
	UnfeasibleOrdersPairsMatrix = getUnfeasibleOrdersPairs(orders) // orders pair i,j cannot be schedule if (i,j) value is 1
	orderPartitions := makeOrderPartition()                        // Each partition is a list of the orders assigned to the mover
	cancelled := list.New()                                        // list of cancelled orders

	costList := make([]int, nMover)
	for o := range costList {
		costList[o] = -1
	}

	// Phase-1: find an assignment of orders to the movers
	toSchedule := orders.Front()
	for h := 0; h < nOrder; h++ {
		minCost := utils.Inf
		bestMover := -1
		var cancelledOrderPresent bool
		var cancelledOrderMin bool

		for mover := 0; mover < nMover; mover++ {

			var cost int
			/*for e := orderPartitions[mover].Front(); e != nil; e = e.Next() {
				if UnfeasibleOrdersPairsMatrix[e.Value.(*Order).id][toSchedule.Value.(*Order).id] == 1 {
					continue
				}
			} */

			// add temporary the order to schedule to the partition of the current mover
			// The output variable are set to zero because this is just a temporary partition
			// Real output variable are computed in phase-2
			cost, cancelledOrderPresent = SingleMoverSchedulingOrders(mover, orderPartitions[mover], toSchedule, nil)
			costList[mover] = cost
			if cost < minCost {
				minCost = cost
				bestMover = mover
				if DEBUG {
					fmt.Println(bestMover)
				}
				cancelledOrderMin = cancelledOrderPresent
			}
		}

		// The schedule is feasible for bestMover
		if !cancelledOrderMin {

			switch moverPolicy {
			case MINIMIZE_ACTIVE_MOVERS:
				bestMover = bestMover
			case MAXIMIZE_ACTIVE_MOVERS:
				mNumOrders := make([]int, nMover)
				mBestNumb := 0
				for k, v := range costList {

					if v == minCost {
						mBestNumb++
						mNumOrders[k] = orderPartitions[k].Len()
					} else {
						mNumOrders[k] = utils.Inf
					}
				}

				min := 0
				for k, v := range mNumOrders {
					if v <= min {
						min = v
						bestMover = k
					}
				}
			}
			insert(orderPartitions[bestMover], toSchedule.Value.(*Order))
			results.nAssigned++
		} else {
			cancelled.PushFront(toSchedule)
			results.totalCost += 10
			results.w[toSchedule.Value.(*Order).id] = 1
			results.nCancelled++
		}

		toSchedule = toSchedule.Next()
	}

	// Phase-2 : Re-compute final schedule for each mover and update output structures
	for mover := 0; mover < nMover; mover++ {
		cost, _ := SingleMoverSchedulingOrders(mover, orderPartitions[mover], nil, &results)
		results.totalCost += cost // cost cannot be inf because we know the partition can be scheduled
	}

	for mover := 0; mover < nMover; mover++ {
		if DEBUG {
			fmt.Printf("Mover-%d", mover)
		}
		for e := orderPartitions[mover].Front(); e != nil; e = e.Next() {
			order := e.Value.(*Order)
			if DEBUG {
				fmt.Printf("[id: %d, t: %d, x: %d]", order.id, order.t, order.x)
			}
		}
		if DEBUG {
			fmt.Printf("\n")
		}
	}

	return results
}

// return the cost to schedule orders in the list
func SingleMoverSchedulingOrders(mover int, orders *list.List, newOrderElem *list.Element, results *SolverResult) (cost int, cancelled bool) {

	cost = 0
	cancelled = false
	// keep the last assigned order
	var lastOrder = new(Order)

	// TODO mover starts 5 min early
	//lastOrder.x = -1 //0
	lastOrder.id = nOrder + mover

	length := orders.Len()
	iteration := length
	if newOrderElem != nil {
		iteration += 1
	}

	for i := 0; i < iteration; i++ {

		var minOrderElem *list.Element // order that minimize the cost
		minCost := utils.Inf           // keep the cost of the favourable order
		bestDeliveryTime := 0

		current := orders.Front()
		for j := 0; j < length; j++ {

			order := current.Value.(*Order)

			newCost, nextDeliveryTime, orderCancelled := computeCost(lastOrder.id, lastOrder.x, order)

			// If costs are equals we choose the order with the lower id
			if newCost < minCost || (newCost == minCost && nextDeliveryTime < bestDeliveryTime) {
				minOrderElem = current
				minCost = newCost
				bestDeliveryTime = nextDeliveryTime

				cancelled = orderCancelled

			} else if newCost == minCost && nextDeliveryTime == bestDeliveryTime &&
				minOrderElem != nil {

				if order.id < minOrderElem.Value.(*Order).id {
					minOrderElem = current
					minCost = newCost
					bestDeliveryTime = nextDeliveryTime

					cancelled = orderCancelled
				}

			}

			current = current.Next()
		}

		// Check if the new order can be scheduled
		if newOrderElem != nil {

			newOrder := newOrderElem.Value.(*Order)
			newCost, nextDeliveryTime, newOrderCancelled := computeCost(lastOrder.id, lastOrder.x, newOrder)

			if newCost < minCost || (newCost == minCost && nextDeliveryTime < bestDeliveryTime) {
				minOrderElem = newOrderElem
				minCost = newCost
				bestDeliveryTime = nextDeliveryTime

				cancelled = newOrderCancelled

			} else if newCost == minCost && nextDeliveryTime == bestDeliveryTime &&
				minOrderElem != nil {

				if newOrder.id < minOrderElem.Value.(*Order).id {
					minOrderElem = newOrderElem
					minCost = newCost
					bestDeliveryTime = nextDeliveryTime

					cancelled = newOrderCancelled
				}

			}

		}

		// schedule is feasible
		cost += minCost
		minOrder := minOrderElem.Value.(*Order)
		minOrder.x = bestDeliveryTime
		minOrder.cost = minCost

		// Remove order from list of orders to schedule
		if minOrderElem != newOrderElem {
			orders.MoveToBack(minOrderElem)
			length--
		} else {
			newOrderElem = nil
		}

		lastOrder = minOrder
	}

	// Update output
	if results != nil {
		moverOrder := new(Order)
		moverOrder.id = mover + nOrder
		moverOrder.x = -1 // TODO check
		orders.PushFront(moverOrder)

		for e := orders.Front(); e != nil; e = e.Next() {
			order := e.Value.(*Order)
			next := e.Next()
			if next != nil {
				next := next.Value.(*Order)
				results.y[order.id][next.id] = 1
			}

			if order.id < nOrder {
				results.x[order.id] = order.x
				if order.cost == 1 {
					results.z[order.id] = 1
					results.n1++
				} else if order.cost == 2 {
					results.z[order.id] = 1
					results.z1[order.id] = 1
					results.n2++
				} else if order.cost == 3 {
					results.z[order.id] = 1
					results.z1[order.id] = 1
					results.z2[order.id] = 1
					results.n3++
				}
			}
		}

	}

	return cost, cancelled

}

// orders allocation and initialization with delivery times
// We use linked list because we must frequently remove assigned orders
func initOrder(deliveryTimes []int, n int) *list.List {

	orderList := list.New()

	for i := 0; i < n; i++ {
		var order = new(Order)
		order.id = i
		order.t = deliveryTimes[i]
		insert(orderList, order)
	}

	return orderList
}

// insert an order in the right position to keep the linked list sorted
func insert(l *list.List, order *Order) {
	if l.Len() == 0 {
		l.PushFront(order)
		return
	}

	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*Order).t > order.t {
			l.InsertBefore(order, e)
			return
		}
	}

	l.PushBack(order)
}

func computeCost(lastOrderId int, lastDeliveryTime int, nextOrder *Order) (cost int, x int, cancelled bool) {

	x = lastDeliveryTime + distances[lastOrderId][nextOrder.id]
	lateness := x - nextOrder.t

	cancelled = false
	switch {
	case lateness <= -3:
		cost = 0
		x = nextOrder.t - 3
	case lateness <= 3 && lateness > -3:
		cost = 0
	case lateness > 3 && lateness <= 6:
		cost = 1
	case lateness > 6 && lateness <= 9:
		cost = 2
	case lateness > 9 && lateness <= 12:
		cost = 3
	case lateness > 12:
		cost = 10
		cancelled = true
	}

	return cost, x, cancelled

}

func getInput() (distMat Distances, deliveryTime DeliveryTimeVector) {


	/* read from file */

	deliveryTimesMap := utils.ReadOrdersTargetTime()

	if nOrder == -1 {
		nOrder = len(deliveryTimesMap)
	}

	orderOrderDisMat, moverOrderDistMat := utils.ReadDistanceMatrix(nOrder)


	if len(moverOrderDistMat) < nMover {
		fmt.Printf("too many mover! The distances matrix contains only %d mover \r\n", len(moverOrderDistMat))
		os.Exit(1)
	} else if nMover == -1 {
		nMover = len(moverOrderDistMat)
	}

	/* init */
	distMat = make(Distances, nOrder+nMover)
	deliveryTime = make(DeliveryTimeVector, nOrder)
	orderIndexToName = make(map[int]string)
	moverIndexToName = make(map[int]string)

	for orderKey, distVector := range orderOrderDisMat {

		if orderKey.I >= nOrder {
			continue
		}
		distMat[orderKey.I] = make([]int, nOrder)

		/* update additional order info */
		orderIndexToName[orderKey.I] = orderKey.N
		deliveryTime[orderKey.I] = deliveryTimesMap[orderKey.N]

		for j, distance := range distVector {
			if j >= nOrder {
				continue
			}
			distMat[orderKey.I][j] = distance
		}

	}

	for moverKey, distVector := range moverOrderDistMat {

		if moverKey.I >= nMover {
			continue
		}

		distMat[nOrder+moverKey.I] = make([]int, nOrder)

		/* update additional mover info */
		moverIndexToName[moverKey.I] = moverKey.N
		for j, distance := range distVector {
			if j >= nOrder {
				continue
			}
			distMat[nOrder+moverKey.I][j] = distance
		}
	}

	/*
	for i,e := range distMat {
		fmt.Printf("index = %d, vector : %d\n",i, e)
	}

	for k,v := range moverIndexToName {
		fmt.Printf("index = %d, order name : %s\n",k, v)
	} */

	return distMat, deliveryTime
}

func getUnfeasibleOrdersPairs(orders *list.List) [][]uint8 {

	notFeasiblePair := make([][]uint8, nOrder+nMover)
	for i := 0; i < nOrder+nMover; i++ {
		notFeasiblePair[i] = make([]uint8, nOrder)
	}

	//alpha := 75
	alpha := 15
	for i := orders.Front(); i != nil; i = i.Next() {
		for j := orders.Front(); j != nil; j = j.Next() {
			first := i.Value.(*Order)
			sec := j.Value.(*Order)
			if i == j || sec.t-first.t+alpha < distances[first.id][sec.id] {
				// ti−ti0 +α < d(i0 , i)
				notFeasiblePair[first.id][sec.id] = 1
			}
		}
	}

	//alpha = 60
	alpha = 12
	for i := nOrder; i < nMover+nOrder; i++ {
		for j := orders.Front(); j != nil; j = j.Next() {
			first := j.Value.(*Order)

			if first.t+alpha < distances[i][first.id] {
				notFeasiblePair[i][first.id] = 1
			}
		}
	}

	return notFeasiblePair
}

func execute() (SolverResult, time.Duration) {

	//distances = utils.CreateOrderMatrix(nOrder, nMover)
	//deliveryTimes = utils.CreateDeliveryTimeVector(nOrder)
	distances, deliveryTimes = getInput()

	validateInput()

	start := time.Now()

	results := GreedySolver(nOrder, nMover)

	elapsed := time.Since(start)

	validateResults(results)

	writeResultsToFile(results)

	if DEBUG {
		printFinal(elapsed, results)
	}

	return results, elapsed
}

func validateInput() {
	if len(deliveryTimes) != nOrder {
		fmt.Printf("len of delivery time vector %d != #orders %d\r\n",len(deliveryTimes),nOrder)
		getopt.Usage()
		os.Exit(1)
	}
}

func init() {
	getopt.FlagLong(&utils.DistanceMatrixFilename, "distanceMat", 'd', "distance matrix filename")
	getopt.FlagLong(&utils.DeliveryTimeFilename, "deliveryTimes", 't', "delivery times vector filename")

	getopt.FlagLong(&nOrder, "nOrder", 'n', "number of orders")
	getopt.FlagLong(&nMover, "nMover", 'm', "number of movers")
	getopt.FlagLong(&DEBUG, "debug", 'i', "execute in debug mode: Extra output info")
	getopt.FlagLong(&moverPolicy,"policy", 'p', "policy to balance number of orders among movers")
}

func printFinal(elapsed time.Duration, results SolverResult) {
	fmt.Printf("#Order, #Mover\r\n")
	fmt.Printf("  %d,     %d  \r\n", nOrder, nMover)
	fmt.Printf("Solver took %s\r\n", elapsed)
	fmt.Printf("Policy %d\r\n", moverPolicy)
	fmt.Printf("Total cost: %d\r\n", results.totalCost)
	fmt.Printf("assigned: %d, cancelled %d\r\n", results.nAssigned, results.nCancelled)
	fmt.Printf("#order in (3,6] %d\r\n", results.n1)
	fmt.Printf("#order in (6,9] %d\r\n", results.n2)
	fmt.Printf("#order in (9,12] %d\r\n", results.n3)
}
func validateResults(results SolverResult) bool {
	if Validate(results, distances, deliveryTimes) {
		if DEBUG {
			fmt.Printf("The solution is admissible\r\n")
		}
		return true
	}
	if DEBUG {
		fmt.Printf("The solution is NOT admissible\r\n")
	}
	return false

}

func createOutputPath() string {
	begin := strings.Index(utils.DeliveryTimeFilename, "ist")
	end := strings.Index(utils.DeliveryTimeFilename, ".")
	ist := utils.DeliveryTimeFilename[begin:end]

	pathList := []string{"results","results/"+ist}

	output := ""
	if moverPolicy == 0 {
		dir := pathList[1]+"/minimize_active_movers/"
		rel := dir+strconv.Itoa(nMover)+"/"
		pathList = append(pathList, dir)
		pathList = append(pathList, rel)
		output += rel
	} else {
		dir := pathList[1]+"/maximize_active_movers/"
		rel := dir+strconv.Itoa(nMover)+"/"
		pathList = append(pathList, dir)
		pathList = append(pathList, rel)
		output += rel
	}
	for i := range pathList {
		if _, err := os.Stat(pathList[i]); os.IsNotExist(err) {
			os.Mkdir(pathList[i], os.ModePerm)
		}
	}
	return output
}

func writeResultsToFile(results SolverResult) {
	OUTPUT := createOutputPath()

	utils.WriteAdjMatOnFile(OUTPUT+"y.csv", results.y, orderIndexToName, moverIndexToName)
	utils.WriteOrderVectorInt(OUTPUT+"x.csv", results.x, orderIndexToName, moverIndexToName, nOrder, []string{"order", "x"})
	utils.WriteOrderVectorUint8(OUTPUT+"w.csv", results.w, orderIndexToName, []string{"order", "w"})
	utils.WriteOrderVectorUint8(OUTPUT+"z.csv", results.z, orderIndexToName, []string{"order", "z"})
	utils.WriteOrderVectorUint8(OUTPUT+"z1.csv", results.z1, orderIndexToName, []string{"order", "z1"})
	utils.WriteOrderVectorUint8(OUTPUT+"z2.csv", results.z2, orderIndexToName, []string{"order", "z2"})
}

func main() {
	getopt.Parse()
	utils.DeliveryTimeFilename = "datasets/" + utils.DeliveryTimeFilename
	utils.DistanceMatrixFilename = "datasets/" + utils.DistanceMatrixFilename
	if !utils.Exist(utils.DeliveryTimeFilename) || !utils.Exist(utils.DistanceMatrixFilename) {
		getopt.Usage()
		fmt.Printf("The given files do not exist.\r\n")
		os.Exit(1)
	}
	execute()
}
