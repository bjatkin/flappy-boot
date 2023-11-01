package main

import (
	"github.com/bjatkin/flappy_boot/gameplay"
	"github.com/bjatkin/flappy_boot/internal/game"
)

func main() {
	harness := game.NewHarness()
	harness.Run(gameplay.NewManager(harness.E))
}
