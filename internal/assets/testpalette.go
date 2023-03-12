package assets

import (
    _ "embed"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed test.pal4
var testPalette []byte

// TestPalette is a 16 color palette
var TestPalette Palette = unsafe.Slice(
    (*memmap.PaletteValue)(unsafe.Pointer(&testPalette)),
    16,
)