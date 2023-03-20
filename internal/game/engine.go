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
	Next() (Runable, bool)
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
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.SetReg(hw_display.Controll, hw_display.Sprite1D|hw_display.ForceBlank)
	return &Engine{
		activeSprites: make(map[*Sprite]struct{}, 128),
		bgPalAlloc:    alloc.NewPal(memmap.Palette[:256]),
		sprPalAlloc:   alloc.NewPal(memmap.Palette[256:]),

		// the first tile is left transparent and can be shared by all tile maps
		bgTileAlloc:  alloc.NewVRAM(memmap.VRAM[memmap.TileOffset4:memmap.CharBlockOffset*2], 16),
		sprTileAlloc: alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*4:], 16),
		mapAlloc:     alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*2:], memmap.HalfKByte*2),
	}
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

			if next, ok := run.Next(); ok {
				sprite.Reset()
				run = next
				break
			}
		}
	}
}

func (e *Engine) drawSprites() {
	var i int
	for s := range e.activeSprites {
		hw_sprite.OAM[i] = *s.attrs()
		i++
	}
}

func (e *Engine) drawBackgrounds() {
	// TODO: only update the backgrouds if something has changed?
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			continue
		}

		controll := e.activeBackgrounds[i].controll()
		switch i {
		case 0:
			memmap.SetReg(hw_display.BG0Controll, controll)
			memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.BG0)
			memmap.SetReg(hw_display.BG0HOffset, e.activeBackgrounds[0].hScroll)
			memmap.SetReg(hw_display.BG0VOffset, e.activeBackgrounds[0].vScroll)
		case 1:
			memmap.SetReg(hw_display.BG1Controll, controll)
			memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.BG1)
			memmap.SetReg(hw_display.BG1HOffset, e.activeBackgrounds[1].hScroll)
			memmap.SetReg(hw_display.BG1VOffset, e.activeBackgrounds[1].vScroll)
		case 2:
			memmap.SetReg(hw_display.BG2Controll, controll)
			memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.BG2)
			memmap.SetReg(hw_display.BG2HOffset, e.activeBackgrounds[2].hScroll)
			memmap.SetReg(hw_display.BG2VOffset, e.activeBackgrounds[2].vScroll)
		case 3:
			memmap.SetReg(hw_display.BG3Controll, controll)
			memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.BG3)
			memmap.SetReg(hw_display.BG3HOffset, e.activeBackgrounds[3].hScroll)
			memmap.SetReg(hw_display.BG3VOffset, e.activeBackgrounds[3].vScroll)
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

// NewSprite returns a new Sprite
func (e *Engine) NewSprite(tileSet *assets.TileSet) *Sprite {
	return &Sprite{
		engine:  e,
		tileSet: tileSet,
		size:    hw_sprite.Medium,
		shape:   hw_sprite.Square,
	}
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
