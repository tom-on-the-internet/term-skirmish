package main

func main() {
	run()
}

func run() {
	g := newGame()
	g.beforeGame()
	g.runGame()
	g.afterGame()
}
