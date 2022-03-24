package main

import "math/rand"

type sniper struct {
	position     position
	prevPosition position
	destination  position
	health       int
	movePower    int
	bulletPower  int
}

func newSniper() sniper {
	return sniper{
		position:    wallPosition(),
		destination: randomPosition(),
		health:      40,
		movePower:   2,
		bulletPower: rand.Intn(10),
	}
}

func (s *sniper) avatar() string {
	switch {
	case s.health > 40:
		return "🟪"
	case s.health > 30:
		return "🟦"
	case s.health > 20:
		return "🟩"
	case s.health > 10:
		return "🟨"
	default:
		return "🟥"
	}
}

func (s *sniper) getPosition() position {
	return s.position
}

func (s *sniper) getPrevPosition() position {
	return s.prevPosition
}

func (s *sniper) shouldRemove() bool {
	return s.health == 0
}

func (s *sniper) takeTurn(entities []entity) []entity {
	if s.movePower >= 50 {
		s.movePower = 0
		s.moveTowardDestination()

		if s.hasReachedDestination() {
			s.destination = randomPosition()
		}
	} else {
		s.movePower += s.health
	}

	var bullet bullet

	if s.bulletPower >= 40 {
		s.bulletPower = 0

		wussOut := rand.Intn(2) == 0

		if !wussOut {
			bullet = s.lockAndFire(entities)
		}
	} else {
		s.bulletPower += s.health
	}

	return []entity{&bullet}
}

func (s *sniper) onCollide(e entity) {
	if s.health == 0 {
		return
	}

	s.health--
}

func (s *sniper) hasReachedDestination() bool {
	return positionsAreSame(s.position, s.destination)
}

func (s *sniper) moveTowardDestination() {
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

func (s *sniper) moveUp() {
	s.prevPosition[1] = s.position[1]
	s.position[1]++
}

func (s *sniper) moveDown() {
	s.prevPosition[1] = s.position[1]
	s.position[1]--
}

func (s *sniper) moveRight() {
	s.prevPosition[0] = s.position[0]
	s.position[0]++
}

func (s *sniper) moveLeft() {
	s.prevPosition[0] = s.position[0]
	s.position[0]--
}

func (s *sniper) onRemoveExplode() bool {
	return true
}

func (s *sniper) lockAndFire(entities []entity) bullet {
	var bullet bullet

	ships := getShipsFromEntities(entities)

	if len(ships) == 0 {
		return bullet
	}

	seen := make(map[*ship]struct{})

	for {
		// no one to shoot at
		if len(seen) == len(ships) {
			break
		}

		ship := ships[rand.Intn(len(ships))]

		// already seen
		if _, ok := seen[ship]; ok {
			continue
		}

		if positionsAreSame(s.position, ship.position) {
			seen[ship] = struct{}{}

			continue
		}

		xDis := abs(s.getPosition()[0] - ship.getPosition()[0])
		yDis := abs(s.getPosition()[1] - ship.getPosition()[1])

		// no straight shot
		if xDis != 0 && yDis != 0 && xDis-yDis != 0 {
			seen[ship] = struct{}{}

			continue
		}

		// now there must be a straight shot
		// make bullet and fire
		xPos := 0
		yPos := 0

		if s.getPosition()[0] > ship.getPosition()[0] {
			xPos = -1
		} else if s.getPosition()[0] < ship.getPosition()[0] {
			xPos = 1
		}

		if s.getPosition()[1] > ship.getPosition()[1] {
			yPos = -1
		} else if s.getPosition()[1] < ship.getPosition()[1] {
			yPos = 1
		}

		pos := position{s.position[0] + xPos, s.position[1] + yPos}
		bullet = newBullet(pos, [2]int{xPos, yPos})

		break
	}

	return bullet
}
