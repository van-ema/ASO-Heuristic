package main

import (
	"orderSchedulingAlgorithm/utils"
	"fmt"
	"container/list"
	"time"
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

var nOrder int
var nMover int

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
	nOrder    int
	nMover    int
	totalCost int
	y         [][]uint8
	x         []int
	w         []uint8
	z         []uint8
	z1        []uint8
	z2        []uint8
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

	nAssigned := 0  // number of orders assigned
	nCancelled := 0 // number of cancelled orders

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
			nAssigned++

		} else {
			cancelled.PushFront(toSchedule)
			results.totalCost += 10
			results.w[toSchedule.Value.(*Order).id] = 1
			nCancelled++
		}

		toSchedule = toSchedule.Next()
	}

	// Phase-2 : Re-compute final schedule for each mover and update output structures
	for mover := 0; mover < nMover; mover++ {
		cost := SingleMoverSchedulingOrders(mover, orderPartitions[mover], nil, &results)
		results.totalCost += cost // cost cannot be inf because we know the partition can be scheduled
	}

	fmt.Printf("assigned: %d ; cancelled: %d\n", nAssigned, nCancelled)

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
				} else if order.cost == 2 {
					results.z[order.id] = 1
					results.z1[order.id] = 1
				} else if order.cost == 3 {
					results.z[order.id] = 1
					results.z1[order.id] = 1
					results.z2[order.id] = 1
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

func
main() {
	nOrder = ORDER_N
	nMover = 20

	//distances = utils.CreateOrderMatrix(nOrder, nMover)
	//deliveryTimes = utils.CreateDeliveryTimeVector(nOrder)

	distances, deliveryTimes = getInput()

	/* TODO put in other place */
	orders := initOrder(deliveryTimes, nOrder)
	UnfeasibleOrdersPairsMatrix = getUnfeasibleOrdersPairs(orders)
	//utils.PrintMatrix(UnfeasibleOrdersPairsMatrix)

	//utils.PrintDistanceMatrix(distances, nOrder)
	fmt.Print("Algorithm 1:\n")
	start := time.Now()
	results := GreedySolver(nOrder, nMover)
	elapsed := time.Since(start)
	//printResults(res)

	utils.PrintAssigmentMatrix(results.y, nOrder)
	fmt.Println(results.x)
	fmt.Println(results.w)
	fmt.Println(results.z, results.z1, results.z2)
	fmt.Printf("Solver took %s\n", elapsed)
	fmt.Printf("Total cost: %d\n", results.totalCost)

	if Validate(results, distances, deliveryTimes) {
		fmt.Printf("VALID")
	} else {
		fmt.Printf("NOT VALID")
	}

	//fmt.Print("\n\nAlgorithm 2:\n")
	//start = time.Now()
	//y, x, w, z, z1, z2 = BaseSolver(&cost)
	//elapsed = time.Since(start)
	////printResults(res)
	//
	//utils.PrintAssigmentMatrix(y, nOrder)
	//fmt.Println(x)
	//fmt.Println(w)
	//fmt.Println(z, z1, z2)
	//fmt.Printf("Solver took %s\n", elapsed)
	//fmt.Printf("Total cost: %d\n", cost)

}
