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
	x  int // delivery time
	t  int // target delivery time
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

	for i := 0; i < *length; i++ {

		var minOrderElem *list.Element // order that minimize the cost
		minCost := utils.Inf           // keep the cost of the favourable order
		newDeliveryTime := 0           // keep the delivery time of the last scheduled order

		current := orders.Front()
		for j := 0; j < *length; j++ {
			order := current.Value.(*Order)

			newCost, nextDeliveryTime := computeCost(lastOrder.id, newDeliveryTime, order)
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
		if minCost == 1 {
			z[minOrder.id] = 1
		} else if minCost == 2 {
			z1[minOrder.id] = 1
		} else if minCost == 3 {
			z2[minOrder.id] = 1
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
func BaseSolver(totalCost *int) ([][]uint8, []int, []uint8, []uint8, []uint8, []uint8) {

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
	// Keep for each mover the last order assigned to him
	lastOrders := make([]*Order, nMover)
	for i := 0; i < nMover; i++ {
		lastOrders[i] = new(Order)
		lastOrders[i].id = nOrder + i
	}

	e := orders.Front()
	for i := 0; i < length; i++ {
		order := e.Value.(*Order)

		// Try to assign the order to a mover
		moverId, cost := schedule(order, lastOrders, y, x, w)
		*totalCost += cost

		// If assignment successes
		if moverId >= 0 {
			e = e.Next()
			orders.MoveToBack(e.Prev())
			length--
			i--

			// Update results
			if cost == 1 {
				z[order.id] = 1
			} else if cost == 2 {
				z1[order.id] = 1
			} else if cost == 3 {
				z2[order.id] = 1
			}
		} else {
			e = e.Next()
		}
	}

	return y, x, w, z, z1, z2

}

func schedule(order *Order, lastOrders []*Order, y [][]uint8, x []int, w []uint8, ) (int, int) {

	minCost := utils.Inf
	bestMover := -1
	bestDeliveryTime := 0

	// Find best mover
	for mover := 0; mover < nMover; mover++ {

		cost, deliveryTime := computeCost(lastOrders[mover].id, lastOrders[mover].x, order)
		if cost < minCost {
			minCost = cost
			bestMover = mover
			bestDeliveryTime = deliveryTime
		}
	}

	if bestMover != -1 {
		x[order.id] = bestDeliveryTime
		y[lastOrders[bestMover].id][order.id] = 1

		order.x = bestDeliveryTime
		lastOrders[bestMover] = order
	} else {
		w[order.id] = 1
		minCost = 10
	}

	return bestMover, minCost
}

func main() {
	nOrder = ORDER_N
	nMover = MOVER_N

	distances = utils.CreateOrderMatrix(nOrder, nMover)
	deliveryTimes = utils.CreateDeliveryTimeVector(nOrder)

	utils.PrintDistanceMatrix(distances, nOrder)
	fmt.Print("Algorithm 1:\n")
	start := time.Now()
	var cost int
	y, x, w, z, z1, z2 := GreedySolver(&cost)
	elapsed := time.Since(start)
	//printResults(res)

	utils.PrintAssigmentMatrix(y, nOrder)
	fmt.Println(x)
	fmt.Println(w)
	fmt.Println(z, z1, z2)
	fmt.Printf("Solver took %s\n", elapsed)
	fmt.Printf("Total cost: %d\n", cost)

	fmt.Print("\n\nAlgorithm 2:\n")
	start = time.Now()
	y, x, w, z, z1, z2 = BaseSolver(&cost)
	elapsed = time.Since(start)
	//printResults(res)

	utils.PrintAssigmentMatrix(y, nOrder)
	fmt.Println(x)
	fmt.Println(w)
	fmt.Println(z, z1, z2)
	fmt.Printf("Solver took %s\n", elapsed)
	fmt.Printf("Total cost: %d\n", cost)

}
