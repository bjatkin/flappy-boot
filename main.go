package main

import (
	"embed"
	"fmt"
	"io"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

//go:embed assets/gba
var fs embed.FS

func main() {
	name := "assets/gba/palette_0.p16"
	palIndex := 0

	palette, err := fs.Open(name)
	if err != nil {
		fmt.Printf("failed to open palette file %s: %s", name, err)
		return
	}

	pal, err := io.ReadAll(palette)
	if err != nil {
		fmt.Printf("failed to read palette file: %s", err)
		return
	}

	if palIndex > 0x0010 {
		fmt.Printf("palette bank %d does not exist must be 0-16", palIndex)
		return
	}

	start := 0x0100 + 0x0010*palIndex
	memmap.Copy16(memmap.Palette[start:start+0x0010], pal)

	// palette, err := assets.Open(name)
	// if err != nil {
	// 	fmt.Printf("failed to open palette file %s: %s", name, err)
	// 	return
	// }

	// pal, err := io.ReadAll(palette)
	// if err != nil {
	// 	fmt.Printf("failed to read palette file: %s", err)
	// 	return
	// }

	sprite.LoadPalette16(fs, name, 0)
	// Palette0 := memmap.Palette[0x0100:0x0110]
	// memmap.Copy16(Palette0, pal)
	// memmap.Copy16(hw_sprite.Palette0, pal)

	// memmap.Copy16(memmap.Palette, a)

	// bFile, err := assets.Open("assets/gba/palette_0.p16")
	// if err != nil {
	// 	fmt.Printf("err 1: %s\n", err)
	// 	return
	// }

	// b, err := io.ReadAll(bFile)
	// if err != nil {
	// 	fmt.Printf("err 2: %s\n", err)
	// }

	// memmap.Copy16(memmap.Palette[0x0100:], b)

	for {
	}

	//game.Run(gameplay.NewDemo(assetFS))
}
