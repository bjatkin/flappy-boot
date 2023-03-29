// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed resart.ts4
var resartTileSet []byte

// ResartTileSet is restart and quit options for the game over screen
var ResartTileSet = &TileSet{
    shape: sprite.Wide,
    size:  sprite.Medium,
    count: 12,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&resartTileSet[0])),
        192,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&resartTileSet[384])),
            16,
        ),
    },

}
