package gb8

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

// NewPal256 converts an image into a valid color.Palette.
// it will also ensure the correct transparent option is at index 0 in the palette
// if the resulting palette contains more than 16 colors and error will be returned
func NewPal256(m image.Image, o *Options) (color.Palette, error) {
	pal := gbaimg.NewPal(m)
	if len(pal) > 256 {
		return nil, fmt.Errorf("palette is too large %d", len(pal))
	}

	var transparent gbacol.RGB15
	if o != nil && o.Trans != nil {
		transparent = *o.Trans
	}

	if i := pal.Index(transparent); i != 0 {
		pal[0], pal[i] = pal[i], pal[0]
	}

	return pal, nil
}

// Options allows the behavior of the Encode function to be adjusted
// the two optoins are TileSize, which can be used to adjust the size of tiles in the image
// and Transparent which can be used to explicitly set the transparent color in the color palette
type Options struct {
	TileSize *tile.Size
	Trans    *gbacol.RGB15
}

// Encode writes the image as a valid 8 bit rgb image (.gb8) to the provided writer
// the provided image has a maximum width and height of 256*TileSize pixels
// Any larger and the funciton will return an error
func Encode(w io.Writer, m image.Image, o *Options) error {
	dx, dy := m.Bounds().Dx(), m.Bounds().Dy()

	// 256 colors at two bytes each
	palSize := 256 * 2

	// this is the worst case tile size where every tile is totally unique so every pixel in the image must be stored
	tileSize := dx * dy

	// this is the worst case tile size where every tile is totally unique so every index will also be unique
	// each tile is at the smallest 8x8 pixles or 64 pixels in total which is why this number is divided by 64
	indexSize := (dx * dy) / 64

	// 2 bytes for the tile size
	// 1 byte for the width of the final image in tiles
	// 1 byte for the height of the final image in tiles
	// 2 bytes for the number of unique tiles
	headerSize := 6

	raw := make([]byte, 0, headerSize+palSize+tileSize+indexSize)

	var size tile.Size // defaults to 8x8
	if o != nil {
		size = *o.TileSize
	}

	pal, err := NewPal256(m, o)
	if err != nil {
		return err
	}

	tiles := tile.NewMetaSlice(m, pal, size)

	uniqueTiles := tile.Unique(tiles)
	tileCount := len(uniqueTiles)
	if tileCount > 0xFFFF {
		return fmt.Errorf("too many unique tiles in the image %d", tileCount)
	}

	raw = append(raw, size.Bytes()...)
	raw = append(raw, byte(dx/8), byte(dy/8))
	raw = append(raw, byteconv.Itoa(uint16(tileCount))...)

	for _, p := range pal {
		p256 := gbaimg.RGB15Model.Convert(p).(gbacol.RGB15)
		raw = append(raw, p256.Bytes()...)
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
