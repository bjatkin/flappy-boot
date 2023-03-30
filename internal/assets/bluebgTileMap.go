// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/display"
    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed bluebg.tm4
var bluebgTileMap []byte

// BluebgTileMap is a small selection block for when you game over
var BluebgTileMap = &TileMap{
    Size:    display.BGSizeSmall,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&bluebgTileMap[0])),
        1024,	
    ),
    tileSet: &TileSet{
        shape: sprite.Square,
        size:  sprite.Small,
        count: 20,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&bluebgTileMap[2048])),
            320,
        ),
        palette: &Palette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&bluebgTileMap[2688])),
                16,
            ),
        },
    },
}
