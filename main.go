package main

import (
	"math/rand"
	"time"
)

func main() {
	beforeGame()
	runGame()
	gameOver()
}

func beforeGame() {
	rand.Seed(time.Now().UnixNano())
	hideCursor()
}

func runGame() {
	entities := []entity{}
	numShips := 50

	for i := 0; i < numShips; i++ {
		ship := newShip()
		entities = append(entities, &ship)
	}

	for numShips > 1 {
		clear()

		newEntities := []entity{}

		for _, entity := range entities {
			es := entity.takeTurn(entities)
			if newEntities != nil {
				newEntities = append(newEntities, es...)
			}
		}

		collidedEntities := make(map[entity]struct{})

		for _, entity := range entities {
			if _, collided := collidedEntities[entity]; collided {
				continue
			}

			for _, otherEntity := range entities {
				if entity == otherEntity {
					continue
				}

				if collided(entity, otherEntity) {
					entity.onCollide(otherEntity)

					collidedEntities[entity] = struct{}{}
				}
			}
		}

		remainingEntities := []entity{}

		for _, entity := range entities {
			if !entity.shouldRemove() {
				remainingEntities = append(remainingEntities, entity)

				continue
			}

			if entity.onRemoveExplode() {
				explosion := newExplosion(entity.getPosition())
				remainingEntities = append(remainingEntities, &explosion)
			}
		}

		entities = remainingEntities

		drawGame(entities)

		time.Sleep(time.Second / 60)

		entities = append(entities, newEntities...)

		if rand.Intn(100) == 99 {
			ship := newShip()

			entities = append(entities, &ship)
		}
		numShips = countShips(entities)
	}
}

func drawGame(entities []entity) {
	// draw
	for _, entity := range entities {
		moveCursor(entity.getPosition())
		draw(entity.avatar())
	}

	width, _ := getSize()

	status := getStatus(entities)
	moveCursor(position{width/2 - (len(status) / 2), 0})
	draw(status)

	render()
}

func gameOver() {
	clear()
	showCursor()
	moveCursor(position{0, 0})
	draw("Thanks for watching.")
}
