package main

import (
	"math/rand"
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

func collided(entityA, entityB entity) bool {
	posA, posB := entityA.getPosition(), entityB.getPosition()

	if positionsAreSame(posA, posB) {
		return true
	}

	prevPosA, prevPosB := entityA.getPrevPosition(), entityB.getPrevPosition()

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

func wallPosition() [2]int {
	maxX, maxY := getSize()

	switch rand.Intn(4) {
	case 0:
		// top
		return [2]int{rand.Intn(maxX), 1}
	case 1:
		// bottom
		return [2]int{rand.Intn(maxX), maxY}
	case 2:
		// left
		return [2]int{1, rand.Intn(maxY)}
	default:
		// right
		return [2]int{maxX, rand.Intn(maxY)}
	}
}

func countShips(entities []entity) int {
	return len(getShipsFromEntities(entities))
}

func getShipsFromEntities(entities []entity) []*ship {
	ships := []*ship{}

	for _, e := range entities {
		if ship, ok := e.(*ship); ok {
			ships = append(ships, ship)
		}
	}

	return ships
}

func abs(i int) int {
	if i < 0 {
		i *= -1
	}

	return i
}
