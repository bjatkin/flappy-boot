package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/game"
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

var pal game.Palette = unsafe.Slice(
	(*memmap.PaletteValue)(unsafe.Pointer(&background[16])),
	16, // 16 is hard coded because a gb4 always has a 16 color palette
)

var BackgroundTileSet = &game.TileSet{
	Count: *(*uint32)(unsafe.Pointer(&background[12])),
	Tiles: unsafe.Slice(
		(*memmap.VRAMValue)(unsafe.Pointer(&background[48])),
		*(*uint32)(unsafe.Pointer(&background[12]))*uint32((background[0]/4)*background[1]),
	),
	Palette: &pal,
}

// TODO: this should be a pointer
var BackgroundTileMap game.TileMap = unsafe.Slice(
	(*memmap.VRAMValue)(unsafe.Pointer(&background[48+*(*uint32)(unsafe.Pointer(&background[12]))*uint32((background[0]/4)*background[1])*2])),
	(*(*uint32)(unsafe.Pointer(&background[4]))/8)*(*(*uint32)(unsafe.Pointer(&background[8]))/8),
)
