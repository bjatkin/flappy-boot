// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed playerAnim.ts4
var playerAnimTileSet []byte

// PlayerAnimTileSet is the sprite sheet for the player character and all it's associated animations
var PlayerAnimTileSet = &TileSet{
    shape: sprite.Wide,
    size:  sprite.Large,
    count: 48,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&playerAnimTileSet[0])),
        768,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&playerAnimTileSet[1536])),
            16,
        ),
    },

}
