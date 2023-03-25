// This is generated code. DO NOT EDIT

package assets

import (
	"github.com/bjatkin/flappy_boot/internal/alloc"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// TileMap is tilemap data for a background
type TileMap struct {
	// dirtyTiles are the tiles that have changed since the tilemap was loaded into memory
	dirtyTiles []int

	// Size is the size of the background
	Size memmap.BGControll

	// tiles is the tile index data for the tile map
	tiles []memmap.VRAMValue

	// tileSet is the tile set pixel data for the tile map
	tileSet *TileSet

	// alloc is the allocated memory for the tile map in VRAM
	alloc *alloc.VMem
}

// ScreenBaseBlock returns the screen base block for the tile map
func (t *TileMap) ScreenBaseBlock() memmap.BGControll {
	// the 16 offset here is because the bottom half of background VRAM is reserved for tile maps
	return memmap.BGControll(t.alloc.Offset+16)<<display.SBBShift
}

// SetTile updates the tile map with the new tile at the coordinates x,y
func (t *TileMap) SetTile(x, y, tile int) {
	var screen int
	switch t.Size {
	case display.BGSizeSmall:
		screen = 0
	case display.BGSizeWide:
		screen = (x/32)*1024
	case display.BGSizeTall:
		screen = (y/32)*1024
	case display.BGSizeLarge:
		screen = ((x/32)+(y/16))*1024
	}

	x%=32
	y%=32
	i := screen+y*32+x
	t.tiles[i] = memmap.VRAMValue(tile)

	t.dirtyTiles = append(t.dirtyTiles, i)
}

// Load loads the tile map data into vram memory, if the tile map has already been loaded into memory
// load will only load those tiles which have changed since the last tile load was called. Therefor it is
// safe to call load repeatedly
func (t *TileMap) Load(mapAlloc, tileAlloc *alloc.VRAM, palAlloc *alloc.Pal) error {
	// make sure the associated tile set is loaded
	err := t.tileSet.Load(tileAlloc, palAlloc)
	if err != nil {
		return err
	}

	if t.alloc == nil {
		var screens int
		switch t.Size {
		case display.BGSizeLarge:
			screens = 4
		case display.BGSizeTall:
			screens = 2
		case display.BGSizeWide:
			screens = 2
		case display.BGSizeSmall:
			screens = 1
		}

		t.alloc, err = mapAlloc.Alloc(screens)
		if err != nil {
			return err
		}

		for i := range t.tiles {
			if t.tiles[i] == 0 {
				t.alloc.Memory[i] = 0
				continue
			}
			t.alloc.Memory[i] = (t.tiles[i] + memmap.VRAMValue(t.tileSet.alloc.Offset)) | t.tileSet.TilePalette()
		}
		t.dirtyTiles = []int{}
	}

	for _, i := range t.dirtyTiles {
		if t.tiles[i] == 0 {
			t.alloc.Memory[i] = 0
			continue
		}
		t.alloc.Memory[i] = (t.tiles[i] + memmap.VRAMValue(t.tileSet.alloc.Offset)) | t.tileSet.TilePalette()
	}
	t.dirtyTiles = []int{}

	return nil
}

// Free frees the space that was allocated for this tile map in vram
func (t *TileMap) Free(mapAlloc, tileAlloc *alloc.VRAM, palAlloc *alloc.Pal) {
	mapAlloc.Free(t.alloc)
	t.alloc = nil

	// TODO: use reference counting so tilesets can be shared
	t.tileSet.Free(tileAlloc, palAlloc)
}

// TileSet is tileset data for a background or sprite
type TileSet struct {
	// the number of 8x8 tiles in this tileset
	count int

	// pixels contains the pixel data for the tileset
	pixels []memmap.VRAMValue

	// palette is the palette data for the tileset
	palette *Palette

	// alloc is the allocated memory in VRAM for the tile set
	alloc *alloc.VMem
}

// Offset is the offset in tiles into tile memory where this tile set was loaded
func (t *TileSet) Offset() sprite.Attr2 {
	return sprite.Attr2(t.alloc.Offset)
}

// SprPalette is the palette number that this tileset uses
func (t *TileSet) SprPalette() sprite.Attr2 {
	return sprite.Attr2(t.palette.alloc.Offset)<<sprite.PalShift
}

// TilePalette is the palette number that this tileset uses
func (t *TileSet) TilePalette() memmap.VRAMValue{
	return memmap.VRAMValue(t.palette.alloc.Offset) << display.PaletteShift
}

// Load the tileset into vram
func (t *TileSet) Load(tileAlloc *alloc.VRAM, palAlloc *alloc.Pal) error {
	t.palette.Load(palAlloc)

	if t.alloc == nil {
		var err error
		t.alloc, err = tileAlloc.Alloc(t.count)
		if err != nil {
			return err
		}

		// don't use copy as it may copy data one byte at a time.
		// pixel data must be coppied 16-bits at a time or the pixels will be corrupted
		for i := range t.pixels {
			t.alloc.Memory[i] = t.pixels[i]
		}
	}

	return nil
}

// Free frees the space that was allocated for this tileset in vram
func (t *TileSet) Free(tileAlloc *alloc.VRAM, palAlloc *alloc.Pal) {
	tileAlloc.Free(t.alloc)
	t.alloc = nil

	// TODO: use reference counting so that palettes can be shared
	t.palette.Free(palAlloc)
}

// Palette is a 16 color palette
type Palette struct {
	colors []memmap.PaletteValue
	alloc  *alloc.PMem
}

// Load loads the palette into the gba's palette memory
func (p *Palette) Load(alloc *alloc.Pal) error {
	if p.alloc == nil {
		var err error
		p.alloc, err = alloc.Alloc()
		if err != nil {
			return err
		}

		// don't use copy as it may copy data one byte at a time.
		// color data must be coppied 16-bits at a time or the value will be corrupted
		for i := range p.colors {
			p.alloc.Memory[i] = p.colors[i]
		}
	}

	return nil
}

// Free frees the space that was allocated for this palette in palette memory
func (p *Palette) Free(palAlloc *alloc.Pal) {
	palAlloc.Free(p.alloc)
	p.alloc = nil
}
