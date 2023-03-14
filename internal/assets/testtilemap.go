package assets

import (
	_ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed test.tm4
var testMap []byte

var TestMap = &TileMap{
	size:    display.BGSizeWide,
	tiles:   unsafe.Slice(
		(*memmap.VRAMValue)(unsafe.Pointer(&testMap)),
		2048,	
	),
	tileSet: TestTileSet,
}
