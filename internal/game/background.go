package game

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// Background represents a normal background layer
type Background struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	loaded bool

	tileMap *assets.TileMap

	controllReg memmap.BGControll
	hScroll     uint16
	vScroll     uint16
}

// NewBackground returns a new Background
func (e *Engine) NewBackground(tilemap *assets.TileMap, priority memmap.BGControll) *Background {
	return &Background{
		engine:      e,
		tileMap:     tilemap,
		controllReg: priority,
	}
}

// Load loads a backgrounds data into memory
// if there is not enough free VRAM to accommodate this background an error will be returned
func (b *Background) Load() error {
	if b.loaded {
		return nil
	}

	err := b.tileMap.Load(b.engine.mapAlloc, b.engine.bgTileAlloc, b.engine.bgPalAlloc)
	if err != nil {
		return err
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

	err = b.engine.addBackground(b)
	if err != nil {
		return err
	}

	return nil
}

// Remove removes the background for the list of active backgrounds.
// removing a background does not unload it's loaded assets from VRAM. To do that you must call Unload
func (b *Background) Remove() {
	b.engine.removeBackground(b)
}

func (b *Background) Unload() {
	b.tileMap.Free(b.engine.mapAlloc, b.engine.bgTileAlloc, b.engine.sprPalAlloc)
	b.loaded = false
}

// controll returns the correct value for the background controll registers for the given background
func (b *Background) controll() memmap.BGControll {
	return b.controllReg |
		b.tileMap.ScreenBaseBlock() |
		b.tileMap.Size
}

// TODO: introduce a v2 type?
func (b *Background) Scroll(dx, dy int) {
	b.hScroll += uint16(dx)
	b.vScroll += uint16(dy)
}

func (b *Background) SetTile(x, y, tile int) {
	b.tileMap.SetTile(x, y, tile)
}
