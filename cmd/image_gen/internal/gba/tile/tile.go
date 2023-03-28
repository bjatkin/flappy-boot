package tile

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"sort"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbaimg"
	"golang.org/x/exp/maps"
)

// Size represents a valid GBA tile size
type Size struct {
	str    string
	width  int
	height int
	size   string
	shape  string
}

var (
	S8x8 = Size{
		str:    "8x8",
		width:  8,
		height: 8,
		size:   "sprite.Small",
		shape:  "sprite.Square",
	}
	S8x16 = Size{
		str:    "16x16",
		width:  8,
		height: 16,
		size:   "sprite.Small",
		shape:  "sprite.Tall",
	}
	S16x8 = Size{
		str:    "16x8",
		width:  16,
		height: 8,
		size:   "sprite.Small",
		shape:  "sprite.Wide",
	}
	S16x16 = Size{
		str:    "16x16",
		width:  16,
		height: 16,
		size:   "sprite.Medium",
		shape:  "sprite.Square",
	}
	S32x8 = Size{
		str:    "32x8",
		width:  32,
		height: 8,
		size:   "sprite.Medium",
		shape:  "sprite.Wide",
	}
	S8x32 = Size{
		str:    "8x32",
		width:  8,
		height: 32,
		size:   "sprite.Medium",
		shape:  "sprite.Tall",
	}
	S32x32 = Size{
		str:    "32x32",
		width:  32,
		height: 32,
		size:   "sprite.Large",
		shape:  "sprite.Square",
	}
	S32x16 = Size{
		str:    "32x16",
		width:  32,
		height: 16,
		size:   "sprite.Large",
		shape:  "sprite.Wide",
	}
	S16x32 = Size{
		str:    "16x32",
		width:  16,
		height: 32,
		size:   "sprite.Large",
		shape:  "sprite.Tall",
	}
	S64x64 = Size{
		str:    "64x64",
		width:  64,
		height: 64,
		size:   "sprite.XL",
		shape:  "sprite.Square",
	}
	S64x32 = Size{
		str:    "64x32",
		width:  64,
		height: 32,
		size:   "sprite.XL",
		shape:  "sprite.Wide",
	}
	S32x64 = Size{
		str:    "32x64",
		width:  32,
		height: 64,
		size:   "sprite.XL",
		shape:  "sprite.Tall",
	}
)

var allSizes = map[string]Size{
	S8x8.str:   S8x8,
	S8x16.str:  S8x16,
	S16x8.str:  S16x8,
	S16x16.str: S16x16,
	S32x8.str:  S32x8,
	S8x32.str:  S8x32,
	S32x32.str: S32x32,
	S32x16.str: S32x16,
	S16x32.str: S16x32,
	S64x64.str: S64x64,
	S64x32.str: S64x32,
	S32x64.str: S32x64,
}

// NewSize creates a new Size from a string
// strings should contain size width and height seperated by an x (e.g. 8x8, 16x32)
func NewSize(size string) (Size, error) {
	if s, ok := allSizes[size]; ok {
		return s, nil
	}
	return Size{}, fmt.Errorf("%s is not a valid tile size", size)
}

// Point returns the dimentions of the tile size as an image.Point
func (s Size) Point() image.Point {
	return image.Point{X: s.width, Y: s.height}
}

// Tiles returns the number of 8x8 tiles that make up a tile of this size
func (s Size) Tiles() int {
	return s.width * s.height / 64
}

// Size returns the sprite size
func (s Size) Size() string {
	return s.size
}

// Shape returns the sprite shape
func (s Size) Shape() string {
	return s.shape
}

// Is returns true if the to tile sizes have equal width and height
func (s Size) Is(comp Size) bool {
	return s.width == comp.width && s.height == comp.height
}

// String is the string representation of the tile size
func (s Size) String() string {
	return s.str
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
	if !m.Size.Is(S8x8) {
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
