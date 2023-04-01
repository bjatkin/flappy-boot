// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/display"
    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed mainmenu.tm4
var mainmenuTileMap []byte

// MainmenuTileMap is the main set for the main menu
var MainmenuTileMap = &TileMap{
    Size:    display.BGSizeSmall,
    tiles:   unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&mainmenuTileMap[0])),
        1024,	
    ),
    tileSet: &TileSet{
        shape: sprite.Square,
        size:  sprite.Small,
        count: 59,
        pixels: unsafe.Slice(
            (*memmap.VRAMValue)(unsafe.Pointer(&mainmenuTileMap[2048])),
            944,
        ),
        palette: &Palette{
            colors: unsafe.Slice(
                (*memmap.PaletteValue)(unsafe.Pointer(&mainmenuTileMap[3936])),
                16,
            ),
        },
    },
}
