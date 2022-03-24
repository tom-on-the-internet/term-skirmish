package main

import (
	"math/rand"
	"strconv"
)

type (
	// position x,y.
	position [2]int
	entity   interface {
		getPosition() position
		getPrevPosition() position
		shouldRemove() bool
		avatar() string
		takeTurn(entities []entity) []entity
		onCollide(e entity)
		onRemoveExplode() bool
	}
)

func collided(a, b entity) bool {
	posA, posB := a.getPosition(), b.getPosition()

	if positionsAreSame(posA, posB) {
		return true
	}

	prevPosA, prevPosB := a.getPrevPosition(), b.getPrevPosition()

	// swapped position
	if positionsAreSame(posA, prevPosB) && positionsAreSame(posB, prevPosA) {
		return true
	}

	return false
}

func positionsAreSame(a, b position) bool {
	return a[0] == b[0] && a[1] == b[1]
}

func randomPosition() [2]int {
	width, height := getSize()
	x := rand.Intn(width-1) + 1
	y := rand.Intn(height-1) + 2

	return [2]int{x, y}
}

func countShips(entities []entity) int {
	numShips := 0

	for _, e := range entities {
		if _, ok := e.(*ship); ok {
			numShips++
		}
	}

	return numShips
}

func getStatus(entities []entity) string {
	x, y := getSize()
	return "ship count: " + strconv.Itoa(
		countShips(entities),
	) + " x: " + strconv.Itoa(
		x,
	) + " y: " + strconv.Itoa(
		y,
	)
}
