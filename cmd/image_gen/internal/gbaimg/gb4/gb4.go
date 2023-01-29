package gb4

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/byteconv"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
)

// NewPal16 converts an image into a valid color.Palette.
// it will also ensure the correct transparent option is at index 0 in the palette
// if the resulting palette contains more than 16 colors and error will be returned
func NewPal16(m image.Image, o *Options) (color.Palette, error) {
	pal := gbaimg.NewPal(m)
	if len(pal) > 16 {
		return nil, fmt.Errorf("palette is too large %d", len(pal))
	}

	for len(pal) < 16 {
		pal = append(pal, gbacol.RGB15(0x0000))
	}

	var tIndex int
	if o != nil && o.Transparent != nil {
		tIndex = pal.Index(o.Transparent)
	}

	if tIndex != 0 {
		pal[0], pal[tIndex] = pal[tIndex], pal[0]
	}

	return pal, nil
}

// Options allows the behavior of the Encode function to be adjusted
// the two options are TileSize, which can be used to adjust the size of tiles in the image
// and Transparent which can be used to explicitly set the transparent color in the color palette
type Options struct {
	TileSize    *tile.Size
	Transparent *gbacol.RGB15
	IncludeMap  bool
}

// Encode writes the image as a valid 4 bit rgb image (.gb4) to the provided writer
// the provided image has a maximum width and height of 256*TileSize pixels
// Any larger and the function will return an error
func Encode(w io.Writer, m image.Image, o *Options) error {
	dx, dy := m.Bounds().Dx(), m.Bounds().Dy()

	// 16 colors at two bytes each
	palSize := 16 * 2

	// this is the worst case tile size where every tile is toally unique so every pixel in the image stored
	tileSize := dx * dy

	// this is the worst case tile size where every tile is totally unique so every index will also be unique
	// each tile is at the smallest 8x8 pixles or 64 pixels in total which is why this number is divided by 64
	indexSize := (dx * dy) / 64

	// 4 bytes for the tile size
	// 4 byte for the width of the final image in tile
	// 4 byte for the height of the final image in tiles
	// 4 bytes for the number of unique tiles
	headerSize := 16

	raw := make([]byte, 0, headerSize+palSize+tileSize+indexSize)

	var size tile.Size // defaults to 8x8
	if o != nil {
		size = *o.TileSize
	}

	includeMap := o.IncludeMap
	if includeMap && size != tile.S8x8 {
		return fmt.Errorf("tile size must be 8x8 if a tile map is included")
	}

	pal, err := NewPal16(m, o)
	if err != nil {
		return err
	}

	tiles := tile.NewMetaSlice(m, pal, size)

	uniqueTiles := tile.Unique(tiles)
	tileCount := len(uniqueTiles)

	raw = append(raw, size.Bytes()...)
	raw = append(raw, 0, 0) // preserve alignment

	raw = append(raw, byteconv.Itoa(uint32(dx))...)
	raw = append(raw, byteconv.Itoa(uint32(dy))...)
	raw = append(raw, byteconv.Itoa(uint32(tileCount))...)

	for _, p := range pal {
		p16 := gbaimg.RGB15Model.Convert(p).(gbacol.RGB15)
		raw = append(raw, p16.Bytes()...)
	}

	for _, tile := range uniqueTiles {
		raw = append(raw, tile.Bytes()...)
	}

	// only include the map data if the tile size is 8x8 and saving map data was requested
	if includeMap {
		vFlip := uint16(0x0800)
		hFlip := uint16(0x0400)
		for _, tile := range tiles {
			for i, match := range uniqueTiles {
				var found bool
				switch {
				case gbaimg.Match(tile.Img, match.Img):
					raw = append(raw, byteconv.Itoa(uint16(i))...)
					found = true
				case gbaimg.Match(gbaimg.Flip(tile.Img, true, false), match.Img):
					raw = append(raw, byteconv.Itoa(uint16(i)|hFlip)...)
					found = true
				case gbaimg.Match(gbaimg.Flip(tile.Img, false, true), match.Img):
					raw = append(raw, byteconv.Itoa(uint16(i)|vFlip)...)
					found = true
				case gbaimg.Match(gbaimg.Flip(tile.Img, true, true), match.Img):
					raw = append(raw, byteconv.Itoa(uint16(i)|hFlip|vFlip)...)
					found = true
				}
				if found {
					break
				}
			}
		}
	}

	_, err = w.Write(raw)
	return err
}
