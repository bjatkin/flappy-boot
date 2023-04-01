// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed start.ts4
var startTileSet []byte

// StartTileSet is the PRESS START text present on the title screen
var StartTileSet = &TileSet{
    shape: sprite.Wide,
    size:  sprite.Small,
    count: 12,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&startTileSet[0])),
        192,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&startTileSet[384])),
            16,
        ),
    },

}
