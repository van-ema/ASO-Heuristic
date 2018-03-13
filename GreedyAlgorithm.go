package main

import (
	"orderSchedulingAlgorithm/utils"
	"fmt"
	"container/list"
	"time"
	"github.com/pborman/getopt/v2"
	"os"
)

const (
	ORDER_N = 205 // #Orders
	MOVER_N = 38  // #Movers
)

type Order struct {
	id   int    // order's id
	name string // alphanumeric name
	x    int    // target delivery time
	t    int    // delivery time
	cost int    // The cost that the order add to the final solution
}

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

	x := make([]int, nOrder)
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

func GreedySolver(nOrder, nMover int) SolverResult {

	results := initResults(nOrder, nMover)

	orders := initOrder(deliveryTimes, nOrder)
	UnfeasibleOrdersPairsMatrix = getUnfeasibleOrdersPairs(orders)

	// Each partition is a list of the orders assigned to the mover
	orderPartitions := make([]*list.List, nMover)
	for i := 0; i < nMover; i++ {
		orderPartitions[i] = list.New()
	}
	// list of cancelled orders
	cancelled := list.New()

	// Phase-1: find an assignment of orders to the movers
	toSchedule := orders.Front()
	for h := 0; h < nOrder; h++ {
		minCost := utils.Inf
		bestMover := -1
		for mover := 0; mover < nMover; mover++ {

			for e := orderPartitions[mover].Front(); e != nil; e = e.Next() {
				if UnfeasibleOrdersPairsMatrix[e.Value.(*Order).id][toSchedule.Value.(*Order).id] == 1 {
					continue
				}
			}

			// add temporary the order to schedule to the partition of the current mover
			// The output variable are set to zero because this is just a temporary partition
			// Real output variable are computed in phase-2
			cost := SingleMoverSchedulingOrders(mover, orderPartitions[mover], toSchedule, nil)
			if cost < minCost {
				minCost = cost
				bestMover = mover
			}
		}

		// The schedule is feasible for bestMover
		if bestMover != -1 {
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
		cost := SingleMoverSchedulingOrders(mover, orderPartitions[mover], nil, &results)
		results.totalCost += cost // cost cannot be inf because we know the partition can be scheduled
	}

	return results
}

// return the cost to schedule orders in the list
func SingleMoverSchedulingOrders(mover int, orders *list.List, newOrderElem *list.Element, results *SolverResult) int {

	cost := 0

	// keep the last assigned order
	var lastOrder = new(Order)
	lastOrder.x = 0
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

			newCost, nextDeliveryTime := computeCost(lastOrder.id, lastOrder.x, order)

			// If costs are equals we choose the order with the lower id
			if newCost < minCost || (newCost == minCost && nextDeliveryTime < bestDeliveryTime) {
				minOrderElem = current
				minCost = newCost
				bestDeliveryTime = nextDeliveryTime

			} else if newCost == minCost && nextDeliveryTime == bestDeliveryTime &&
				minOrderElem != nil {

				if order.id < minOrderElem.Value.(*Order).id {
					minOrderElem = current
					minCost = newCost
					bestDeliveryTime = nextDeliveryTime
				}

			}

			current = current.Next()
		}

		// Check if the new order can be scheduled
		if newOrderElem != nil {

			newOrder := newOrderElem.Value.(*Order)
			newCost, nextDeliveryTime := computeCost(lastOrder.id, lastOrder.x, newOrder)

			if newCost < minCost || (newCost == minCost && nextDeliveryTime < bestDeliveryTime) {
				minOrderElem = newOrderElem
				minCost = newCost
				bestDeliveryTime = nextDeliveryTime

			} else if newCost == minCost && nextDeliveryTime == bestDeliveryTime &&
				minOrderElem != nil {

				if newOrder.id < minOrderElem.Value.(*Order).id {
					minOrderElem = newOrderElem
					minCost = newCost
					bestDeliveryTime = nextDeliveryTime
				}

			}

		}

		// schedule not feasible
		if minCost == utils.Inf {
			return minCost
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
			iteration--
			i--
		} else {
			newOrderElem = nil
		}

		lastOrder = minOrder
	}

	// Update output
	if results != nil {
		moverOrder := new(Order)
		moverOrder.id = mover + nOrder
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

	return cost

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
		if e.Value.(*Order).t >= order.t {
			l.InsertBefore(order, e)
			return
		}
	}

	l.PushBack(order)
}

func computeCost(lastOrderId int, lastDeliveryTime int, nextOrder *Order) (int, int) {

	x := lastDeliveryTime + distances[lastOrderId][nextOrder.id]
	lateness := x - nextOrder.t
	var cost int
	switch {
	case lateness <= -15:
		//cost = utils.Inf
		cost = 0
		x = nextOrder.t - 15
	case lateness <= 15 && lateness > -15:
		cost = 0
	case lateness > 15 && lateness <= 30:
		cost = 1
	case lateness > 30 && lateness <= 45:
		cost = 2
	case lateness > 45 && lateness <= 60:
		cost = 3
	case lateness > 60:
		cost = utils.Inf
	}

	return cost, x

}

func getInput() (distMat Distances, deliveryTime DeliveryTimeVector) {
	/* init */
	distMat = make(Distances, ORDER_N+MOVER_N)
	deliveryTime = make(DeliveryTimeVector, ORDER_N)
	orderIndexToName = make(map[int]string)
	moverIndexToName = make(map[int]string)

	/* read from file */
	orderOrderDisMat, moverOrderDistMat := utils.ReadDistanceMatrix(ORDER_N)
	deliveryTimesMap := utils.ReadOrdersTargetTime()

	for orderKey, distVector := range orderOrderDisMat {
		distMat[orderKey.I] = make([]int, ORDER_N)

		/* update additional order info */
		orderIndexToName[orderKey.I] = orderKey.N
		deliveryTime[orderKey.I] = deliveryTimesMap[orderKey.N]

		for j, distance := range distVector {
			distMat[orderKey.I][j] = distance
		}

	}

	for moverKey, distVector := range moverOrderDistMat {
		distMat[ORDER_N+moverKey.I] = make([]int, ORDER_N)

		/* update additional mover info */
		moverIndexToName[moverKey.I] = moverKey.N
		for j, distance := range distVector {
			distMat[ORDER_N+moverKey.I][j] = distance
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

	alpha := 75
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

	alpha = 60
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

func init() {
	getopt.FlagLong(&utils.DistanceMatrixFilename, "distanceMat", 'd', "distance matrix filename")
	getopt.FlagLong(&utils.DeliveryTimeFilename, "deliveryTimes", 't', "delivery times vector filename")

	getopt.FlagLong(&nOrder, "nOrder", 'n', "number of orders")
	getopt.FlagLong(&nMover, "nMover", 'm', "number of movers")
}

func main() {

	getopt.Parse()
	if !utils.Exist(utils.DeliveryTimeFilename) || !utils.Exist(utils.DistanceMatrixFilename) {
		getopt.Usage()
		os.Exit(1)
	}

	start := time.Now()
	distances, deliveryTimes = getInput()
	if len(deliveryTimes) != nOrder {
		fmt.Errorf("len of delivery time vector != #orders\r\n")
		getopt.Usage()
		os.Exit(1)
	}

	results := GreedySolver(nOrder, nMover)

	utils.WriteAdjMatOnFile("y.csv", results.y, orderIndexToName, moverIndexToName)
	utils.WriteOrderVectorInt("x.csv", results.x, orderIndexToName, []string{"order", "x"})
	utils.WriteOrderVectorUint8("w.csv", results.w, orderIndexToName, []string{"order", "w"})
	utils.WriteOrderVectorUint8("z.csv", results.z, orderIndexToName, []string{"order", "z"})
	utils.WriteOrderVectorUint8("z1.csv", results.z1, orderIndexToName, []string{"order", "z1"})
	utils.WriteOrderVectorUint8("z2.csv", results.z2, orderIndexToName, []string{"order", "z2"})

	elapsed := time.Since(start)

	if Validate(results, distances, deliveryTimes) {
		fmt.Printf("The solution is admissible\r\n")
	} else {
		fmt.Printf("The solution is not admissible\r\n")
	}

	fmt.Printf("#Order, #Mover\r\n")
	fmt.Printf("%d,%d\r\n", nOrder, nMover)
	fmt.Printf("Solver took %s\r\n", elapsed)
	fmt.Printf("Total cost: %d\r\n", results.totalCost)
	fmt.Printf("#order in (15,30] %d\r\n", results.n1)
	fmt.Printf("#order in (30,45] %d\r\n", results.n2)
	fmt.Printf("#order in (45,60] %d\r\n", results.n3)

}
