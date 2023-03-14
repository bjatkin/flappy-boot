package main

import (
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"os"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/config"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
)

func main() {
	// look for an input yaml file
	if len(os.Args) != 2 {
		fmt.Println("invalie usage, missing config file", os.Args)
		return
	}

	file := os.Args[1]
	cfg, err := config.NewConfigFromFile(file)
	if err != nil {
		fmt.Printf("failed to read in config %s\n", err)
		return
	}

	fmt.Println(cfg)
	// convert build up all the color palettes first
}

// Convert is an image that needs to be converted
type Convert struct {
	Img      image.Image
	TileSize tile.Size
	Path     string
	File     string
}
