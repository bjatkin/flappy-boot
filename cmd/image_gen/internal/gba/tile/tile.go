package tile

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/byteconv"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbaimg"
	"golang.org/x/exp/maps"
)

// TODO: it would be nice to be able to auto-generate a lot of this code

// TileSize represents a valid GBA tile size
type Size uint16

// valid tile sizes
const (
	// TODO make these bitfields that match the ones in the sprite attrs
	S8x8   Size = 0x08_08
	S16x8  Size = 0x10_08
	S8x16  Size = 0x08_10
	S16x16 Size = 0x10_10
	S32x8  Size = 0x20_08
	S8x32  Size = 0x08_20
	S32x32 Size = 0x20_20
	S32x16 Size = 0x20_10
	S16x32 Size = 0x10_20
	S64x64 Size = 0x40_40
	S64x32 Size = 0x40_20
	S32x64 Size = 0x20_40
)

// NewSize creates a new Size from a string
// strings should contain size width and height seperated by an x (e.g. 8x8, 16x32)
func NewSize(size string) (Size, error) {
	switch size {
	case "8x8":
		return S8x8, nil
	case "8x16":
		return S8x16, nil
	case "16x8":
		return S16x8, nil
	case "16x16":
		return S16x16, nil
	case "32x8":
		return S32x8, nil
	case "8x32":
		return S8x32, nil
	case "32x32":
		return S32x32, nil
	case "32x16":
		return S32x16, nil
	case "16x32":
		return S16x32, nil
	case "64x64":
		return S64x64, nil
	case "64x32":
		return S64x32, nil
	case "32x64":
		return S32x64, nil
	default:
		return 0, fmt.Errorf("%s is not a valid tile size", size)
	}
}

// Point returns the Size as an image.Point
func (s Size) Point() image.Point {
	switch s {
	case S8x8:
		return image.Point{X: 8, Y: 8}
	case S8x16:
		return image.Point{X: 8, Y: 16}
	case S16x8:
		return image.Point{X: 16, Y: 8}
	case S16x16:
		return image.Point{X: 16, Y: 16}
	case S32x8:
		return image.Point{X: 32, Y: 8}
	case S8x32:
		return image.Point{X: 8, Y: 32}
	case S32x32:
		return image.Point{X: 32, Y: 32}
	case S32x16:
		return image.Point{X: 32, Y: 16}
	case S16x32:
		return image.Point{X: 16, Y: 32}
	case S64x64:
		return image.Point{X: 64, Y: 64}
	case S64x32:
		return image.Point{X: 64, Y: 32}
	case S32x64:
		return image.Point{X: 32, Y: 64}
	default:
		return image.Point{X: 8, Y: 8}
	}
}

// Tiles returns the number of 8x8 tiles the make up a tile of this size
func (s Size) Tiles() int {
	switch s {
	case S8x8:
		return 1
	case S8x16:
		return 2
	case S16x8:
		return 2
	case S16x16:
		return 4
	case S32x8:
		return 4
	case S8x32:
		return 4
	case S32x32:
		return 16
	case S32x16:
		return 8
	case S16x32:
		return 8
	case S64x64:
		return 64
	case S64x32:
		return 32
	case S32x64:
		return 32
	default:
		return 1
	}
}

// Bytes returns the Size as a byte slice
func (s Size) Bytes() []byte {
	return byteconv.Itoa(uint16(s))
}

// Meta is a meta tile. It consists of smaller 8x8 tiles and must use either a 16 color or 256 color palette
type Meta struct {
	// Size is the size of the full meta tile
	Size Size

	// Img is the raw image data of the full meta tile
	Img image.Image

	// Pal is the color.Palette used by the meta tile, the Img will be converted to use this color space
	Pal color.Palette

	// Tiles is the smaller 8x8 tiles that make up the meta tile
	Tiles []image.Image
}

// NewMeta creates a new meta tile with the given image and palette.
// if the provided image is different from the provided size, the image will be cropped or padded so the size matches
func NewMeta(img image.Image, pal color.Palette, size Size) *Meta {
	var tiles []image.Image
	gbaimg.WalkN(img, image.Point{X: 8, Y: 8}, func(x, y int) {
		tile := image.NewRGBA(image.Rect(0, 0, 8, 8))
		gbaimg.Copy(gbaimg.SubImage(img, image.Rect(x, y, x+8, y+8)), tile)

		tiles = append(tiles, tile)
	})

	return &Meta{
		Size:  size,
		Img:   img,
		Pal:   pal,
		Tiles: tiles,
	}
}

// Hash, provides a single hash for the meta tile
// The has will be consistent for tiles that share the same color data
// as well as for all all tiles that are simply mirriors of each other
func (m *Meta) Hash() [md5.Size]byte {
	hash := func(img image.Image) [md5.Size]byte {
		var payload []byte

		gbaimg.Walk(img, func(x, y int) {
			c := gbaimg.RGB15Model.Convert(img.At(x, y)).(gbacol.RGB15)
			payload = append(payload, c.Bytes()...)
		})

		return md5.Sum(payload)
	}

	hashes := [][md5.Size]byte{
		hash(m.Img),
		hash(gbaimg.Flip(m.Img, true, false)),
		hash(gbaimg.Flip(m.Img, false, true)),
		hash(gbaimg.Flip(m.Img, true, true)),
	}

	sum := func(data [16]byte) int {
		var total int
		for _, d := range data {
			total += int(d)
		}
		return total
	}

	sort.Slice(hashes, func(i, j int) bool {
		return sum(hashes[i]) > sum(hashes[j])
	})

	return hashes[0]
}

// Bytes returns the meta tile as bytes
// each tile within the meta tile has it's color data converted into the palettes color space and from there into palette indexes
// these indexes with take 8 bits each if the color pallete is contains more than 16 colors
// otherwise indexes will only take on nibble (4 bits) each.
func (m *Meta) Bytes() []byte {
	var data []byte
	for _, tile := range m.Tiles {
		gbaimg.Walk(tile, func(x, y int) {
			col := tile.At(x, y)
			data = append(data, byte(m.Pal.Index(col)))
		})
	}

	if len(m.Pal) <= 16 {
		var nibbles []byte
		for i := 0; i < len(data); i += 2 {
			nibbles = append(nibbles, data[i+1]<<4|data[i])
		}
		data = nibbles
	}

	return data
}

// IsTransparent returns true if the tile is 8x8 and completely transparent.
// This is usefule for building layerd backgrounds where many tiles can be fully transparent
func (m *Meta) IsTransparent() bool {
	if m.Size != S8x8 {
		return false
	}

	var hasColor bool
	transparent := gbaimg.RGB15Model.Convert(m.Pal[0])
	gbaimg.Walk(m.Img, func(x, y int) {
		if gbaimg.RGB15Model.Convert(m.Img.At(x, y)) != transparent {
			hasColor = true
		}
	})

	return !hasColor
}

// NewMetaSlice creates a new slice of meta tiles from an underlying image
// Each Meta tile will have the specified size and share the provided palette
// if the image size is not evenly divisible by the provided size the image will be padded to include all the image data
// for example an image that is 10x10 and uses a size of 8x8 will end up padding the provided image to be 16x16
func NewMetaSlice(img image.Image, pal color.Palette, size Size) []*Meta {
	pt := size.Point()

	var metas []*Meta
	gbaimg.WalkN(img, pt, func(x, y int) {
		tile := image.NewRGBA(image.Rect(0, 0, pt.X, pt.Y))
		gbaimg.Copy(gbaimg.SubImage(img, image.Rect(x, y, x+pt.X, y+pt.Y)), tile)

		metas = append(metas, NewMeta(tile, pal, size))
	})

	return metas
}

// Unique returns only the unique tiles in the slice of meta tiles
// tiles that only differ in that they are mirriors of one another are not considerd to be unique from each other
func Unique(tiles []*Meta) []*Meta {
	unique := make(map[[md5.Size]byte]*Meta)
	for _, tile := range tiles {
		unique[tile.Hash()] = tile
	}
	uniqueTiles := maps.Values(unique)

	sum := func(data [16]byte) int {
		var total int
		for _, d := range data {
			total += int(d)
		}
		return total
	}

	// sort the tiles slice for consistent indexes
	sort.Slice(uniqueTiles, func(i, j int) bool {
		return sum(uniqueTiles[i].Hash()) > sum(uniqueTiles[j].Hash())
	})

	return uniqueTiles
}
