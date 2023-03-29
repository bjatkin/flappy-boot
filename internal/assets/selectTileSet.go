// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed select.ts4
var selectTileSet []byte

// SelectTileSet is a simple spinning select arrow
var SelectTileSet = &TileSet{
    shape: sprite.Square,
    size:  sprite.Small,
    count: 3,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&selectTileSet[0])),
        48,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&selectTileSet[96])),
            16,
        ),
    },

}
