// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed debug.ts4
var debugTileSet []byte

// DebugTileSet is small tileset useful for debugging
var DebugTileSet = &TileSet{
    shape: sprite.Square,
    size:  sprite.Small,
    count: 4,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&debugTileSet[0])),
        64,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&debugTileSet[128])),
            16,
        ),
    },

}
