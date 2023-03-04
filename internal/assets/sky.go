package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed sky.gb4
var sky []byte

func NewSky() *Asset {
	// 4 because this is a gb4 file
	u16PerTile := uint32((sky[0] / 4) * sky[1])
	width := *(*uint32)(unsafe.Pointer(&sky[4]))
	height := *(*uint32)(unsafe.Pointer(&sky[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&sky[12]))

	return &Asset{
		Width:  width,
		Height: height,
		Palette: unsafe.Slice(
			(*memmap.PaletteValue)(unsafe.Pointer(&sky[16])),
			16, // 16 is hard coded because a gb4 always has a 16 color palette
		),

		Tiles: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&sky[48])),
			tileCount*u16PerTile,
		),

		TileMap: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&sky[48+tileCount*u16PerTile*2])),
			(width/8)*(height/8), // divide by 8 since tilemaps must use 8x8 pixel tiles
		),
	}
}

var SkyPal game.Palette = unsafe.Slice(
	(*memmap.PaletteValue)(unsafe.Pointer(&background[16])),
	16,
)

var SkyTileSet = &game.TileSet{
	Count: *(*uint32)(unsafe.Pointer(&sky[12])),
	Tiles: unsafe.Slice(
		(*memmap.VRAMValue)(unsafe.Pointer(&sky[48+*(*uint32)(unsafe.Pointer(&sky[12]))*uint32((sky[0]/4)*sky[1])*2])),
		(*(*uint32)(unsafe.Pointer(&sky[4]))/8)*(*(*uint32)(unsafe.Pointer(&sky[8]))/8),
	),
	Palette: &SkyPal,
}

var SkyTileMap = &game.TileMap{
	ScreenSize: display.BGSizeWide,
	Data: unsafe.Slice(
		(*memmap.VRAMValue)(unsafe.Pointer(&sky[48+*(*uint32)(unsafe.Pointer(&sky[12]))*uint32((sky[0]/4)*sky[1])*2])),
		(*(*uint32)(unsafe.Pointer(&sky[4]))/8)*(*(*uint32)(unsafe.Pointer(&sky[8]))/8),
	),
}
