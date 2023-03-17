// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/display"
    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed pillars.tm4
var pillarsTileMap []byte

var PillarsTileMap = &TileMap{
   size:    display.BGSizeLarge,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&pillarsTileMap)),
        2048,	
    ),
    tileSet: &TileSet{
        count: 29,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&pillarsTileMap[4096])),
            464,
        ),
        palette: CPalette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&pillarsTileMap[5024])),
                16,
            ),
        },
    },
}
