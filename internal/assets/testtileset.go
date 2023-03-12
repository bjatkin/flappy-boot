package assets

import (
    _ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed test.ts4
var testPixels []byte

// TestTileSet is 
var TestTileSet = &TileSet{
    count: 23,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&testPixels)),
        368,
    ),
    palette: testPalette,
}