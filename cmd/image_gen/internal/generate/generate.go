package generate

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/config"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gb4"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
)

// assetDir is the base dir within which all asset files are generated
var assetDir = filepath.Join("internal", "assets")

// PaletteData metadata for a specific palette
type PaletteData struct {
	Palette     color.Palette
	PaletteFile string
	GoFile      string
	Name        string
	Description string
}

func Palette(palette config.Palette) (*PaletteData, error) {
	imgFile, err := os.Open(palette.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %w", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s | %w", palette.File, err)
	}

	pal, err := gb4.NewPal16(img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new 16 color palette %w", err)
	}

	return &PaletteData{
		Palette:     pal,
		PaletteFile: filepath.Join(assetDir, palette.Name+"pal4"),
		GoFile:      filepath.Join(assetDir, palette.Name+"palette.go"),
		Name:        palette.Name,
		Description: palette.Description,
	}, nil
}

type TileSetData struct {
	Tiles      []*tile.Meta
	Name       string
	TileLength int
	Length     int
	Palette    *PaletteData
}

func TileSet(tileSet config.TileSet, palettes map[string]*PaletteData) (*TileSetData, error) {
	imgFile, err := os.Open(tileSet.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %w", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s | %w", tileSet.File, err)
	}

	size, err := tile.NewSize(tileSet.Size)
	if err != nil {
		return nil, fmt.Errorf("invalid tile size %s", tileSet.Size)
	}

	var pal *PaletteData
	if tileSet.Palette == "" {
		pal, err = Palette(config.Palette{
			File: tileSet.File,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create valid palette from image %s | %w", tileSet.File, err)
		}
	} else {
		var ok bool
		pal, ok = palettes[tileSet.Palette]
		if !ok {
			return nil, fmt.Errorf("palette %s does not exists", tileSet.Palette)
		}
	}

	tiles := tile.NewMetaSlice(img, pal.Palette, size)
	uniqueTiles := tile.Unique(tiles)

	return &TileSetData{
		Tiles:      uniqueTiles,
		Name:       tileSet.Name,
		TileLength: len(uniqueTiles),
		Length:     len(uniqueTiles) * 16,
		Palette:    pal,
	}, nil
}

type TileMapData struct {
	Name      string
	Width     int
	Height    int
	TileCount int
	TileSet   *TileSetData
}

func TileMap(tileMap config.TileMap, tileSets map[string]*TileSetData) (*TileMapData, error) {
	imgFile, err := os.Open(tileMap.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %s | %w", tileMap.File, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s | %w", tileMap.File, err)
	}

	if tileMap.TileSet == "" {
		tileSet()
	} else {

	}

}
