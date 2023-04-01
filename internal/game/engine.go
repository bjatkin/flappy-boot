package game

import (
	"github.com/bjatkin/flappy_boot/internal/alloc"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/display"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

type Runable interface {
	Init(*Engine) error
	Update(*Engine, int) error
}

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// activeBackgrounds are the backgrounds that need to be drawn each frame
	activeBackgrounds [4]*Background

	// Allocators
	bgPalAlloc   *alloc.Pal
	sprPalAlloc  *alloc.Pal
	bgTileAlloc  *alloc.VRAM
	sprTileAlloc *alloc.VRAM
	mapAlloc     *alloc.VRAM

	// Debug sprites
	debug [10]*Sprite
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.SetReg(hw_display.Controll, hw_display.Sprite1D|hw_display.ForceBlank)

	e := &Engine{
		activeSprites: make(map[*Sprite]struct{}, 128),
		bgPalAlloc:    alloc.NewPal(memmap.Palette[:256]),
		sprPalAlloc:   alloc.NewPal(memmap.Palette[256:]),

		// the first tile is left transparent and can be shared by all tile maps
		bgTileAlloc:  alloc.NewVRAM(memmap.VRAM[memmap.TileOffset4:memmap.CharBlockOffset*2], 16),
		sprTileAlloc: alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*4:], 16),
		mapAlloc:     alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*2:], memmap.HalfKByte*2),
	}

	debugSprites := [10]*Sprite{}
	for i := range debugSprites {
		debugSprites[i] = e.NewSprite(assets.DebugTileSet)
	}

	e.debug = debugSprites

	return e
}

func (e *Engine) Run(run Runable) error {
	// enable sprites
	memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.Sprites)
	sprite.Reset()

	for {
		err := run.Init(e)
		if err != nil {
			exit(err)
		}

		var frame int
		for {
			key.KeyPoll()
			err := run.Update(e, frame)
			if err != nil {
				exit(err)
			}

			frame++

			vSyncWait()

			// copy active sprite data into OAM memory
			e.drawSprites()

			// copy active background data into the background registers
			e.drawBackgrounds()
		}
	}
}

func (e *Engine) drawSprites() {
	var i int
	for s := range e.activeSprites {
		hw_sprite.OAM[i] = *s.attrs()
		i++
	}

	clear := hw_sprite.Attrs{
		Attr0: hw_sprite.Attr0(255) | hw_sprite.Hide,
		Attr1: hw_sprite.Attr1(511),
	}
	for ; i < 128; i++ {
		hw_sprite.OAM[i] = clear
	}
}

func (e *Engine) drawBackgrounds() {
	// TODO: only update the backgrouds if something has changed?

	// 0xF0FF masks out all the 'active backgrounds' bits from the controll register
	// these bits are then added back to the controll value only if the background is still active
	backgroundControll := memmap.GetReg(hw_display.Controll) & 0xF0FF
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			continue
		}

		switch i {
		case 0:
			backgroundControll |= hw_display.BG0
		case 1:
			backgroundControll |= hw_display.BG1
		case 2:
			backgroundControll |= hw_display.BG2
		case 3:
			backgroundControll |= hw_display.BG3
		}
	}
	memmap.SetReg(hw_display.Controll, backgroundControll)

	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			continue
		}

		controll := e.activeBackgrounds[i].controll()
		switch i {
		case 0:
			memmap.SetReg(hw_display.BG0Controll, controll)
			memmap.SetReg(hw_display.BG0HOffset, e.activeBackgrounds[0].HScroll.Uint16())
			memmap.SetReg(hw_display.BG0VOffset, e.activeBackgrounds[0].VScroll.Uint16())
		case 1:
			memmap.SetReg(hw_display.BG1Controll, controll)
			memmap.SetReg(hw_display.BG1HOffset, e.activeBackgrounds[1].HScroll.Uint16())
			memmap.SetReg(hw_display.BG1VOffset, e.activeBackgrounds[1].VScroll.Uint16())
		case 2:
			memmap.SetReg(hw_display.BG2Controll, controll)
			memmap.SetReg(hw_display.BG2HOffset, e.activeBackgrounds[2].HScroll.Uint16())
			memmap.SetReg(hw_display.BG2VOffset, e.activeBackgrounds[2].VScroll.Uint16())
		case 3:
			memmap.SetReg(hw_display.BG3Controll, controll)
			memmap.SetReg(hw_display.BG3HOffset, e.activeBackgrounds[3].HScroll.Uint16())
			memmap.SetReg(hw_display.BG3VOffset, e.activeBackgrounds[3].VScroll.Uint16())
		}
	}
}

// addBackground adds a new background to the list of active backgrounds.
func (e *Engine) addBackground(bg *Background) error {
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			e.activeBackgrounds[i] = bg
			return nil
		}
	}

	return alloc.ErrOOM
}

// removeBackground removes the current background from the list of active backgrounds. It will not unload the
// background from memory so you must do that yourself if the background is no longer needed
func (e *Engine) removeBackground(bg *Background) {
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == bg {
			e.activeBackgrounds[i] = nil
			return
		}
	}
}

// addSprite adds a new sprite to the list of active sprites.
func (e *Engine) addSprite(sprite *Sprite) {
	e.activeSprites[sprite] = struct{}{}
}

// removeSprite removes a sprite from the list of active sprites. It will not unload the sprites assets
// from memory so you must do that yourself if the sprite is no longer needed
func (e *Engine) removeSprite(sprite *Sprite) {
	delete(e.activeSprites, sprite)
}

// exit exits the game loop and draws error infromation to the screen
func exit(err error) {
	// TODO: this should be updated to use a 'system' font to write out the
	// error data to make debugging easier
	memmap.SetReg(hw_display.Controll, hw_display.Mode3|hw_display.BG2)

	// Draw red to the screen so we can tell there was an error
	blue := display.RGB15(0, 0, 31)
	for i := 10; i < 240*160; i++ {
		memmap.VRAM[i] = memmap.VRAMValue(blue)
	}

	// block forever
	for {
	}
}

// vSyncWait blocks while it waits for the screen to enter the vertical blank and then returns
func vSyncWait() {
	// TODO: leverage hardware interrupts rather than spinning the GBA cpu like this.
	// my guess is this is going to lead to pretty high power usage for no real benefit

	// wait till VDraw
	for display.VCount() >= 160 {
	}

	// wailt tile VBlank
	for display.VCount() < 160 {
	}
}
