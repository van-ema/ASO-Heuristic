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

	CANCELLED = -1
)

type Order struct {
	id int // order's id
	x  int // target delivery time
	t  int // delivery time
}

// return for each mover a list of orders
func GreedySolver(nOrder int, nMover int, distances [][]int, deliveryTimes []int) map[int][]*Order {

	// alias for D - D* : contains the orders not assigned
	orders := initOrder(deliveryTimes)
	length := orders.Len()
	// alias for  D* : contains the orders already assigned
	assigned := make([]*Order, nOrder)
	// key: mover ; value: array of orders assigned to him
	results := make(map[int][]*Order)

	totalCost := 0
	// Solving for each mover
	for mover := 0; mover < nMover; mover++ {

		// array of orders assigned to the current mover
		partition := make([]*Order, 0, int(nOrder/nMover)+1)
		for i := 0; i < length; i++ {

			var minOrderElem *list.Element
			minCost := utils.Inf
			newDeliveryTime := 0

			current := orders.Front()
			for j := 0; j < length; j++ {
				order := current.Value.(*Order)

				// Get last order in the list assigned to the mover to compute distance with
				// the next order
				// If the mover has no order assigned the distance must be computed
				// with the mover initial position
				var lastOrder int
				var lastDeliveryTime int
				if len(partition) == 0 {
					lastOrder = nOrder + mover
					lastDeliveryTime = 0
				} else {
					lastOrder = partition[len(partition)-1].id
					lastDeliveryTime = partition[len(partition)-1].x
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
			totalCost += minCost
			minOrder := minOrderElem.Value.(*Order)
			minOrder.x = newDeliveryTime
			partition = append(partition, minOrder)
			assigned = append(assigned, minOrder)
			orders.MoveToBack(minOrderElem)

			length--
			i--

		}

		results[mover] = partition
	}

	var cancelled []*Order
	e := orders.Front()
	for k := 0; k < length; k++ {
		cancelled = append(cancelled, e.Value.(*Order))
		e = e.Next()

	}

	results[CANCELLED] = cancelled

	return results
}

// orders allocation and initialization with delivery times
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

func printResults(results map[int][]*Order) {
	for k := range results {
		if k >= 0 {
			fmt.Printf("%s%d : ", "Mover-", k)
		} else {
			fmt.Printf("cancelled:")
		}

		for _, v := range results[k] {
			fmt.Printf("[id:%d,t:%d,x:%d]", v.id, v.t, v.x)
		}

		fmt.Print("\n")

	}
}

func main() {
	m1 := utils.CreateOrderMatrix(ORDER_N, MOVER_N)
	t := utils.CreateDeliveryTimeVector(ORDER_N)
	start := time.Now()
	res := GreedySolver(ORDER_N, MOVER_N, m1, t)
	elapsed := time.Since(start)
	log.Printf("Solver took %s", elapsed)
	printResults(res)
}
