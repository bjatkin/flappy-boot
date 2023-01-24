package assets

import (
	_ "embed"
	"unsafe"
)

//go:embed grass_sky_bg_2.gb4
var grassSkyBG []byte

type Palette []uint16
type Tile [16]uint16

type Asset struct {
	Width   uint32
	Height  uint32
	Tiles   []uint16
	TileMap []Tile
	Palette Palette
}

func NewBG() *Asset {
	width := *(*uint32)(unsafe.Pointer(&grassSkyBG[4]))
	height := *(*uint32)(unsafe.Pointer(&grassSkyBG[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&grassSkyBG[12]))

	return &Asset{
		Width:  width,
		Height: height,
		Palette: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[16])),
			16,
		),
		Tiles: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[32])),
			tileCount,
		),
		TileMap: unsafe.Slice(
			(*Tile)(unsafe.Pointer(&grassSkyBG[32+tileCount*2])),
			width*height,
		),
	}
}
