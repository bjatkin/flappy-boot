package gameplay

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/mode0"
)

var (
	// gravity influences how quickly the boot falls
	gravity = fix.New(0, 64)

	// jump influences how powerful each jump is
	jump = fix.New(-4, 0)

	// grassY is the y level that the grass ground sits at
	grassY = fix.New(131, 0)
)

// Demo is a test node used for prototyping basic mechanics
type Demo struct {
	assets embed.FS
	bootDy fix.P8
	bootY  fix.P8
	scroll int
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
	player.Load(0, 0)

	// TODO: this should be managed through a sprite package
	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | 72,
		Attr1: hw_sprite.Medium | 40,
	}

	pillars := assets.NewPillars()
	pillars.Load(1, 1)

	// TOP
	hw_sprite.OAM[1] = hw_sprite.Attrs{
		Attr0: hw_sprite.Wide | hw_sprite.Color16 | hw_sprite.Normal | 50,
		Attr1: hw_sprite.Medium | 50,
		Attr2: 0x0200 | 0x1<<0xC,
	}

	// BOTTOM
	hw_sprite.OAM[2] = hw_sprite.Attrs{
		Attr0: hw_sprite.Wide | hw_sprite.Color16 | hw_sprite.Normal | 80,
		Attr1: hw_sprite.Medium | 50,
		Attr2: 0x0204 | 0x1<<0xC,
	}

	// MIDLE
	hw_sprite.OAM[3] = hw_sprite.Attrs{
		Attr0: hw_sprite.Wide | hw_sprite.Color16 | hw_sprite.Normal | 50,
		Attr1: hw_sprite.Medium | 80,
		Attr2: 0x0208 | 0x1<<0xC,
	}

	mode0.Enable(
		mode0.WithBG(true, true, false, false),
		mode0.With1DSprites(),
	)

	return nil
}

func (d *Demo) Update(frame uint) (game.Node, error) {
	if key.IsPressed(key.A) {
		d.bootDy = jump
	}

	d.bootDy += gravity
	d.bootY += d.bootDy

	if d.bootY > grassY {
		d.bootY = grassY
		d.bootDy = 0
	}

	if d.bootY < 0 {
		d.bootY = 0
		d.bootDy = 0
	}

	// scroll the background
	d.scroll++
	memmap.SetReg(display.BG0HOffset, uint16(d.scroll))
	memmap.SetReg(display.BG1HOffset, uint16(d.scroll>>2))

	return nil, nil
}

func (d *Demo) Draw() error {
	// TODO: this shouldn't be nessisary since the game.Run function should just copy the sprite data automatically
	hw_sprite.OAM[0] = hw_sprite.Attrs{
		Attr0: hw_sprite.Square | hw_sprite.Color16 | hw_sprite.Normal | hw_sprite.Attr0(d.bootY.Int()),
		Attr1: hw_sprite.Medium | 40,
	}
	return nil
}

func (d *Demo) Unload() error {
	return nil
}
