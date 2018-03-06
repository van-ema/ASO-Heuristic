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
	id int // order's id
	x  int // deprecated
	t  int // delivery time
}

var nOrder int
var nMover int

var distances [][]int
var deliveryTimes []int

// return for each mover a list of orders
func GreedySolver(totalCost *int) ([][]uint8, []int, []uint8, []uint8, []uint8, []uint8) {

	nRow := nMover + nOrder
	nCol := nOrder
	y := make([][]uint8, nRow)
	for i := 0; i < nRow; i++ {
		y[i] = make([]uint8, nCol)
	}

	x := make([]int, nOrder)
	w := make([]uint8, nOrder)

	z := make([]uint8, nOrder)
	z1 := make([]uint8, nOrder)
	z2 := make([]uint8, nOrder)

	// alias for D - D* : contains the orders not assigned
	orders := initOrder(deliveryTimes)
	// keep the effective length of the list
	// from position length to the end of the list we can find scheduled orders
	length := orders.Len()

	// keep the total cost of the solution
	*totalCost = 0
	// Solving for each mover
	for mover := 0; mover < nMover; mover++ {

		// Orders are never removed from the orders list but they are moved on the back
		// We use the var length to keep the number of orders not already scheduled
		// placed at the beginning of the list
		*totalCost += SingleMoverSchedulingOrders(mover, y, x, z, z1, z2, orders, &length)

	}

	// The orders still in the list are all those that cannot be scheduled
	e := orders.Front()
	for k := 0; k < length; k++ {
		*totalCost += 10
		cancelled := e.Value.(*Order)
		w[cancelled.id] = 1
		e = e.Next()

	}

	return y, x, w, z, z1, z2
}

// return a list of orders that can be scheduled by the mover
func SingleMoverSchedulingOrders(mover int, y [][]uint8, x []int, z, z1, z2 []uint8, orders *list.List, length *int) int {

	cost := 0
	var lastOrder = new(Order)
	lastOrder.id = nOrder + mover
	lastOrder.t = 0

	for i := 0; i < *length; i++ {

		var minOrderElem *list.Element
		minCost := utils.Inf
		newDeliveryTime := 0

		current := orders.Front()
		for j := 0; j < *length; j++ {
			order := current.Value.(*Order)

			newCost, nextDeliveryTime := computeCost(lastOrder.id, lastOrder.t, order)
			if newCost < minCost {
				minOrderElem = current
				minCost = newCost
				newDeliveryTime = nextDeliveryTime
			}

			current = current.Next()
		}

		// schedule not feasible, try with the next mover
		if minCost == utils.Inf {
			break
		}

		// schedule is feasible
		cost += minCost
		minOrder := minOrderElem.Value.(*Order)

		// Update results
		switch minCost {
		case 1:
			z[minOrder.id] = 1
		case 2:
			z1[minOrder.id] = 1
		case 3:
			z2[minOrder.id] = 1
		default:
			z[minOrder.id] = 0
		}
		x[minOrder.id] = newDeliveryTime
		y[lastOrder.id][minOrder.id] = 1
		// Remove order from list of orders to schedule
		orders.MoveToBack(minOrderElem)
		lastOrder = minOrder
		*length--
		i--
	}

	return cost

}

// orders allocation and initialization with delivery times
// We use linked list because we must frequently remove assigned orders
func initOrder(deliveryTimes []int) *list.List {

	n := len(deliveryTimes)
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
	case lateness < -15:
		cost = utils.Inf
	case lateness < 15 && lateness > -15:
		cost = 0
	case lateness >= 15 && lateness < 30:
		cost = 1
	case lateness >= 30 && lateness < 45:
		cost = 2
	case lateness >= 45 && lateness < 60:
		cost = 3
	case lateness >= 60:
		cost = utils.Inf
	}

	return cost, x

}

func printResults(results [][]*Order) {
	for k := range results {
		/*if k >= 0 {
			fmt.Printf("%s%d : ", "Mover-", k)
		} else {
			fmt.Printf("cancelled:")
		} */
		if k == len(results)-1 { /* */
			fmt.Printf("cancelled:")
		} else {
			fmt.Printf("%s%d : ", "Mover-", k)
		}

		printArray(results[k])

		fmt.Print("\n")
	}
}

func printArray(p []*Order) {
	for _, v := range p {
		fmt.Printf("[id:%d,x:%d,t:%d]", v.id, v.x, v.t)
	}
}

// return for each mover a list of orders that he can schedule
// For each order it find the mover that minimize the cost
func BaseSolver() ([][]*Order, int) {

	// alias for D - D* : contains the orders not assigned
	orders := initOrder(deliveryTimes)
	// keep the effective length of the list
	// from position length to the end of the list we can find scheduled orders
	length := orders.Len()

	// At the position i we find the schedule for the mover i-th
	results := make([][]*Order, nMover+1)
	for i := 0; i < nMover; i++ {
		results[i] = make([]*Order, 0, int(nOrder/nMover))
	}

	// keep the total cost of the solution
	totalCost := 0
	for i := 0; i < length; i++ {
		e := orders.Front()
		order := e.Value.(*Order)
		moverId, cost, x := findBestMover(order, results, distances)

		if moverId >= 0 {
			order.x = x
			results[moverId] = append(results[moverId], order)
			totalCost += cost

			orders.MoveToBack(e)
			length--
			i--
		}

		e = e.Next()
	}

	var cancelled []*Order
	e := orders.Front()
	for i := 0; i < length; i++ {
		totalCost += 10
		cancelled = append(cancelled, e.Value.(*Order))
		e.Next()
	}

	results[nMover] = cancelled

	return results, totalCost

}

func findBestMover(order *Order, results [][]*Order, dist [][]int, ) (int, int, int) {

	minCost := utils.Inf
	bestMover := -1
	bestDeliveryTime := 0
	for mover := 0; mover < nMover; mover++ {

		res := results[mover]
		var lastOrder int
		var lastDeliveryTime int
		if len(res) == 0 {
			lastOrder = nOrder + mover
			lastDeliveryTime = 0
		} else {
			lastOrder = results[mover][len(res)-1].id
			lastDeliveryTime = results[mover][len(res)-1].x
		}

		cost, x := computeCost(lastOrder, lastDeliveryTime, order)
		if cost < minCost {
			minCost = cost
			bestMover = mover
			bestDeliveryTime = x
		}
	}

	return bestMover, minCost, bestDeliveryTime

}

func main() {
	nOrder = 10
	nMover = 2

	distances = utils.CreateOrderMatrix(nOrder, nMover)
	deliveryTimes = utils.CreateDeliveryTimeVector(nOrder)

	utils.PrintDistanceMatrix(distances, nOrder)
	fmt.Print("Algorithm 1:\n")
	start := time.Now()
	var cost int
	y, x, w, z, z1, z2 := GreedySolver(&cost)
	elapsed := time.Since(start)

	utils.PrintAssigmentMatrix(y, nOrder)
	fmt.Println(x)
	fmt.Println(w)
	fmt.Println(z, z1, z2)
	fmt.Printf("Solver took %s\n", elapsed)
	fmt.Printf("Total cost: %d\n", cost)

	//fmt.Print("\n\nAlgorithm 2:\n")
	//start = time.Now()
	//res1, cost1 := BaseSolver()
	//elapsed = time.Since(start)
	////printResults(res1)
	//fmt.Printf("Solver took %s\n", elapsed)
	//fmt.Printf("Total cost: %d\n", cost1)

}
