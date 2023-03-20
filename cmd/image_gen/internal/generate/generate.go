package generate

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"path/filepath"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/config"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/raw"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/tile"
)

// File is an interface that can be used to create both raw data files and coresponding go files
type File interface {
	Raw() ([]byte, error)
	Go() ([]byte, error)
}

// PaletteData holds metadata for a specific palette
type PaletteData struct {
	Palette        color.Palette
	Name           string
	Description    string
	Shared         int
	SetTransparent *gbacol.RGB15
}

// NewPaletteData creates new palette data from a palette config
func NewPaletteData(palette config.Palette, setTransparent *gbacol.RGB15) (*PaletteData, error) {
	imgFile, err := os.Open(palette.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %w", err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s | %w", palette.File, err)
	}

	pal, err := raw.NewPal16(img, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new 16 color palette %w", err)
	}

	return &PaletteData{
		Palette:        pal,
		Name:           palette.Name,
		Description:    palette.Description,
		SetTransparent: setTransparent,
	}, nil
}

// Raw returns the raw palette data as a byte slice
func (t *PaletteData) Raw() ([]byte, error) {
	if t.SetTransparent != nil {
		return raw.Palette(append(color.Palette{*t.SetTransparent}, t.Palette[1:]...)), nil
	}

	return raw.Palette(t.Palette), nil
}

// Go returns a go file that contains the specified palette
func (t *PaletteData) Go() ([]byte, error) {
	b := &bytes.Buffer{}
	err := goTemplates.ExecuteTemplate(b, "palette.go.tmpl", t)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// TileSetData contains metadata for a specific tile set
type TileSetData struct {
	Name      string
	Tiles     []*tile.Meta
	TileCount int
	Length    int
	Bytes     int
	Palette   *PaletteData
	Size      tile.Size
	Shared    int
}

// NewTileSetData creates TileSetData from tileSet configuration and a map of PaletteData
func NewTileSetData(tileSet config.TileSet, setTransparent *gbacol.RGB15, palettes map[string]*PaletteData) (*TileSetData, error) {
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
		pal, err = NewPaletteData(config.Palette{
			File: tileSet.File,
		}, setTransparent)
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

	pal.Shared++
	tiles := tile.NewMetaSlice(img, pal.Palette, size)
	uniqueTiles := tile.Unique(tiles)

	return &TileSetData{
		Name:      tileSet.Name,
		Tiles:     uniqueTiles,
		TileCount: len(uniqueTiles),
		Length:    len(uniqueTiles) * 16,
		Bytes:     len(uniqueTiles) * 32,
		Palette:   pal,
		Size:      size,
	}, nil
}

// Raw returns the raw tile set data. If the tile set is the only user of its palette the
// palette data will be appended to the end of the tile set data as well
func (t *TileSetData) Raw() ([]byte, error) {
	raw := raw.Tiles(t.Tiles)

	if t.Palette.Shared == 1 {
		pal, err := t.Palette.Raw()
		if err != nil {
			return nil, err
		}

		raw = append(raw, pal...)
	}

	return raw, nil
}

// Go returns a go file that contains the tile set. If the tile set is the only user of its
// palette the palette data will also be contained in the go file
func (t *TileSetData) Go() ([]byte, error) {
	b := &bytes.Buffer{}
	err := goTemplates.ExecuteTemplate(b, "tileset.go.tmpl", t)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// TileMapData contains metadata for a specific tile map
type TileMapData struct {
	Name      string
	Width     int
	Height    int
	Tiles     []*tile.Meta
	TileCount int
	TileSet   *TileSetData
	Bytes     int
}

// BGSize returns the correct display.BGSize constant that corresponds to the given width and height
func (t *TileMapData) BGSize(width, height int) string {
	switch {
	case width > 32*8 && height > 32*8:
		return "display.BGSizeLarge"
	case width > 32*8:
		return "display.BGSizeWide"
	case height > 32*8:
		return "display.BGSizeTall"
	default:
		return "display.BGSizeSmall"
	}
}

// NewTileMapData creates a new TileMapData from a tileMap configuration and a tileSet and palette map
func NewTileMapData(tileMap config.TileMap, setTransparent *gbacol.RGB15, tileSets map[string]*TileSetData, palettes map[string]*PaletteData) (*TileMapData, error) {
	imgFile, err := os.Open(tileMap.File)
	if err != nil {
		return nil, fmt.Errorf("failed to read image file %s | %w", tileMap.File, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file %s | %w", tileMap.File, err)
	}

	var tileSet *TileSetData
	switch {

	// TODO: this section of code doesn't actually support using a custom palette yet. I should add that in
	case tileMap.TileSet == "" && tileMap.Palette == "":
		tileSet, err = NewTileSetData(config.TileSet{
			Name: tileMap.Name,
			File: tileMap.File,
			Size: "8x8",
		}, setTransparent, palettes)
		if err != nil {
			return nil, fmt.Errorf("failed to create tile set %w", err)
		}
	case tileMap.TileSet == "":
		tileSet, err = NewTileSetData(config.TileSet{
			Name:    tileMap.Name,
			File:    tileMap.File,
			Palette: tileMap.Palette,
			Size:    "8x8",
		}, setTransparent, palettes)
		if err != nil {
			return nil, fmt.Errorf("failed to create tile set %w", err)
		}
	default:
		var ok bool
		tileSet, ok = tileSets[tileMap.TileSet]
		if !ok {
			return nil, fmt.Errorf("tile set %s does not exists", tileMap.TileSet)
		}
	}

	if tileSet.Size != tile.S8x8 {
		return nil, fmt.Errorf("tile set size must be 8x8")
	}

	tileSet.Shared++
	tiles := tile.NewMetaSlice(img, tileSet.Palette.Palette, tile.S8x8)

	return &TileMapData{
		Name:      tileMap.Name,
		Width:     img.Bounds().Dx(),
		Height:    img.Bounds().Dy(),
		Tiles:     tiles,
		TileCount: len(tiles),
		TileSet:   tileSet,
		Bytes:     len(tiles) * 2,
	}, nil
}

// Raw returns the raw tile map data. If the tile map is the only user of it's tile set
// the tile set will be appended to the end of the data
func (t *TileMapData) Raw() ([]byte, error) {
	raw, err := raw.MapData(t.Tiles, t.TileSet.Tiles, t.Width, t.Height)
	if err != nil {
		return nil, err
	}

	if t.TileSet.Shared == 1 {
		ts, err := t.TileSet.Raw()
		if err != nil {
			return nil, err
		}
		raw = append(raw, ts...)
	}

	return raw, nil
}

// Go returns a go file that contains the the tile map. If the tile map is the only user of it's
// tile set the tile set data will also be included
func (t *TileMapData) Go() ([]byte, error) {
	b := &bytes.Buffer{}
	err := goTemplates.ExecuteTemplate(b, "tilemap.go.tmpl", t)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func WriteAssetFile(dir string) error {
	assetFile := filepath.Join(dir, "assets.go")
	f, err := os.Create(assetFile)
	if err != nil {
		return err
	}

	err = goTemplates.ExecuteTemplate(f, "assets.go.tmpl", nil)
	if err != nil {
		return err
	}

	return nil
}
