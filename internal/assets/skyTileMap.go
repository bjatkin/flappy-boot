// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/display"
    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed sky.tm4
var skyTileMap []byte

var SkyTileMap = &TileMap{
    Size:    display.BGSizeWide,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&skyTileMap[0])),
        2048,	
    ),
    tileSet: &TileSet{
        count: 10,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&skyTileMap[4096])),
            160,
        ),
        palette: &Palette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&skyTileMap[4416])),
                16,
            ),
        },
    },
}