package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed grass_sky_bg.gb4
var grassSkyBG []byte

type Palette []memmap.PaletteValue

type Asset struct {
	Width   uint32
	Height  uint32
	Tiles   []memmap.VRAMValue
	TileMap []memmap.VRAMValue
	Palette Palette
}

func NewBG() *Asset {
	// 4 becuase this is a gb4 file
	u16PerTile := uint32((grassSkyBG[0] / 4) * grassSkyBG[1])
	width := *(*uint32)(unsafe.Pointer(&grassSkyBG[4]))
	height := *(*uint32)(unsafe.Pointer(&grassSkyBG[8]))
	tileCount := *(*uint32)(unsafe.Pointer(&grassSkyBG[12]))

	return &Asset{
		Width:  width,
		Height: height,
		Palette: unsafe.Slice(
			(*memmap.PaletteValue)(unsafe.Pointer(&grassSkyBG[16])),
			16, // 16 is hard coded becuase a gb4 always has a 16 color palette
		),
		Tiles: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&grassSkyBG[48])),
			tileCount*u16PerTile,
		),
		TileMap: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&grassSkyBG[48+tileCount*u16PerTile*2])),
			(width/8)*(height/8), // divide by 8 since tilemaps must use 8x8 pixel tiles
		),
	}
}

func (a *Asset) LoadMap(bgWidth int) {
	for i := range a.Palette {
		memmap.Palette[i] = a.Palette[i]
	}

	for i := range a.Tiles {
		memmap.VRAM[i] = a.Tiles[i]
	}

	// 8 because the tiles are 8x8
	// tileWidth := a.Width / 8
	// for i := range a.TileMap {
	// 	// 32 because the background size is 32x32
	// 	memmap.VRAM[memmap.ScreenBlockOffset+((i/int(tileWidth))*bgWidth)+(i%8)] = a.TileMap[i]
	// }
}
