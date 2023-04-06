package game

import (
	"github.com/bjatkin/flappy_boot/internal/alloc"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/display"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/math"
)

const (
	// White is the color white as a memmap.PaletteValue
	White = memmap.PaletteValue(0x7FFF)

	// Black is the color black as a memmap.PaletteValue
	Black = memmap.PaletteValue(0x0000)
)

// Runable is an interface for a type that can be run by the engine
type Runable interface {
	Init(*Engine) error
	Update(*Engine) error
}

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// activeBackgrounds are the backgrounds that need to be drawn each frame
	activeBackgrounds [4]*Background

	// palBuff is the palette buffer that is coppied into the palette memory every frame
	palBuff  [512]memmap.PaletteValue
	fadeCol  memmap.PaletteValue
	fadeFrac math.Fix8
	doFade   bool

	// Allocators
	bgPalAlloc   *alloc.Pal
	sprPalAlloc  *alloc.Pal
	bgTileAlloc  *alloc.VRAM
	sprTileAlloc *alloc.VRAM
	mapAlloc     *alloc.VRAM

	// previousKey holds the state of the hardware key input register durring the last KeyPoll
	// it is used to check key transition states
	previousKeys memmap.Input

	// currentKey holds the state of the hardware key input register the current KeyPoll
	// it is used to check key transition states
	currentKeys memmap.Input

	// frame is the current engine frame
	frame int

	// Debug contains some simple sprites for debugging
	Debug [10]*Sprite
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.Palette[0] = White
	memmap.SetReg(hw_display.Controll, hw_display.Sprite1D|hw_display.ForceBlank)

	e := &Engine{
		activeSprites: make(map[*Sprite]struct{}, 128),

		// the first tile is left transparent and can be shared by all tile maps
		bgTileAlloc:  alloc.NewVRAM(memmap.VRAM[memmap.TileOffset4:memmap.CharBlockOffset*2], 16),
		sprTileAlloc: alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*4:], 16),
		mapAlloc:     alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*2:], memmap.HalfKByte*2),
	}

	e.bgPalAlloc = alloc.NewPal(e.palBuff[:256])
	e.sprPalAlloc = alloc.NewPal(e.palBuff[256:])

	debugSprites := [10]*Sprite{}
	for i := range debugSprites {
		debugSprites[i] = e.NewSprite(assets.DebugTileSet)
	}

	e.Debug = debugSprites

	return e
}

// Frame return the current engine frame
func (e *Engine) Frame() int {
	return e.frame
}

// Run runs the provided Runable
func (e *Engine) Run(run Runable) error {
	// enable sprites
	memmap.SetReg(hw_display.Controll, *hw_display.Controll|hw_display.Sprites)
	// hide all the sprites before the engine starts
	e.drawSprites()

	for {
		err := run.Init(e)
		if err != nil {
			exit(err)
		}

		for {
			e.keyPoll()
			err := run.Update(e)
			if err != nil {
				exit(err)
			}

			e.frame++

			vSyncWait()

			// update the palette if needed
			if e.doFade || e.bgPalAlloc.IsDirty() || e.sprPalAlloc.IsDirty() {
				e.updatePalette()
				e.doFade = false
				e.bgPalAlloc.MarkClean()
				e.sprPalAlloc.MarkClean()
			}

			// copy active sprite data into OAM memory
			e.drawSprites()

			// copy active background data into the background registers
			e.drawBackgrounds()
		}
	}
}

// PalFade fades the current color palette towards the specified color
// t is clamped to between 0 and 1. At 0 the color palette is completely unchanged.
// At 1 the palette is completely the provided color.
func (e *Engine) PalFade(color memmap.PaletteValue, t math.Fix8) {
	t = math.Clamp(t, 0, math.FixOne)
	if e.fadeCol != color || e.fadeFrac != t {
		e.doFade = true
	}
	e.fadeCol = color
	e.fadeFrac = t
}

// updatePalette will copy the current palette into palette RAM
func (e *Engine) updatePalette() {
	for i := range e.palBuff {
		memmap.Palette[i] = lerpColor(e.palBuff[i], e.fadeCol, e.fadeFrac)
	}
}

// lerpColor lerps from the src color to the dest color. t should be between 0 and 1
// at t=0, src is the returned color. at t=1, dest is the returned color.
func lerpColor(src, dest memmap.PaletteValue, t math.Fix8) memmap.PaletteValue {
	switch {
	case t < 15:
		return src
	case t > math.FixOne-15:
		return dest
	case src == dest:
		return src
	}

	// TODO this could be made a lot more efficient if it didin't use math.Lerp
	redMask := memmap.PaletteValue(0b0_00000_00000_11111)
	redSrc := math.NewFix8(int(src&redMask), 0)
	redDest := math.NewFix8(int(dest&redMask), 0)
	red := math.Lerp(redSrc, redDest, t).Int()

	greenMask := memmap.PaletteValue(0b0_00000_11111_00000)
	greenSrc := math.NewFix8(int(src&greenMask)>>5, 0)
	greenDest := math.NewFix8(int(dest&greenMask)>>5, 0)
	green := math.Lerp(greenSrc, greenDest, t).Int()

	blueSrc := math.NewFix8(int(src>>10), 0)
	blueDest := math.NewFix8(int(dest>>10), 0)
	blue := math.Lerp(blueSrc, blueDest, t).Int()

	return memmap.PaletteValue(red | green<<5 | blue<<10)
}

// drawSprites copies all the engines active sprites into OAM memory
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

// drawBackgrounds updates all the background registers and the display controll register based on the
// engines active background
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
