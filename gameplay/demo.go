package gameplay

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

// Demo is a test node used for prototyping basic mechanics
type Demo struct {
	assets embed.FS
}

func NewDemo(assets embed.FS) *Demo {
	return &Demo{
		assets: assets,
	}
}

func (d *Demo) Init() error {
	// Load in the sprite palettes
	sprite.LoadPalette16(d.assets, "palette_0.p16", 0)
	sprite.LoadPalette16(d.assets, "palette_1.p16", 1)
	sprite.LoadPalette16(d.assets, "palette_2.p16", 2)
	sprite.LoadPalette16(d.assets, "palette_3.p16", 3)

	return nil
}

func (d *Demo) Update(frame uint) (game.Node, error) {
	return nil, nil
}

func (d *Demo) Draw() error {
	return nil
}

func (d *Demo) Unload() error {
	return nil
}
