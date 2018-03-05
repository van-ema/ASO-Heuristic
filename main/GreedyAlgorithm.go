package main

import (
	"orderSchedulingAlgorithm/main/utils"
	"fmt"
	"container/list"
	"time"
	"log"
)

const (
	ORDER_N = 300 // #Orders
	MOVER_N = 30  // #Movers
)

type Order struct {
	id int // order's id
	x  int // target delivery time
	t  int // delivery time
}

var nOrder int
var nMover int

// return for each mover a list of orders
func GreedySolver(distances [][]int, deliveryTimes []int) ([][]*Order, int) {

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
	// Solving for each mover
	for mover := 0; mover < nMover; mover++ {

		// Orders are never removed from the orders list but they are moved on the back
		// We use the var length to keep the number of orders not already scheduled
		// placed at the beginning of the list
		totalCost += SingleMoverSchedulingOrders(mover, &results[mover], orders, &length, distances)

	}

	// The orders still in the list are all those that cannot be scheduled
	var cancelled []*Order
	e := orders.Front()
	for k := 0; k < length; k++ {
		cancelled = append(cancelled, e.Value.(*Order))
		e = e.Next()

	}
	results[nMover] = cancelled

	return results, totalCost
}

// return a list of orders that can be scheduled by the mover
func SingleMoverSchedulingOrders(mover int, result *[]*Order, orders *list.List, length *int, distances [][]int) int {

	cost := 0
	for i := 0; i < *length; i++ {

		var minOrderElem *list.Element
		minCost := utils.Inf
		newDeliveryTime := 0

		current := orders.Front()
		for j := 0; j < *length; j++ {
			order := current.Value.(*Order)

			// Get last order in the list assigned to the mover to compute distance with
			// the next order
			// If the mover has no order assigned the distance must be computed
			// with the mover initial position
			var lastOrder int
			var lastDeliveryTime int
			if len(*result) == 0 {
				lastOrder = nOrder + mover
				lastDeliveryTime = 0
			} else {
				lastOrder = (*result)[len(*result)-1].id
				lastDeliveryTime = (*result)[len(*result)-1].x
			}

			newCost, nextDeliveryTime := computeCost(lastOrder, lastDeliveryTime, order, distances)
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
		minOrder.x = newDeliveryTime
		*result = append(*result, minOrder)
		orders.MoveToBack(minOrderElem)

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

func computeCost(lastOrderId int, lastDeliveryTime int, nextOrder *Order, distances [][]int) (int, int) {

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
		if k >= 0 {
			fmt.Printf("%s%d : ", "Mover-", k)
		} else {
			fmt.Printf("cancelled:")
		}

		printArray(results[k])

		fmt.Print("\n")
	}
}

func printArray(p []*Order) {
	for _, v := range p {
		fmt.Printf("[id:%d,t:%d,x:%d]", v.id, v.t, v.x)
	}
}

func main() {
	nOrder = ORDER_N
	nMover = MOVER_N

	distances := utils.CreateOrderMatrix(nOrder, nMover)
	t := utils.CreateDeliveryTimeVector(nOrder)

	start := time.Now()
	res, cost := GreedySolver(distances, t)
	elapsed := time.Since(start)
	printResults(res)
	log.Printf("Solver took %s", elapsed)
	log.Printf("Total cost: %d", cost)
}
