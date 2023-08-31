package main

type Case struct {
	X     int
	Y     int
	State string // "empty", "hit", "flop"
}

const GameSize = 10

var PlayerGrid [GameSize][GameSize]Case

func initGrid(grid *[GameSize][GameSize]Case) {
	for x := 0; x < GameSize; x++ {
		for y := 0; y < GameSize; y++ {
			grid[x][y] = Case{X: x, Y: y, State: "empty"}
		}
	}
}
