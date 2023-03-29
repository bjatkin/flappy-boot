// This is generated code. DO NOT EDIT

package assets

import (
    _ "embed"
    "unsafe"

    "github.com/bjatkin/flappy_boot/internal/hardware/memmap"
    "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

//go:embed banners.ts4
var bannersTileSet []byte

// BannersTileSet is score banners for the game over screen
var BannersTileSet = &TileSet{
    shape: sprite.Wide,
    size:  sprite.Large,
    count: 32,
    pixels: unsafe.Slice(
        (*memmap.VRAMValue)(unsafe.Pointer(&bannersTileSet[0])),
        512,
    ),

    palette: &Palette{
        colors: unsafe.Slice(
            (*memmap.PaletteValue)(unsafe.Pointer(&bannersTileSet[1024])),
            16,
        ),
    },

}
