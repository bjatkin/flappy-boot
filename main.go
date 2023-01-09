package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/gameplay"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

//go:embed assets/gba
var assets embed.FS

func main() {
	err := sprite.LoadPalette16(assets, "assets/gba/palette_0.p16", 0)
	if err != nil {
		return
	}

	game.Run(gameplay.NewDemo(assets))
}
