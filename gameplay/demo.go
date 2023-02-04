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

const (
	// gravity influences how quickly the boot falls
	gravity = 20

	// maxDy is the terminal velocity for the boot
	maxDy = 512

	// grassY is the y level that the grass ground sits at
	grassY = 144
)

// Demo is a test node used for prototyping basic mechanics
type Demo struct {
	assets embed.FS
	dy     int32
	booty  int32
}

func NewDemo(assets embed.FS) *Demo {
	return &Demo{
		assets: assets,
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

	// TODO: this should be managed through a sprite package
	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | 72,
		Attr1: hw_sprite.Medium | hw_sprite.Attr1(20),
	}

	mode0.Enable(
		mode0.WithBG(true, true, false, false),
		mode0.With1DSprites(),
	)

	return nil
}

func (d *Demo) Update(frame uint) (game.Node, error) {
	if key.Pressed(key.A) {
		d.dy -= 10
	}
	if d.dy < maxDy {
		d.dy += gravity
	}

	d.booty += d.dy

	if d.booty > grassY<<8 {
		d.booty = grassY << 8
	}

	return nil, nil
}

func (d *Demo) Draw() error {
	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | hw_sprite.Attr0(d.booty>>8),
		Attr1: hw_sprite.Medium | 8,
	}
	return nil
}

func (d *Demo) Unload() error {
	return nil
}
