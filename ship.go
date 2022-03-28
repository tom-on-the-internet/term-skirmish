package main

import "math/rand"

type ship struct {
	position     position
	prevPosition position
	destination  position
	alive        bool
	movePower    int
	bulletPower  int
	team         team
}

type team int

const (
	BLUE = iota
	RED
	YELLOW
	GREEN
	ORANGE
	BROWN
	PURPLE
	WHITE
)

func (d team) Ok() bool {
	switch d {
	case BLUE,
		BROWN,
		GREEN,
		ORANGE,
		PURPLE,
		RED,
		WHITE,
		YELLOW:
		return true
	}

	return false
}

func newShip() ship {
	wallPos := wallPosition()

	return ship{
		position:     wallPos,
		prevPosition: wallPos,
		destination:  randomPosition(),
		alive:        true,
		movePower:    3,
		bulletPower:  rand.Intn(10),
		team:         team(rand.Intn(2)),
	}
}

func (s *ship) avatar() string {
	switch s.team {
	case BLUE:
		return "🔵"
	case BROWN:
		return "🟤"
	case GREEN:
		return "🟢"
	case ORANGE:
		return "🟠"
	case PURPLE:
		return "🟣"
	case RED:
		return "🔴"
	case WHITE:
		return "⚪"
	case YELLOW:
		return "🟡"
	}

	// indicates error
	return "❌"
}

func (s *ship) getPosition() position {
	return s.position
}

func (s *ship) getPrevPosition() position {
	return s.prevPosition
}

func (s *ship) shouldRemove() bool {
	return !s.alive
}

func (s *ship) takeTurn(entities []entity) []entity {
	s.move(entities)

	bullet := s.shoot(entities)

	if bullet == nil {
		return []entity{}
	}

	return []entity{bullet}
}

func (s *ship) shoot(entities []entity) *bullet {
	if s.bulletPower != 15 {
		s.bulletPower++

		return nil
	}

	s.bulletPower = 0

	if wussOut := rand.Intn(2) == 0; wussOut {
		return nil
	}

	ships := getShipsFromEntities(entities)

	if len(ships) == 0 {
		return nil
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

		// same team
		if s.team == ship.team {
			seen[ship] = struct{}{}

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
		bullet := newBullet(pos, [2]int{xPos, yPos})

		return &bullet
	}

	return nil
}

func (s *ship) move(entities []entity) {
	if s.movePower != 3 {
		s.movePower++

		return
	}

	s.movePower = 0
	s.moveTowardDestination()

	if s.hasReachedDestination() {
		s.destination = s.getDestination(entities)
	}
}

func (s *ship) getDestination(entities []entity) position {
	if rand.Intn(2) == 0 {
		return randomPosition()
	}

	for _, e := range entities {
		ship, ok := e.(*ship)
		if ok && ship.team != s.team {
			return ship.getPosition()
		}
	}

	return randomPosition()
}

func (s *ship) onCollide(e entity) {
	s.alive = false
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
