package gameplay

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
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
	test := assets.NewBG()
	for i := range test.Palette {
		memmap.Palette[i] = memmap.PaletteValue(test.Palette[i])
	}

	for i := range test.Tiles {
		memmap.VRAM[i] = memmap.VRAMValue(test.Tiles[i])
	}

	// Load in the sprite palettes
	err := sprite.LoadPalette16(d.assets, "assets/gba/palette_0.p16", 0)
	if err != nil {
		return err
	}

	err = sprite.LoadPalette16(d.assets, "assets/gba/palette_1.p16", 1)
	if err != nil {
		return err
	}

	err = sprite.LoadPalette16(d.assets, "assets/gba/palette_2.p16", 2)
	if err != nil {
		return err
	}

	err = sprite.LoadPalette16(d.assets, "assets/gba/palette_3.p16", 3)
	if err != nil {
		return err
	}

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
