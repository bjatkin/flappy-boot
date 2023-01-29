package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed player.gb4
var playerData []byte

type SpriteSheet struct {
	Count   uint32
	Sprites []memmap.VRAMValue
	Palette Palette
}

func NewPlayer() *SpriteSheet {
	// fmt.Println("len: ", len(playerData))
	// 4 because this is a gb4 file
	u16PerTile := uint32((playerData[0] / 4) * playerData[1])
	tileCount := *(*uint32)(unsafe.Pointer(&playerData[12]))

	return &SpriteSheet{
		Count: tileCount,
		Palette: unsafe.Slice(
			(*memmap.PaletteValue)(unsafe.Pointer(&playerData[16])),
			16, // 16 is hard coded becuase a gb4 always has a 16 color palette
		),
		Sprites: unsafe.Slice(
			(*memmap.VRAMValue)(unsafe.Pointer(&playerData[48])),
			tileCount*u16PerTile,
		),
	}
}

func (s *SpriteSheet) Load() {
	for i := range s.Palette {
		memmap.Palette[i+memmap.PaletteOffset*16] = s.Palette[i]
	}

	for i := range s.Sprites {
		memmap.VRAM[i+memmap.CharBlockOffset*4] = s.Sprites[i]
	}
}
