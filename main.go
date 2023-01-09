package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/gameplay"
	"github.com/bjatkin/flappy_boot/internal/game"
)

//go:embed assets/gba
var assets embed.FS

func main() {
	game.Run(gameplay.NewDemo(assets))
}
