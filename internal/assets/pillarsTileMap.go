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
    Size:    display.BGSizeWide,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&pillarsTileMap[0])),
        2048,	
    ),
    tileSet: &TileSet{
        count: 30,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&pillarsTileMap[4096])),
            480,
        ),
        palette: &Palette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&pillarsTileMap[5056])),
                16,
            ),
        },
    },
}
