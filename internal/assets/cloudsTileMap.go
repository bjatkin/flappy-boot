// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/display"
    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

//go:embed clouds.tm4
var cloudsTileMap []byte

// CloudsTileMap is the background clouds
var CloudsTileMap = &TileMap{
    Size:    display.BGSizeWide,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&cloudsTileMap[0])),
        2048,	
    ),
    tileSet: &TileSet{
        count: 3,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&cloudsTileMap[4096])),
            48,
        ),
        palette: &Palette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&cloudsTileMap[4192])),
                16,
            ),
        },
    },
}
