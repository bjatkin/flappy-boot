// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed player.ts4
var playerTileSet []byte

// PlayerTileSet is the sprite sheet for the player character
var PlayerTileSet = &TileSet{
    shape: sprite.Square,
    size:  sprite.Medium,
    count: 4,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&playerTileSet[0])),
        64,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&playerTileSet[128])),
            16,
        ),
    },

}
