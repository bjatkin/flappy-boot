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
	var raw []byte

	size := tile.S8x8
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

	raw = append(raw, generateHeader(size, dx, dy, len(uniqueTiles))...)
	raw = append(raw, rawPalette(pal)...)
	raw = append(raw, rawTiles(uniqueTiles)...)

	if includeMap {
		raw = append(raw, rawMapData(tiles, uniqueTiles, dx, dy)...)
	}

	_, err = w.Write(raw)
	return err
}

// generateHeader generates raw header data for a .gb4 image
func generateHeader(size tile.Size, dx, dy, tileCount int) []byte {
	var header []byte
	header = append(header, size.Bytes()...)

	// preserve alignment since a tile size is only 2 bytes in length
	header = append(header, 0, 0)

	header = append(header, byteconv.Itoa(uint32(((256+dx)/246)*256))...)
	header = append(header, byteconv.Itoa(uint32(dy))...)
	header = append(header, byteconv.Itoa(uint32(tileCount))...)

	return header
}

// rawPalette converts a color.Palette into a byte slices for a .gb4 image
func rawPalette(pal color.Palette) []byte {
	var raw []byte
	for _, p := range pal {
		p16 := gbaimg.RGB15Model.Convert(p).(gbacol.RGB15)
		raw = append(raw, p16.Bytes()...)
	}
	return raw
}

// rawTiles converts a slice of tile.Meta tiles to a raw byte slice for a .gb4 image
func rawTiles(tiles []*tile.Meta) []byte {
	var raw []byte
	for _, tile := range tiles {
		raw = append(raw, tile.Bytes()...)
	}
	return raw
}

// rawMapData converts map data for a .gb4 image into a raw byte slice.
// tiles are mapped using 32x32 tile screen base blocks.
func rawMapData(tiles []*tile.Meta, uniqueTiles []*tile.Meta, dx, dy int) []byte {
	// pitch is the number of tiles per horizontal row
	pitch := ((256 + dx) / 256) * 32
	raw := make([]byte, pitch*dy)

	// vFlip is used to indicate that the tile should be flipped vertically when drawn
	vFlip := uint16(0x0800)
	// hFlip is used to indicate that the tile should be flipped horzontally when drawn
	hFlip := uint16(0x0400)

	for i, tile := range tiles {
		for ii, match := range uniqueTiles {
			var found bool
			var add uint16

			// all flipped orientations need to be checked since the unique tile set may be optimized to
			// take advantage of flip bits
			switch {
			case gbaimg.Match(tile.Img, match.Img):
				add = uint16(ii)
				found = true
			case gbaimg.Match(gbaimg.Flip(tile.Img, true, false), match.Img):
				add = uint16(ii) | hFlip
				found = true
			case gbaimg.Match(gbaimg.Flip(tile.Img, false, true), match.Img):
				add = uint16(ii) | vFlip
				found = true
			case gbaimg.Match(gbaimg.Flip(tile.Img, true, true), match.Img):
				add = uint16(ii) | hFlip | vFlip
				found = true
			}

			if found {
				tileX, tileY := i%(dx/8), i/(dx/8)
				screenBaseBlock := (tileY/32)*(pitch/32) + (tileX / 32)
				index := (screenBaseBlock*1024 + (tileY%32)*32 + tileX%32) * 2
				fmt.Println("\t-", tileX, tileY, "|", index)
				tBytes := byteconv.Itoa(add)
				raw[index] = tBytes[0]
				raw[index+1] = tBytes[1]
				break
			}
		}
	}

	return raw
}
