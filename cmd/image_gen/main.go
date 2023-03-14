package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"path/filepath"

	"os"
	"strings"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gb4"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
)

var (
	// p256 is a command flag. If p256 is true generated images will 256 color palette
	p256 bool

	// transparent is a command flag. This specifies the transparent color to use in the generated image
	transparent *gbacol.RGB15
)

func main() {
	err := gb4.GeneratePalette(gb4.Palette{
		Name: "test",
		File: "internal/assets/background.png",
	})
	if err != nil {
		panic(err)
	}

	imgFile, err := os.Open("internal/assets/background.png")
	if err != nil {
		panic(err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	pal, err := gb4.NewPal16(img, nil)
	if err != nil {
		panic(err)
	}

	err = gb4.GenerateTileSet(gb4.TileSet{
		Name:    "test",
		File:    "internal/assets/background.png",
		Size:    "8x8",
		Palette: "test",
	},
		pal,
	)
	if err != nil {
		panic(err)
	}

	err = gb4.GenerateTileMap(gb4.TileMap{
		Name: "test",
		File: "internal/assets/background.png",
	}, gb4.TileSet{
		Name:    "test",
		File:    "internal/assets/background.png",
		Size:    "8x8",
		Palette: "test",
	}, gb4.Palette{
		Name: "test",
		File: "internal/assets/background.png",
	})
	if err != nil {
		panic(err)
	}

	err = gb4.GenerateAssets()
	if err != nil {
		panic(err)
	}

	return
	// look for flags first
	args, err := ExtractFlags(os.Args[1:])
	if err != nil {
		fmt.Println(err)
		return
	}

	argsLen := len(args)
	if argsLen%2 != 0 {
		fmt.Printf("number of arguments must be even, but got %d arguments\n", argsLen)
		return
	}

	var toConvert []*Convert
	for i := 0; i < argsLen; i += 2 {
		path := args[i]
		imgFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("failed to open image file %s %s\n", path, err)
			return
		}

		img, _, err := image.Decode(imgFile)
		if err != nil {
			fmt.Printf("falied to decode image file %s %s\n", path, err)
			return
		}

		size, err := tile.NewSize(args[i+1])
		if err != nil {
			fmt.Printf("%s is not a valid GBA tile size %s\n", args[i+1], err)
			return
		}

		fileName := filepath.Base(path)
		ext := filepath.Ext(fileName)
		toConvert = append(toConvert, &Convert{
			Img:      img,
			TileSize: size,
			Path:     strings.TrimSuffix(path, fileName),
			File:     strings.TrimSuffix(fileName, ext),
		})
	}

	for _, c := range toConvert {
		file, err := os.Create(filepath.Join(c.Path, c.File+".gb4"))
		if err != nil {
			fmt.Printf("filed to create gb4 image file %s", err)
			return
		}

		fmt.Println("to convert: ", c, c.Img.Bounds(), c.Img.At(0, 255))
		err = gb4.Encode(file, c.Img, &gb4.Options{TileSize: &c.TileSize, Transparent: transparent, IncludeMap: c.TileSize == tile.S8x8})
		if err != nil {
			fmt.Printf("%s failed to encode .gb4 file", err)
			return
		}
	}
}

// ExtractFlags gets flags from the provided args list and returns a copy of that list without those flags
// If any flag falis to parse an error will be returned
func ExtractFlags(args []string) ([]string, error) {
	var newArgs []string
	for _, arg := range args {
		switch {
		case arg == "-p256":
			p256 = true
		case strings.HasPrefix(arg, "-transparent="):
			col, err := ParseHexColor(strings.TrimPrefix(arg, "-transparent="))
			if err != nil {
				return nil, fmt.Errorf("invalid hex color format for flag %s", arg)
			}
			transparent = &col
		case strings.HasPrefix(arg, "-t="):
			col, err := ParseHexColor(strings.TrimPrefix(arg, "-t="))
			if err != nil {
				return nil, fmt.Errorf("invalid hex color format for flag %s", arg)
			}
			transparent = &col
		default:
			newArgs = append(newArgs, arg)
		}
	}
	return newArgs, nil
}

// Convert is an image that needs to be converted
type Convert struct {
	Img      image.Image
	TileSize tile.Size
	Path     string
	File     string
}

// ParseHexColor parses a hex color string an returns a valid RGB15 color
// the hex string can be a 3 digit hex color #F0F, or a 6 digit hex color #FFF00FF
func ParseHexColor(s string) (gbacol.RGB15, error) {
	c := color.RGBA{A: 0xFF}
	switch len(s) {
	case 7:
		_, err := fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		if err != nil {
			// TODO
			return 0, fmt.Errorf("invaid color %s", err)
		}

	case 4:
		_, err := fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		if err != nil {
			// TODO
			return 0, fmt.Errorf("invalid color %s", err)
		}
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		// TODO
		return 0, fmt.Errorf("invalid color len %d", len(s))
	}

	return gbacol.NewRGB15(c), nil
}

// Config is a config file for the command
type Config struct {
	Palettes       []gb4.Palette
	TileSets       []gb4.TileSet
	TileMaps       []gb4.TileMap
	SetTransparent string
}
