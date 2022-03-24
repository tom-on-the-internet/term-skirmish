package main

type explosion struct {
	position position
	health   int
}

func newExplosion(pos position) explosion {
	return explosion{
		position: pos,
		health:   100,
	}
}

func (e *explosion) getPosition() position {
	return e.position
}

func (e *explosion) getPrevPosition() position {
	return e.position
}

func (e *explosion) shouldRemove() bool {
	return e.health == 0
}

func (e *explosion) avatar() string {
	return "💥"
}

func (e *explosion) takeTurn(entities []entity) []entity {
	if e.health > 0 {
		e.health--
	}

	return nil
}

func (e *explosion) onCollide(otherEntity entity) {
	if _, collidedWithShip := otherEntity.(*ship); collidedWithShip {
		e.health = 0

		return
	}

	e.health += 10
}

func (e *explosion) onRemoveExplode() bool {
	return false
}
