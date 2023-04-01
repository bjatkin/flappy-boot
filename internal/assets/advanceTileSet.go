// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed advance.ts4
var advanceTileSet []byte

// AdvanceTileSet is the advance badge that goes with the flappy boot logo
var AdvanceTileSet = &TileSet{
    shape: sprite.Square,
    size:  sprite.Medium,
    count: 12,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&advanceTileSet[0])),
        192,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&advanceTileSet[384])),
            16,
        ),
    },

}
