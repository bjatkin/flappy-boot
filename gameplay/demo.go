package gameplay

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/mode0"
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
	mode0.Enable(
		mode0.WithBG(true, false, false, false),
		mode0.With1DSprites(),
	)
	memmap.SetReg(hw_display.BG0Controll, 1<<hw_display.SBBShift|hw_display.Priority3)

	tile := assets.NewBackground()
	tile.LoadMap(64)

	player := assets.NewPlayer()
	player.Load()

	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | 0xA,
		Attr1: hw_sprite.Medium | 0xA,
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
