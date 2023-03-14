package game

import (
	"fmt"

	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

type Runable interface {
	Init(*Engine) error
	Update(*Engine, int) error
	Draw(*Engine) error
	Done() (Runable, bool)
}

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// activeBackgrounds are the backgrounds that need to be drawn each frame
	activeBackgrounds [4]*Background

	// spritePtr points to the next available sprite tile in VRAM
	spritePtr int

	// spritePalPtr points to the next available sprite palette in palette memory
	spritePalPtr int

	// bgPtr points to the next available background tile in VRAM
	bgPtr int

	// bgPalPtr points to the next available background palette in palette memory
	bgPalPtr int

	// screenBlockPtr points to the next available screen block
	screenBlockPtr int

	// Allocators
	bgPals   [8]bool
	sprPals  [8]bool
	bgTiles  Allocator
	bgMaps   Allocator
	sprTiles Allocator
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.SetReg(display.Controll, display.Sprite1D|display.ForceBlank)
	return &Engine{
		activeSprites:  make(map[*Sprite]struct{}, 128),
		spritePtr:      memmap.CharBlockOffset * 4,
		spritePalPtr:   16,
		screenBlockPtr: 24,
	}
}

func (e *Engine) Run(run Runable) error {
	// enable sprites
	memmap.SetReg(display.Controll, *display.Controll|display.Sprites)
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
			var i int
			for s := range e.activeSprites {
				hw_sprite.OAM[i] = *s.attrs()
				i++
			}

			// copy active background data into the background registers
			for i := range e.activeBackgrounds {
				if e.activeBackgrounds[i] == nil {
					continue
				}
				controll := e.activeBackgrounds[i].Controll()

				switch i {
				case 0:
					memmap.SetReg(display.BG0Controll, controll)
					memmap.SetReg(display.Controll, *display.Controll|display.BG0)
					memmap.SetReg(display.BG0HOffset, e.activeBackgrounds[0].hScroll)
					memmap.SetReg(display.BG0VOffset, e.activeBackgrounds[0].vScroll)
				case 1:
					memmap.SetReg(display.BG1Controll, controll)
					memmap.SetReg(display.Controll, *display.Controll|display.BG1)
					memmap.SetReg(display.BG1HOffset, e.activeBackgrounds[1].hScroll)
					memmap.SetReg(display.BG1VOffset, e.activeBackgrounds[1].vScroll)
				case 2:
					memmap.SetReg(display.BG2Controll, controll)
					memmap.SetReg(display.Controll, *display.Controll|display.BG2)
					memmap.SetReg(display.BG2HOffset, e.activeBackgrounds[2].hScroll)
					memmap.SetReg(display.BG2VOffset, e.activeBackgrounds[2].vScroll)
				case 3:
					memmap.SetReg(display.BG3Controll, controll)
					memmap.SetReg(display.Controll, *display.Controll|display.BG3)
					memmap.SetReg(display.BG3HOffset, e.activeBackgrounds[3].hScroll)
					memmap.SetReg(display.BG3VOffset, e.activeBackgrounds[3].vScroll)
				}
			}

			err = run.Draw(e)
			if err != nil {
				exit(err)
			}

			if next, ok := run.Done(); ok {
				sprite.Reset()
				run = next
				break
			}
		}
	}
}

type PalAllocation struct {
	Pal    memmap.VRAMValue
	Memory []memmap.PaletteValue
}

func (e *Engine) AllocSprPal() (PalAllocation, error) {
	for i := range e.sprPals {
		if !e.sprPals[i] {
			return PalAllocation{
				Pal:    memmap.VRAMValue(i) << display.PaletteShift,
				Memory: memmap.Palette[memmap.PaletteOffset*(i+16) : memmap.PaletteOffset*(i+17)],
			}, nil
		}
	}
	return PalAllocation{}, ErrOOM
}

func (e *Engine) FreeSprPal(alloc PalAllocation) {
	e.sprPals[alloc.Pal>>display.PaletteShift] = false
}

func (e *Engine) AllocBGPal() (*PalAllocation, error) {
	for i := range e.sprPals {
		if !e.sprPals[i] {
			return &PalAllocation{
				Pal:    memmap.VRAMValue(i) << display.PaletteShift,
				Memory: memmap.Palette[memmap.PaletteOffset*i : memmap.PaletteOffset*(i+1)],
			}, nil
		}
	}
	return nil, ErrOOM

}

func (e *Engine) FreeBGPall(alloc *PalAllocation) {
	e.bgPals[alloc.Pal>>display.PaletteShift] = false
}

type TileAllocation struct {
	Memory []memmap.VRAMValue
	Index  memmap.VRAMValue
}

func (e *Engine) AllocBGTile(tiles int) (*TileAllocation, error) {
	i, err := e.bgTiles.Alloc(tiles)
	if err != nil {
		return nil, err
	}

	start := i * 16
	end := start + (tiles * 16)
	return &TileAllocation{
		Memory: memmap.VRAM[start:end],
		Index:  memmap.VRAMValue(i),
	}, nil
}

func (e *Engine) FreeBGTile(alloc *TileAllocation) {
	e.bgTiles.Free(int(alloc.Index))
}

type MapAllocation struct {
	Memory         []memmap.VRAMValue
	CharacterBlock memmap.BGControll
	ScreenBlock    memmap.BGControll
}

func (e *Engine) AllocBGMap(screens int) (*MapAllocation, error) {
	i, err := e.bgMaps.Alloc(screens)
	if err != nil {
		return nil, err
	}

	characterBlock := memmap.BGControll(2) << display.CBBShift
	if i > 8 {
		characterBlock = 3 << display.CBBShift
		i -= 8
	}

	start := (memmap.CharBlockOffset * 2)
	end := start + (screens * memmap.ScreenBlockOffset)
	return &MapAllocation{
		Memory:         memmap.VRAM[start:end],
		CharacterBlock: characterBlock,
		ScreenBlock:    memmap.BGControll(i) << display.SBBShift,
	}, nil
}

func (e *Engine) FreeBGMap(alloc *MapAllocation) {
	i := alloc.ScreenBlock >> display.SBBShift
	if alloc.CharacterBlock == 3<<display.CBBShift {
		i += 8
	}

	e.bgMaps.Free(int(i))
}

type SprAllocation struct {
	Memory []memmap.VRAMValue
	Offset hw_sprite.Attr2
}

func (e *Engine) AllocSprs(sprites int) (SprAllocation, error) {
	i, err := e.sprTiles.Alloc(sprites)
	if err != nil {
		return SprAllocation{}, ErrOOM
	}

	start := 4*memmap.CharBlockOffset + i*16
	end := start + sprites*16
	return SprAllocation{
		Memory: memmap.VRAM[start:end],
		Offset: hw_sprite.Attr2(i),
	}, nil
}

func (e *Engine) FreeSprs(alloc SprAllocation) {
	e.sprTiles.Free(int(alloc.Offset))
}

func (e *Engine) loadBGPal(data []memmap.PaletteValue) int {
	palID := e.bgPalPtr
	for i := range data {
		memmap.Palette[i+(memmap.PaletteOffset*palID)] = data[i]
	}

	e.bgPalPtr += 1

	return palID
}

func (e *Engine) loadSprPal(data []memmap.PaletteValue) int {
	palID := e.spritePalPtr
	for i := range data {
		memmap.Palette[i+memmap.PaletteOffset*palID] = data[i]
	}
	e.spritePalPtr++

	return palID - 16
}

func (e *Engine) loadCB(data []memmap.VRAMValue) (int, int) {
	tileOffset := e.bgPtr
	for i := range data {
		memmap.VRAM[i+tileOffset] = data[i]
	}

	e.bgPtr += len(data)
	return 0, tileOffset / 16
}

func (e *Engine) loadSB(data []memmap.VRAMValue, offset, palID int) int {
	screenID := e.screenBlockPtr
	for i := range data {
		memmap.VRAM[i+memmap.ScreenBlockOffset*screenID] = (data[i] + memmap.VRAMValue(offset)) | memmap.VRAMValue(palID)<<0x0C
	}

	e.screenBlockPtr += len(data) / memmap.ScreenBlockOffset
	return screenID
}

func (e *Engine) updateSB(screenBlock, offset int, tile memmap.VRAMValue) {
	memmap.VRAM[offset+memmap.ScreenBlockOffset*screenBlock] = tile
}

func (e *Engine) loadSprite(data []memmap.VRAMValue) int {
	tileOffset := e.spritePtr
	for i := range data {
		memmap.VRAM[i+tileOffset] = data[i]
	}

	e.spritePtr += len(data)
	return tileOffset - memmap.CharBlockOffset*4
}

// addBG adds a new background to the list of active backgrounds
func (e *Engine) addBG(background *Background) error {
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			e.activeBackgrounds[i] = background
			return nil
		}
	}

	return fmt.Errorf("OOM: no unused backgrounds available")
}

func (e *Engine) addSprite(sprite *Sprite) {
	e.activeSprites[sprite] = struct{}{}
}

// NewBackground returns a new Background
func (e *Engine) NewBackground(tilemap *TileMap, priority memmap.BGControll) *Background {
	return &Background{
		engine:   e,
		tileMap:  tilemap,
		controll: priority,
	}
}

// NewSprite returns a new Sprite
func (e *Engine) NewSprite(tileSet *TileSet) *Sprite {
	return &Sprite{
		engine:  e,
		tileSet: tileSet,
		size:    hw_sprite.Medium,
		shape:   hw_sprite.Square,
	}
}

// TODO move to it's own file
// Background represents a normal background layer
type Background struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	loaded bool

	tileMap *TileMap

	controll memmap.BGControll
	hScroll  uint16
	vScroll  uint16
}

// Load loads a backgrounds data into memory
// if there is not enough free VRAM to accommodate this background an error will be returned
func (b *Background) Load() error {
	if b.loaded {
		return nil
	}

	b.tileMap.Load(b.engine)
	b.loaded = true

	return nil
}

// Add adds the background to the list of active backgrounds.
// if the background has not yet been loaded, Add will automatically attempt to load them.
// all active backgrounds are drawn every frame, if the maximum number of backgrounds are already
// active an error will be returned
func (b *Background) Add() error {
	err := b.Load()
	if err != nil {
		return err
	}

	err = b.engine.addBG(b)
	if err != nil {
		return err
	}

	return nil
}

// Remove removes the background for the list of active backgrounds.
// removing a background does not unload it's loaded assets from VRAM. To do that you must call Unload
func (b *Background) Remove() {}

// TODO make this private
// Controll returns the correct value for the background controll registers for the given background
func (b *Background) Controll() memmap.BGControll {
	return b.controll |
		memmap.BGControll(b.tileMap.screenBaseBlock)<<display.SBBShift |
		memmap.BGControll(b.tileMap.TileSet.charBaseBlock)<<display.CBBShift |
		b.tileMap.ScreenSize
}

// TODO: introduce a v2 type?
func (b *Background) Scroll(dx, dy int) {
	b.hScroll += uint16(dx)
	b.vScroll += uint16(dy)
}

func (b *Background) SetTile(x, y, tile int) {
	b.tileMap.Data[y*64+x] = memmap.VRAMValue(tile)
	if b.loaded {
		// this should live inside the engine loader I think?
		// or maybe background should have a .Tile() method that handsl most of this and then just hands stuff to the engine to load...
		t := (memmap.VRAMValue(tile) + memmap.VRAMValue(b.tileMap.TileSet.TileIndex)) | memmap.VRAMValue(b.tileMap.TileSet.PaletteIndex)<<0x0C

		b.engine.updateSB(b.tileMap.screenBaseBlock, y*32+x, t)
	}
}

// Sprite is a game engine sprite
type Sprite struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	Y         fix.P8
	X         fix.P8
	TileIndex int
	Hide      bool
	HFlip     bool
	VFlip     bool
	Priority  int
	size      hw_sprite.Attr1
	shape     hw_sprite.Attr0

	tileSet *TileSet
}

func (s *Sprite) attrs() *hw_sprite.Attrs {
	var hideAttr hw_sprite.Attr0
	if s.Hide {
		hideAttr = hw_sprite.Hide
	}

	var vFlipAttr hw_sprite.Attr1
	if s.VFlip {
		vFlipAttr = hw_sprite.HMirrior
	}
	var hFlipAttr hw_sprite.Attr1
	if s.HFlip {
		hFlipAttr = hw_sprite.VMirrior
	}

	var priorityAttr hw_sprite.Attr2
	switch s.Priority {
	case 0:
		priorityAttr = hw_sprite.Priority0
	case 1:
		priorityAttr = hw_sprite.Priority1
	case 2:
		priorityAttr = hw_sprite.Priority2
	case 3:
		priorityAttr = hw_sprite.Priority3
	}

	return &hw_sprite.Attrs{
		Attr0: hw_sprite.Attr0(s.Y.Int()) | s.shape | hideAttr,
		Attr1: hw_sprite.Attr1(s.X.Int()) | vFlipAttr | hFlipAttr | s.size,
		Attr2: hw_sprite.Attr2(s.TileIndex+s.tileSet.TileIndex) |
			priorityAttr |
			hw_sprite.Attr2(s.tileSet.PaletteIndex)<<hw_sprite.PalShift,
	}
}

// Add adds the sprite to the list of active sprites.
// if the sprites associated assets have not been loaded yet, Add will automatically attempt to load them.
// all active sprites are drawn every frame, if more than 128 sprites are active at a time all active
// sprites will be randomly flickered to ensure all sprites continue to be drawn
func (s *Sprite) Add() error {
	err := s.Load()
	if err != nil {
		return err
	}

	s.engine.addSprite(s)
	return nil
}

// Remove removes the sprites from the list of active sprites.
// removing a sprites does not unload it's loaded assets from VRAM. To do that you must call Unload
func (s *Sprite) Remove() {
	delete(s.engine.activeSprites, s)
}

// Load loads a sprites graphics data into memory
// if there is not enough free VRAM to accomodate the sprite an error will be returned
func (s *Sprite) Load() error {
	err := s.tileSet.Load(s.engine)
	if err != nil {
		return err
	}

	return nil
}

// TileSet is a set of 8x8 tiles that can be loaded into VRAM for use by a background or sprite
type TileSet struct {
	// loaded is true if this tileset has been loaded into VRAM
	loaded bool

	// the number of 8x8 tiles in this tileset
	Count uint32

	// tiles is the pixel data for the tile set
	Tiles     []memmap.VRAMValue
	TileIndex int

	charBaseBlock int
	Sprite        bool

	// palette is the palette data for the tileset
	Palette      Palette
	PaletteIndex int
}

func (ts *TileSet) Load(e *Engine) error {
	if ts.loaded {
		return nil
	}

	if ts.Sprite {
		ts.PaletteIndex = e.loadSprPal(ts.Palette)

		ts.TileIndex = e.loadSprite(ts.Tiles)
	} else {
		ts.PaletteIndex = e.loadBGPal(ts.Palette)

		ts.charBaseBlock, ts.TileIndex = e.loadCB(ts.Tiles)
	}

	// // TODO: get rid of these magic numbers
	// if palBase > 32 {
	// 	return fmt.Errorf("OOM: invalid palette base %d", palBase)
	// }

	// // TODO: get rid of these magic numbers
	// if tileBase > 512*6 {
	// 	return fmt.Errorf("OOM: invalid tile base %d", tileBase)
	// }
	ts.loaded = true

	return nil
}

// TileMap is tilemap data for a background
type TileMap struct {
	// loaded is true if this tilemap has been loaded into VRAM
	loaded bool

	screenBaseBlock int
	ScreenSize      memmap.BGControll
	Data            []memmap.VRAMValue
	TileSet         *TileSet
}

func (t *TileMap) Load(e *Engine) {
	if t.loaded {
		return
	}

	// make sure the tile set is loaded into memory
	t.TileSet.Load(e)

	t.screenBaseBlock = e.loadSB(t.Data, t.TileSet.TileIndex, t.TileSet.PaletteIndex)
	t.loaded = true
}

// Paletet is a 16 color palette
type Palette []memmap.PaletteValue
