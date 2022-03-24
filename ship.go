package main

import "math/rand"

type ship struct {
	position     position
	prevPosition position
	destination  position
	health       int
	movePower    int
	bulletPower  int
}

func newShip() ship {
	wallPos := wallPosition()

	return ship{
		position:     wallPos,
		prevPosition: wallPos,
		destination:  randomPosition(),
		health:       4,
		movePower:    6,
		bulletPower:  rand.Intn(10),
	}
}

func (s *ship) avatar() string {
	switch s.health {
	case 4:
		return "🟢"
	case 3:
		return "🟡"
	case 2:
		return "🔴"
	default:
		return "⭕"
	}
}

func (s *ship) getPosition() position {
	return s.position
}

func (s *ship) getPrevPosition() position {
	return s.prevPosition
}

func (s *ship) shouldRemove() bool {
	return s.health == 0
}

func (s *ship) takeTurn(entities []entity) []entity {
	if s.movePower >= 12 {
		s.movePower = 0
		s.moveTowardDestination()

		if s.hasReachedDestination() {
			s.destination = randomPosition()
		}
	} else {
		s.movePower += s.health
	}

	var bullet bullet

	if s.bulletPower >= 100 {
		s.bulletPower = 0

		x := rand.Intn(3) - 1
		y := rand.Intn(3) - 1
		wussOut := rand.Intn(2)

		if wussOut == 0 && (x != 0 || y != 0) {
			pos := position{s.position[0] + x, s.position[1] + y}
			bullet = newBullet(pos, [2]int{x, y})
		}

	} else {
		s.bulletPower += s.health
	}

	// try to shoot at target
	// can shoot at a target if it is in a direction

	return []entity{&bullet}
}

func (s *ship) onCollide(e entity) {
	if s.health == 0 {
		return
	}

	s.health--
}

func (s *ship) hasReachedDestination() bool {
	return positionsAreSame(s.position, s.destination)
}

func (s *ship) moveTowardDestination() {
	if s.position[0] < s.destination[0] {
		s.moveRight()
	} else if s.position[0] > s.destination[0] {
		s.moveLeft()
	}

	if s.position[1] < s.destination[1] {
		s.moveUp()
	} else if s.position[1] > s.destination[1] {
		s.moveDown()
	}
}

func (s *ship) moveUp() {
	s.prevPosition[1] = s.position[1]
	s.position[1]++
}

func (s *ship) moveDown() {
	s.prevPosition[1] = s.position[1]
	s.position[1]--
}

func (s *ship) moveRight() {
	s.prevPosition[0] = s.position[0]
	s.position[0]++
}

func (s *ship) moveLeft() {
	s.prevPosition[0] = s.position[0]
	s.position[0]--
}

func (s *ship) onRemoveExplode() bool {
	return true
}
