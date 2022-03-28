package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

var gameOver = false

func main() {
	beforeGame()
	runGame()
}

func beforeGame() {
	rand.Seed(time.Now().UnixNano())
	hideCursor()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			gameOver = true
		}
	}()
}

func runGame() {
	shipCount := 0
	entities := []entity{}

	reinforce := func() {
		for i := 0; i < rand.Intn(4); i++ {
			ship := newShip()
			entities = append(entities, &ship)

			shipCount++
		}
	}

	for !gameOver {
		clear()

		newEntities := takeTurns(entities)

		checkCollisions(entities)

		entities = removeEntities(entities)

		drawGame(entities, shipCount)

		// 60 fps
		time.Sleep(time.Second / 60)

		entities = append(entities, newEntities...)

		if rand.Intn(100) == 0 {
			reinforce()
		}

	}

	onGameOver()
}

func takeTurns(entities []entity) []entity {
	newEntities := []entity{}

	for _, entity := range entities {
		es := entity.takeTurn(entities)
		if newEntities != nil {
			newEntities = append(newEntities, es...)
		}
	}

	return newEntities
}

func checkCollisions(entities []entity) {
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
}

func removeEntities(entities []entity) []entity {
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

	return remainingEntities
}

func drawGame(entities []entity, shipCount int) {
	for _, entity := range entities {
		moveCursor(entity.getPosition())
		draw(entity.avatar())
	}

	width, _ := getSize()

	status := getStatus(entities, shipCount)
	moveCursor(position{width/2 - (len(status) / 2), 0})

	draw(status)

	render()
}

func onGameOver() {
	clear()
	showCursor()

	moveCursor(position{1, 1})
	draw("🟦 🟢 🟡 🔴 ⭕")

	moveCursor(position{0, 3})

	username, _ := os.LookupEnv("USER")
	draw(fmt.Sprintf("See you again soon %v!", username))

	moveCursor(position{0, 5})

	render()
	os.Exit(0)
}

func allSameTeam(entities []entity) bool {
	seenTeams := make(map[team]struct{})
	ships := getShipsFromEntities(entities)

	for _, ship := range ships {
		seenTeams[ship.team] = struct{}{}
		if len(seenTeams) == 2 {
			return false
		}
	}

	return true
}
