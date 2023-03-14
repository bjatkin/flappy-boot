// This is generated code. DO NOT EDIT

package assets

import (
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// TileMap is tilemap data for a background
type TileMap struct {
	// dirtyTiles are the tiles that have changed since the tilemap was loaded into memory
	dirtyTiles []int

	// size is the size of the background
	size memmap.BGControll

	// tiles is the tile index data for the tile map
	tiles []memmap.VRAMValue

	// tileSet is the tile set pixel data for the tile map
	tileSet *TileSet

	// alloc is the allocated memory for the tile map in vram
	alloc *game.MapAllocation
}

// SetTile updates the tile map with the new tile at the coordinates x,y
func (t *TileMap) SetTile(x, y, tile int) {
	var width int
	switch t.size {
	case display.BGSizeSmall:
		width = 32
	case display.BGSizeWide:
		width = 64
	case display.BGSizeTall:
		width = 32
	case display.BGSizeLarge:
		width = 32
	}
	i := y*width + x
	t.tiles[i] = memmap.VRAMValue(tile)

	t.dirtyTiles = append(t.dirtyTiles, i)
}

// Load loads the tile map data into vram memory, if the tile map has already been loaded into memory
// load will only load those tiles which have changed since the last tile load was called. Therefor it is
// safe to call load repeatedly
func (t *TileMap) Load(e *game.Engine) error {
	// make sure the associated tile set is loaded
	err := t.tileSet.Load(e)
	if err != nil {
		return err
	}

	if t.alloc == nil {
		var screens int
		switch t.size {
		case display.BGSizeLarge:
			screens = 4
		case display.BGSizeTall:
			screens = 2
		case display.BGSizeWide:
			screens = 2
		case display.BGSizeSmall:
			screens = 1
		}

		t.alloc, err = e.AllocBGMap(screens)
		if err != nil {
			return err
		}

		for i := range t.alloc.Memory {
			t.alloc.Memory[i] = (t.tiles[i] + t.tileSet.alloc.Index) | t.tileSet.paletteIndex
		}
		t.dirtyTiles = []int{}
	}

	for _, i := range t.dirtyTiles {
		t.alloc.Memory[i] = (t.tiles[i] + t.tileSet.alloc.Index) | t.tileSet.paletteIndex
	}
	t.dirtyTiles = []int{}

	return nil
}

// Free frees the space that was allocated for this tile map in vram
func (t *TileMap) Free(e *game.Engine) {
	e.FreeBGMap(t.alloc)
	// TODO: use reference counting so tilesets can be shared
	t.tileSet.Free(e)
	t.alloc = nil
}

// TileSet is tileset data for a background or sprite
type TileSet struct {
	// the number of 8x8 tiles in this tileset
	count int

	// pixels contains the pixel data for the tileset
	pixels []memmap.VRAMValue

	// palette is the palette data for the tileset
	palette CPalette

	// paletteIndex is the palette palette index of the tilesets in use palette
	paletteIndex memmap.VRAMValue

	// alloc is the allocated memory in vram for the tile set
	alloc *game.TileAllocation
}

// Load the tileset into vram
func (t *TileSet) Load(e *game.Engine) error {
	t.palette.Load(e)

	if t.alloc == nil {
		var err error
		t.alloc, err = e.AllocBGTile(t.count)
		if err != nil {
			return err
		}

		for i := range t.alloc.Memory {
			t.alloc.Memory[i] = t.pixels[i]
		}
	}

	return nil
}

// Free frees the space that was allocated for this tileset in vram
func (t *TileSet) Free(e *game.Engine) {
	e.FreeBGTile(t.alloc)
	// TODO: use reference counting so that palettes can be shared
	t.palette.Free(e)
	t.alloc = nil
}

// CPalette is a 16 color palette
type CPalette struct {
	colors []memmap.PaletteValue
	alloc  *game.PalAllocation
}

// Load loads the palette into the gba's palette memory
func (p *CPalette) Load(e *game.Engine) error {
	if p.alloc != nil {
		var err error
		p.alloc, err = e.AllocBGPal()
		if err != nil {
			return err
		}

		for i := range p.alloc.Memory {
			p.alloc.Memory[i] = p.colors[i]
		}
	}

	return nil
}

// Free frees the space that was allocated for this palette in palette memory
func (p *CPalette) Free(e *game.Engine) {
	e.FreeBGPall(p.alloc)
	p.alloc = nil
}
