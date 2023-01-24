package main

import (
	_ "embed"
	"unsafe"
)

//go:embed grass_sky_bg.gb4
var grassSkyBG []byte

type Palette []uint16
type Tile [16]uint16

type Asset struct {
	width   uint32
	height  uint32
	tiles   []uint16
	tileMap []Tile
	palette Palette
}

func NewBG() *Asset {
	width := *(*uint32)(unsafe.Pointer(&grassSkyBG[4]))
	height := *(*uint32)(unsafe.Pointer(&grassSkyBG[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&grassSkyBG[12]))

	return &Asset{
		width:  width,
		height: height,
		palette: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[16])),
			16,
		),
		tiles: unsafe.Slice(
			(*uint16)(unsafe.Pointer(&grassSkyBG[32])),
			tileCount,
		),
		// tileMap: unsafe.Slice(
		// 	(*Tile)(unsafe.Pointer(&grassSkyBG[32+tileCount*2])),
		// 	width*height,
		// ),
	}
}
