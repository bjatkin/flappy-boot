package main

import (
	"github.com/bjatkin/flappy_boot/gameplay/fly"
	"github.com/bjatkin/flappy_boot/internal/game"
)

func main() {
	engine := game.NewEngine()
	engine.Run(fly.NewStage())
}
