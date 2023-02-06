package game

import (
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// spriteBlocSize is the size of each VRAM sprite block in 8x8 tiles
const spriteBlockSize = 4

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// spriteBlocks indicates which sprite block are free and which are in use
	spriteBlocks []bool

	// spritePalettes indicates which sprite palettes are free and which are in use
	spritePalettes []bool
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	return &Engine{
		activeSprites:  make(map[*Sprite]struct{}, 128),
		spriteBlocks:   make([]bool, 1024/spriteBlockSize),
		spritePalettes: make([]bool, 16),
	}
}

// loadSprite loads a sprites assets into VRAM, and it's palette into palette memory
// if there is not enough memory for either the palette or the graphics, an error is returned
func (e *Engine) loadSprite(sprite *Sprite) error {
	// TODO: load the palette

	// TODO: load the sprite tiles
	return nil
}

func (e *Engine) freeSprite() {

}

// NewBackground returns a new Background
func (e *Engine) NewBackground() *Background {
	return &Background{}
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

	// loaded is true if the sprites associated graphics have been loaded into memory
	loaded bool
}

// Load loads a backgrounds data into memory
// if there is not enough free VRAM to accommodate this background an error will be returned
func (b *Background) Load() error {
	return nil
}

// Add adds the background to the list of active backgrounds.
// if the background has not yet been loaded, Add will automatically attempt to load them.
// all active backgrounds are drawn every frame, if the maximum number of backgrounds are already
// active an error will be returned
func (b *Background) Add() error {
	return nil
}

// Remove removes the background for the list of active backgrounds.
// removing a background does not unload it's loaded assets from VRAM. To do that you must call Unload
func (b *Background) Remove() {}

// Unload removes all the backgrounds assets from VRAM.
// if the background is currently active Unload will first deactivate the background
func (b *Background) Unload() {}

// Sprite is a game engine sprite
type Sprite struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	// loaded is true if the sprites associated graphics have been loaded into memory
	loaded bool
}

// Load loads a sprites graphics data into memory
// if there is not enough free VRAM to accomodate the sprite an error will be returned
func (s *Sprite) Load() error {
	if s.loaded {
		return nil
	}

	s.loaded = true
	return nil
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

// Unload removes all the sprits graphis from VRAM.
// if the sprite is currently active Unload will first ensure deactivate the sprite
func (s *Sprite) Unload() {
	s.Remove()
}

// Paletet is a 16 color palette
type Palette struct {
	id     uint32
	colors [16]memmap.PaletteValue
}
