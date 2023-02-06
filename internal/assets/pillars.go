package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed pillars.gb4
var pillarData []byte

func NewPillars() *SpriteSheet {
	// 4 because this is a gb4 file
	u16PerTile := uint32((pillarData[0] / 4) * pillarData[1])
	tileCount := *(*uint32)(unsafe.Pointer(&pillarData[12]))

	return &SpriteSheet{
		Count: tileCount,
		Palette: unsafe.Slice(
			(*memmap.PaletteValue)(unsafe.Pointer(&pillarData[16])),
			16,
		),
		Sprites: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&pillarData[48])),
			tileCount*u16PerTile,
		),
	}

}
