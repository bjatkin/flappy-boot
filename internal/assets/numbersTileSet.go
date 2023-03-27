// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed numbers.ts4
var numbersTileSet []byte

// NumbersTileSet is a 16x16 set of digits (0-9)
var NumbersTileSet = &TileSet{
    count: 40,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&numbersTileSet[0])),
        640,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&numbersTileSet[1280])),
            16,
        ),
    },

}
