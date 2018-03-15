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


import "orderSchedulingAlgorithm/utils"

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
	orders := initOrder(deliveryTimes, nOrder)
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

		cost, deliveryTime ,_:= computeCost(lastOrders[mover].id, lastOrders[mover].x, order)
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