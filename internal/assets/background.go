package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed background.gb4
var background []byte

func NewBackground() *Asset {
	// 4 because this is a gb4 file
	u16PerTile := uint32((background[0] / 4) * background[1])
	width := *(*uint32)(unsafe.Pointer(&background[4]))
	height := *(*uint32)(unsafe.Pointer(&background[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&background[12]))

	return &Asset{
		Width:  width,
		Height: height,
		Palette: unsafe.Slice(
			(*memmap.PaletteValue)(unsafe.Pointer(&background[16])),
			16, // 16 is hard coded because a gb4 always has a 16 color palette
		),

		// TODO: looks like this is loading in too much data?
		// why does it look like it's loading part of the map data?
		Tiles: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&background[48])),
			tileCount*u16PerTile,
		),

		TileMap: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&background[48+tileCount*u16PerTile*2])),
			(width/8)*(height/8), // divide by 8 since tilemaps must use 8x8 pixel tiles
		),
	}
}
