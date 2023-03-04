package game

import (
	"fmt"

	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// activeBackgrounds are the backgrounds that need to be drawn each frame
	activeBackgrounds [4]bool

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
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.SetReg(display.Controll, display.Sprite1D|display.ForceBlank)
	return &Engine{
		activeSprites:  make(map[*Sprite]struct{}, 128),
		spritePtr:      2048,
		spritePalPtr:   16,
		screenBlockPtr: 24,
	}
}

// loadSpriteTileSet loads a tileset into object VRAM, and it's palette into object palette memory
// if there is not enough memory for either the palette or the graphics, an error is returned
func (e *Engine) loadSpriteTileSet(tileSet *TileSet) error {
	err := tileSet.Load(e.spritePalPtr, e.spritePtr)
	if err != nil {
		return err
	}

	e.spritePalPtr++
	e.spritePtr += int(tileSet.Count)

	return nil
}

// loadBGTileSet loads a tileset into background VRAM, and it's palette into background palette memory
// if there is not enough memory for either the palette or the graphics, an error is returned
func (e *Engine) loadBGTileSet(tileSet *TileSet) error {
	err := tileSet.Load(e.bgPalPtr, e.bgPtr)
	if err != nil {
		return err
	}

	e.bgPalPtr++
	e.bgPtr += len(tileSet.Tiles)

	return nil
}

// loadTileMap loads a tilemap into background VRAM. if there is not enought memory for either
// the palette or the graphics, an error is returned
func (e *Engine) loadTileMap(tileMap *TileMap) error {
	if e.screenBlockPtr > 32 {
		return fmt.Errorf("OOM: not enough screen blocks for a new map")
	}

	tileMap.Load(e.screenBlockPtr)
	switch tileMap.ScreenSize {
	case display.BGSizeSmall:
		e.screenBlockPtr += 1
	case display.BGSizeWide:
		e.screenBlockPtr += 2
	case display.BGSizeTall:
		e.screenBlockPtr += 2
	case display.BGSizeLarge:
		e.screenBlockPtr += 4
	}

	return nil
}

// addBG adds a new background to the list of active backgrounds
func (e *Engine) addBG(background *Background) error {
	use := -1
	for i := range e.activeBackgrounds {
		if !e.activeBackgrounds[i] {
			use = i
			break
		}
	}

	sbb := memmap.BGControll(background.tileMap.screenBaseBlock) << display.SBBShift
	cbb := memmap.BGControll(background.tileSet.charBaseBlock) << display.CBBShift
	controll := background.controll | sbb | cbb | background.tileMap.ScreenSize

	switch use {
	case 0:
		memmap.SetReg(display.BG0Controll, controll)
		memmap.SetReg(display.Controll, display.BG0)
	case 1:
		memmap.SetReg(display.BG1Controll, controll)
		memmap.SetReg(display.Controll, display.BG1)
	case 2:
		memmap.SetReg(display.BG2Controll, controll)
		memmap.SetReg(display.Controll, display.BG2)
	case 3:
		memmap.SetReg(display.BG3Controll, controll)
		memmap.SetReg(display.Controll, display.BG3)
	default:
		return fmt.Errorf("OOM: no unused backgrounds available")
	}

	e.activeBackgrounds[use] = true

	return nil
}

// NewBackground returns a new Background
func (e *Engine) NewBackground(tileSet *TileSet, tilemap *TileMap, priority memmap.BGControll) *Background {
	return &Background{
		engine:  e,
		tileMap: tilemap,
		tileSet: tileSet,
	}
}

// NewSprite returns a new Sprite
func (e *Engine) NewSprite() *Sprite {
	return &Sprite{
		engine: e,
	}
}

// Background represents a normal background layer
type Background struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	loaded bool

	tileSet *TileSet

	tileMap *TileMap

	controll memmap.BGControll
	hScroll  uint16
	vScroll  uint16
}

// Load loads a backgrounds data into memory
// if there is not enough free VRAM to accommodate this background an error will be returned
func (b *Background) Load() error {
	if !b.tileSet.loaded {
		err := b.engine.loadBGTileSet(b.tileSet)
		if err != nil {
			return err
		}
	}

	if !b.loaded {
		err := b.engine.loadTileMap(b.tileMap)
		if err != nil {
			return err
		}
	}

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

// Sprite is a game engine sprite
type Sprite struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	tileSet *TileSet
}

// Add adds the sprite to the list of active sprites.
// if the sprites associated assets have not been loaded yet, Add will automatically attempt to load them.
// all active sprites are drawn every frame, if more than 128 sprites are active at a time all active
// sprites will be randomly flickered to ensure all sprites continue to be drawn
func (s *Sprite) Add() {
	s.Load()
	s.engine.activeSprites[s] = struct{}{}
}

// Remove removes the sprites from the list of active sprites.
// removing a sprites does not unload it's loaded assets from VRAM. To do that you must call Unload
func (s *Sprite) Remove() {
	delete(s.engine.activeSprites, s)
}

// Load loads a sprites graphics data into memory
// if there is not enough free VRAM to accomodate the sprite an error will be returned
func (s *Sprite) Load() error {
	if !s.tileSet.loaded {
		return nil
	}

	return s.engine.loadSpriteTileSet(s.tileSet)
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

	// palette is the palette data for the tileset
	Palette      *Palette
	PaletteIndex int
}

func (ts *TileSet) Load(palBase, tileBase int) error {
	// TODO: get rid of these magic numbers
	if palBase > 32 {
		return fmt.Errorf("OOM: invalid palette base %d", palBase)
	}

	// TODO: get rid of these magic numbers
	if tileBase > 512*6 {
		return fmt.Errorf("OOM: invalid tile base %d", tileBase)
	}

	ts.PaletteIndex = palBase
	for i := range *ts.Palette {
		memmap.Palette[i+memmap.PaletteOffset*palBase] = (*ts.Palette)[i]
	}

	ts.TileIndex = tileBase
	for i := range ts.Tiles {
		memmap.VRAM[i+tileBase] = ts.Tiles[i]
	}

	ts.charBaseBlock = tileBase / 512
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
}

func (t *TileMap) Load(screenBlock int) {
	for i := range t.Data {
		memmap.VRAM[i+memmap.ScreenBlockOffset*screenBlock] = t.Data[i]
	}
	t.screenBaseBlock = screenBlock

	t.loaded = true
}

// Paletet is a 16 color palette
type Palette []memmap.PaletteValue
