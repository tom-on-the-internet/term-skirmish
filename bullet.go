package main

type bullet struct {
	position     position
	prevPosition position
	direction    [2]int
	active       bool
	bulletPower  int
}

func newBullet(pos position, direction [2]int) bullet {
	return bullet{
		position:     pos,
		prevPosition: pos,
		direction:    direction,
		active:       true,
		bulletPower:  1,
	}
}

func (b *bullet) avatar() string {
	return "🔸"
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
	if b.bulletPower == 1 {
		b.move()
		b.bulletPower = 0

		return nil
	}

	b.bulletPower = 1

	return nil
}

func (b *bullet) onCollide(e entity) {
	b.active = false
}

func (b *bullet) move() {
	b.prevPosition = b.getPosition()
	xPos := b.position[0] + b.direction[0]
	yPos := b.position[1] + b.direction[1]

	maxX, maxY := getSize()
	if xPos < 1 || xPos > maxX || yPos < 1 || yPos > maxY {
		b.active = false

		return
	}

	b.position = position{xPos, yPos}
}

func (b *bullet) onRemoveExplode() bool {
	return false
}
