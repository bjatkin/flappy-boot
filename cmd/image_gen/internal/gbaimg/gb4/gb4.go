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
	"golang.org/x/exp/slices"
)

// NewPal16 converts an image into a valid color.Palette.
// it will also ensure the correct transparent option is at index 0 in the palette
// if the resulting palette contains more than 16 colors and error will be returned
func NewPal16(m image.Image, o *Options) (color.Palette, error) {
	pal := gbaimg.NewPal(m)
	if len(pal) > 16 {
		return nil, fmt.Errorf("palette is too large %d", len(pal))
	}

	var transparent gbacol.RGB15
	if o != nil && o.Transparent != nil {
		transparent = *o.Transparent
	}

	if i := pal.Index(transparent); i != 0 {
		pal[0], pal[i] = pal[i], pal[0]
	}

	return pal, nil
}

// Options allows the behavior of the Encode function to be adjusted
// the two options are TileSize, which can be used to adjust the size of tiles in the image
// and Transparent which can be used to explicitly set the transparent color in the color palette
type Options struct {
	TileSize    *tile.Size
	Transparent *gbacol.RGB15
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

	pal, err := NewPal16(m, o)
	if err != nil {
		return err
	}

	tiles := tile.NewMetaSlice(m, pal, size)

	uniqueTiles := tile.Unique(tiles)
	tileCount := len(uniqueTiles)

	raw = append(raw, size.Bytes()...)
	raw = append(raw, 0, 0) // preserver alignment

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

	for _, tile := range tiles {
		raw = append(raw, byte(slices.Index(uniqueTiles, tile)))
	}

	_, err = w.Write(raw)
	return err
}
