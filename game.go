package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type game struct {
	over            bool
	numTeams        int
	maxShipsPerWave int
	entities        []entity
	shipCount       int
}

func newGame() game {
	const (
		defaultTeamCount    = 2
		defaultMaxWaveCount = 8
	)

	var (
		numTeams        int
		maxShipsPerWave int
	)

	flag.IntVar(&numTeams, "teams", defaultTeamCount, "number of teams (1-8)")
	flag.IntVar(
		&maxShipsPerWave,
		"wave",
		defaultMaxWaveCount,
		"maximum number of ships in each reinforcement wave (1-100)",
	)

	flag.Parse()

	if numTeams < 1 || numTeams > 8 {
		numTeams = 2
	}

	if maxShipsPerWave < 1 || maxShipsPerWave > 100 {
		maxShipsPerWave = 20
	}

	return game{
		over:            false,
		numTeams:        numTeams,
		maxShipsPerWave: maxShipsPerWave,
		entities:        []entity{},
		shipCount:       0,
	}
}

func (g *game) beforeGame() {
	rand.Seed(time.Now().UnixNano())

	hideCursor()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			g.over = true
		}
	}()

	// initial wave
	g.reinforce()
}

func (g *game) reinforce() {
	team := team(rand.Intn(g.numTeams))
	shipCount := rand.Intn(g.maxShipsPerWave) + 1

	for i := 0; i < shipCount; i++ {
		ship := newShip(team)
		g.entities = append(g.entities, &ship)

		g.shipCount++
	}
}

func (g *game) runGame() {
	for !g.over {
		clear()

		newEntities := g.takeTurns()

		g.checkCollisions()
		g.removeEntities()
		g.drawGame()

		// 60 fps
		time.Sleep(time.Second / 60)

		g.entities = append(g.entities, newEntities...)

		// 0.5% chance of reinforcements
		if rand.Intn(200) == 0 {
			g.reinforce()
		}
	}
}

func (g *game) takeTurns() []entity {
	newEntities := []entity{}

	for _, entity := range g.entities {
		es := entity.takeTurn(g.entities)
		if newEntities != nil {
			newEntities = append(newEntities, es...)
		}
	}

	return newEntities
}

func (g *game) checkCollisions() {
	collidedEntities := make(map[entity]struct{})

	for _, entity := range g.entities {
		if _, collided := collidedEntities[entity]; collided {
			continue
		}

		for _, otherEntity := range g.entities {
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

func (g *game) removeEntities() {
	remainingEntities := []entity{}

	for _, entity := range g.entities {
		if !entity.shouldRemove() {
			remainingEntities = append(remainingEntities, entity)

			continue
		}

		if entity.onRemoveExplode() {
			explosion := newExplosion(entity.getPosition())
			remainingEntities = append(remainingEntities, &explosion)
		}
	}

	g.entities = remainingEntities
}

func (g *game) drawGame() {
	for _, entity := range g.entities {
		moveCursor(entity.getPosition())
		draw(entity.avatar())
	}

	width, _ := getSize()

	status := g.getStatus()
	moveCursor(position{width/2 - (len(status) / 2), 0})

	draw(status)

	render()
}

func (g *game) getStatus() string {
	currentShipCount := countShips(g.entities)

	message := "current ship count: " + strconv.Itoa(
		countShips(g.entities),
	) + "     destroyed count: " + strconv.Itoa(
		g.shipCount-currentShipCount,
	)

	return message
}

func (g *game) afterGame() {
	clear()
	showCursor()

	moveCursor(position{1, 1})
	draw(" 🔵  🔸 🔸 🔥")

	moveCursor(position{0, 3})

	username, _ := os.LookupEnv("USER")
	if username == "" {
		username = "friend"
	}

	draw(fmt.Sprintf("See you again soon %v!", username))

	moveCursor(position{0, 5})

	render()
	os.Exit(0)
}
