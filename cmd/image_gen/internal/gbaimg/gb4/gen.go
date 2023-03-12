package gb4

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
)

//go:embed templates
var templates embed.FS

// Palette is a named palette
type Palette struct {
	Name        string
	File        string
	Description string
}

type paletteData struct {
	Name        string
	Description string
}

func Public(name string) string {
	runes := []rune(name)
	if unicode.IsLower(runes[0]) {
		return strings.ToUpper(string(runes[0])) + string(runes[1:])
	}

	return name
}

func Private(name string) string {
	runes := []rune(name)
	if unicode.IsUpper(runes[0]) {
		return strings.ToLower(string(runes[0])) + string(runes[1:])
	}

	return name
}

var paletteTemplate = template.Must(
	template.New("go_templates").
		Funcs(
			map[string]any{
				"private": Private,
				"public":  Public,
			}).
		ParseFS(templates, "templates/*.tmpl"),
)

func GeneratePalette(palette Palette) error {
	imgFile, err := os.Open(palette.File)
	if err != nil {
		return fmt.Errorf("failed to read image file %s", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s", err)
	}

	p, err := NewPal16(img, nil)
	if err != nil {
		return fmt.Errorf("failed to create new 16 color palette %s", err)
	}

	err = os.WriteFile(filepath.Join("internal", "assets", palette.Name+".pal4"), rawPalette(p), 0o0666)
	if err != nil {
		return fmt.Errorf("failed to write %s", err)
	}

	f, err := os.Create(filepath.Join("internal", "assets", palette.Name+"palette.go"))
	if err != nil {
		return fmt.Errorf("failed to create tilemap.go file %s", err)
	}

	description := palette.Description
	if description == "" {
		description = "a 16 color palette"
	}

	err = paletteTemplate.ExecuteTemplate(f, "palette.go.tmpl", &paletteData{
		Name:        palette.Name,
		Description: description,
	})
	if err != nil {
		return fmt.Errorf("failed to execulte palette template")
	}

	return nil
}

// TileSet is a named tile set
type TileSet struct {
	Name    string
	Palette string
	File    string
	Size    string
}

type tileSetData struct {
	Name        string
	Description string
	TileCount   int
	PixelCount  int
	PaletteName string
}

func GenerateTileSet(tileSet TileSet, palette color.Palette) error {
	imgFile, err := os.Open(tileSet.File)
	if err != nil {
		return fmt.Errorf("failed to read image file %s", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s", err)
	}

	size, err := tile.NewSize(tileSet.Size)
	if err != nil {
		return fmt.Errorf("failed to convert size %s", err)
	}

	tiles := tile.NewMetaSlice(img, palette, size)
	uniqueTiles := tile.Unique(tiles)

	err = os.WriteFile(filepath.Join("internal", "assets", tileSet.Name+".ts4"), rawTiles(uniqueTiles), 0o0666)
	if err != nil {
		return fmt.Errorf("failed to create .ts4 file %s", err)
	}

	f, err := os.Create(filepath.Join("internal", "assets", tileSet.Name+"tileset.go"))
	if err != nil {
		return fmt.Errorf("failed to create tilemap.go file %s", err)
	}

	err = paletteTemplate.ExecuteTemplate(f, "tileset.go.tmpl", &tileSetData{
		Name:        tileSet.Name,
		TileCount:   len(uniqueTiles),
		PixelCount:  len(uniqueTiles) * 16,
		PaletteName: tileSet.Palette + "Palette",
	})
	if err != nil {
		return fmt.Errorf("failed to execulte palette template %s", err)
	}

	return nil
}

// TileMap is a named tile map
type TileMap struct {
	Name    string
	File    string
	TileSet string
	Palette string
}

type tileMapData struct {
	Name        string
	Size        string
	TileCount   int
	TileSetName string
}

func GenerateTileMap(tileMap TileMap, tileSet TileSet, palette Palette) error {
	fmt.Println("here 1")
	// create the palette
	imgFile, err := os.Open(palette.File)
	if err != nil {
		return fmt.Errorf("failed to read image file %s", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s", err)
	}

	fmt.Println("here 2")
	p, err := NewPal16(img, nil)
	if err != nil {
		return fmt.Errorf("failed to create new 16 color palette %s", err)
	}

	// create the tile set
	tileSetFile, err := os.Open(tileSet.File)
	if err != nil {
		return fmt.Errorf("failed to read image file %s", err)
	}

	imgTS, _, err := image.Decode(tileSetFile)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s", err)
	}

	size, err := tile.NewSize(tileSet.Size)
	if err != nil {
		return fmt.Errorf("failed to convert size %s", err)
	}

	fmt.Println("here 3")
	tiles := tile.NewMetaSlice(imgTS, p, size)
	uniqueTiles := tile.Unique(tiles)

	// create the tile map
	tileMapFile, err := os.Open(tileMap.File)
	if err != nil {
		return fmt.Errorf("failed to read image file %s", err)
	}

	imgTM, _, err := image.Decode(tileMapFile)
	if err != nil {
		return fmt.Errorf("failed to decode image file %s", err)
	}

	tiles = tile.NewMetaSlice(imgTM, p, size)
	mapData := rawMapData(tiles, uniqueTiles, imgTM.Bounds().Dx(), imgTM.Bounds().Dy())
	fmt.Println("here 4", filepath.Join("internal", "asset", tileMap.Name+".tm4"), len(mapData))
	err = os.WriteFile(filepath.Join("internal", "assets", tileMap.Name+".tm4"), mapData, 0o0666)
	if err != nil {
		return fmt.Errorf("failed to create new map data %s", err)
	}
	fmt.Println("HERE 5, successfully wrote file")

	f, err := os.Create(filepath.Join("internal", "assets", tileSet.Name+"tilemap.go"))
	if err != nil {
		return fmt.Errorf("failed to create tilemap.go file %s", err)
	}

	err = paletteTemplate.ExecuteTemplate(f, "tilemap.go.tmpl", &tileMapData{
		Name:        tileMap.Name,
		Size:        "stil to do",
		TileCount:   len(tiles),
		TileSetName: tileSet.Name,
	})
	if err != nil {
		return fmt.Errorf("failed to execulte palette template %s", err)
	}

	return nil
}
