// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed logo.ts4
var logoTileSet []byte

// LogoTileSet is the main logo for flappy boot
var LogoTileSet = &TileSet{
    shape: sprite.Square,
    size:  sprite.Large,
    count: 48,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&logoTileSet[0])),
        768,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&logoTileSet[1536])),
            16,
        ),
    },

}
