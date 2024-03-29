package raw

import (
	"errors"
	"fmt"
	"image"
	"image/color"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/byteconv"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbaimg"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/tile"
)

// NewPal16 converts an image into a valid color.Palette.
// it will also ensure the correct transparent option is at index 0 in the palette
// if the resulting palette contains more than 16 colors and error will be returned
func NewPal16(m image.Image, transparent *gbacol.RGB15) (color.Palette, error) {
	pal := gbaimg.NewPal(m)
	if len(pal) > 16 {
		return nil, fmt.Errorf("palette is too large %d", len(pal))
	}

	for len(pal) < 16 {
		pal = append(pal, gbacol.RGB15(0x0000))
	}

	var tIndex int
	if transparent != nil {
		tIndex = pal.Index(transparent)
	}

	if tIndex != 0 {
		pal[0], pal[tIndex] = pal[tIndex], pal[0]
	}

	return pal, nil
}

// Palette converts a color.Palette into a byte slices for a .gb4 image
func Palette(pal color.Palette) []byte {
	var raw []byte
	for _, p := range pal {
		p16 := gbaimg.RGB15Model.Convert(p).(gbacol.RGB15)
		raw = append(raw, p16.Bytes()...)
	}
	return raw
}

// Tiles converts a slice of tile.Meta tiles to a raw byte slice for a .gb4 image
func Tiles(tiles []*tile.Meta) []byte {
	var raw []byte
	for _, tile := range tiles {
		raw = append(raw, tile.Bytes()...)
	}
	return raw
}

// MapData converts map data for a .gb4 image into a raw byte slice.
// tiles are mapped using 32x32 tile screen base blocks.
func MapData(tiles []*tile.Meta, uniqueTiles []*tile.Meta, dx, dy int) ([]byte, error) {
	// pitch is the number of tiles per horizontal row
	pitch := paddedPitch(dx)
	tileDy := dy / 8
	raw := make([]byte, pitch*tileDy*2) // multiply by 2 since each index will take 2 bytes

	// vFlip is used to indicate that the tile should be flipped vertically when drawn
	vFlip := uint16(0x0800)
	// hFlip is used to indicate that the tile should be flipped horzontally when drawn
	hFlip := uint16(0x0400)

	for i, tile := range tiles {
		if tile.IsTransparent() {
			index := getIndex(i, dx, dy, pitch)
			raw[index] = 0
			raw[index+1] = 0
			continue
		}

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
				index := getIndex(i, dx, dy, pitch)
				// +1 here because the 0th tile is reserved as a transparency tile
				tBytes := byteconv.Itoa(add + 1)

				raw[index] = tBytes[0]
				raw[index+1] = tBytes[1]
				break
			}

			if ii == len(uniqueTiles)-1 {
				return nil, errors.New("tile not found in tile set")
			}
		}
	}

	return raw, nil
}

// getIndex calculates the index in the buffer where the tile should be set
func getIndex(i, dx, dy, pitch int) int {
	tileX, tileY := i%(dx/8), i/(dx/8)
	screenBaseBlock := (tileY/32)*(pitch/32) + (tileX / 32)
	return (screenBaseBlock*1024 + (tileY%32)*32 + tileX%32) * 2
}

// paddedPitch returns the pitch of the image in tiles.
// it is padded to the nearest screen block
func paddedPitch(dx int) int {
	return ((255 + dx) / 256) * 32
}
