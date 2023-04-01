package main

import (
	"github.com/bjatkin/flappy_boot/gameplay"
	"github.com/bjatkin/flappy_boot/internal/game"
)

func main() {
	engine := game.NewEngine()
	engine.Run(gameplay.NewManager(engine))
}
