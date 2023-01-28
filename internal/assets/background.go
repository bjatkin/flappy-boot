package assets

import (
	_ "embed"
	"unsafe"
)

//go:embed grass_sky_bg.gb4
var grassSkyBG []byte

// TODO: should these be VRAM values? Or a more specific type?
type Palette []uint16

type Asset struct {
	Width   uint32
	Height  uint32
	Tiles   []uint16
	TileMap []uint16
	Palette Palette
}

func NewBG() *Asset {
	bitsPerTile := uint32((grassSkyBG[0] / 4) * grassSkyBG[1])
	width := *(*uint32)(unsafe.Pointer(&grassSkyBG[4]))
	height := *(*uint32)(unsafe.Pointer(&grassSkyBG[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&grassSkyBG[12]))

	return &Asset{
		Width:  width,
		Height: height,
		Palette: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[16])),
			16, // 16 is hard coded becuase a gb4 always has a 16 color palette
		),
		Tiles: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[48])),
			tileCount*bitsPerTile,
		),
		TileMap: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[48+tileCount*bitsPerTile])),
			(width/8)*(height/8), // divide by 8 since tilemaps must use 8x8 pixel tiles
		),
	}
}
