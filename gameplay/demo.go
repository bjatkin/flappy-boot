package gameplay

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/mode0"
)

// Demo is a test node used for prototyping basic mechanics
type Demo struct {
	assets embed.FS
	tmp    uint
}

func NewDemo(assets embed.FS) *Demo {
	return &Demo{
		assets: assets,
		tmp:    0xA,
	}
}

func (d *Demo) Init() error {
	// force a blank screen while we're loading all this data in
	memmap.SetReg(hw_display.Controll, hw_display.ForceBlank)

	// Load in the basic background
	screenBaseBlock := memmap.BGControll(2)
	charBase := memmap.BGControll(0)
	memmap.SetReg(hw_display.BG0Controll, charBase<<hw_display.CBBShift|screenBaseBlock<<hw_display.SBBShift|hw_display.Priority2|hw_display.BGSizeWide)

	tile := assets.NewBackground()
	tile.LoadMap(0, charBase, screenBaseBlock, 64)

	screenBaseBlock += 2
	charBase += 1
	memmap.SetReg(hw_display.BG1Controll, charBase<<hw_display.CBBShift|screenBaseBlock<<hw_display.SBBShift|hw_display.Priority2|hw_display.BGSizeWide)

	sky := assets.NewSky()
	sky.LoadMap(1, charBase, screenBaseBlock, 64)

	player := assets.NewPlayer()
	player.Load()

	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | 0xA,
		Attr1: hw_sprite.Medium | hw_sprite.Attr1(d.tmp),
	}

	mode0.Enable(
		mode0.WithBG(true, true, false, false),
		mode0.With1DSprites(),
	)

	return nil
}

func (d *Demo) Update(frame uint) (game.Node, error) {
	if key.PressedDown(key.Right) {
		d.tmp += 256
	}
	d.tmp += 64
	d.tmp %= 256 << 8

	return nil, nil
}

func (d *Demo) Draw() error {
	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | 0xA,
		Attr1: hw_sprite.Medium | hw_sprite.Attr1(d.tmp>>8),
	}
	return nil
}

func (d *Demo) Unload() error {
	return nil
}
