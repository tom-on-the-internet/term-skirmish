package main

type bullet struct {
	position     position
	prevPosition position
	direction    [2]int
	active       bool
}

func newBullet(pos position, direction [2]int) bullet {
	return bullet{
		position:  pos,
		direction: direction,
		active:    true,
	}
}

func (b *bullet) avatar() string {
	return "🔹"
}

func (b *bullet) getPosition() position {
	return b.position
}

func (b *bullet) getPrevPosition() position {
	return b.prevPosition
}

func (b *bullet) shouldRemove() bool {
	return !b.active
}

func (b *bullet) takeTurn(entities []entity) []entity {
	b.move()

	return nil
}

func (b *bullet) onCollide(e entity) {
	b.active = false
}

func (b *bullet) move() {
	x := b.position[0] + b.direction[0]
	y := b.position[1] + b.direction[1]

	maxX, maxY := getSize()
	if x < 1 || x > maxX || y < 1 || y > maxY {
		b.active = false

		return
	}

	b.position = position{x, y}
}

func (b *bullet) onRemoveExplode() bool {
	return false
}
