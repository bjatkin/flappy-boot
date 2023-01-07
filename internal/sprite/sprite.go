package sprite

import (
	"embed"
	"fmt"
	"io"
	"math/rand"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

// LoadPalette256 loads a 256 color palette from an embedded file into memory
func LoadPalette256(fs embed.FS, name string) error {
	palette, err := fs.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open palette file %s: %s", name, err)
	}

	pal, err := io.ReadAll(palette)
	if err != nil {
		return fmt.Errorf("failed to read palette file: %s", err)
	}

	memmap.Copy16(hw_sprite.Palette, pal)
	return nil
}

// LoadPalette16 loads a 16 color palette from an embedded file into memory
func LoadPalette16(fs embed.FS, name string, palIndex int) error {
	if palIndex > 0x0010 {
		return fmt.Errorf("palette bank %d does not exist must be 0-16", palIndex)
	}

	palette, err := fs.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open palette file %s: %s", name, err)
	}

	pal, err := io.ReadAll(palette)
	if err != nil {
		return fmt.Errorf("failed to read palette file: %s", err)
	}

	start := 0x0100 + 0x0010*palIndex
	memmap.Copy16(memmap.Palette[start:start+0x0010], pal)

	return nil
}

// buffer is the sprite buffer used to double buffer sprite data
var buffer = make(map[uint32]*Instance, hw_sprite.MaxAttrs)

// Reset deletes all the current sprite instances from the screen and and prepares OAM data
// to hold new sprite data. It should be called as part of unloading any node that uses sprites
func Reset() {
	// reset the buffer
	buffer = make(map[uint32]*Instance, hw_sprite.MaxAttrs)

	// move all the sprites off screen and hide them to make sure old sprites don't get drawn
	for i := range hw_sprite.OAM {
		hw_sprite.OAM[i].Attr0 = hw_sprite.Attr0(255) | hw_sprite.Hide
		hw_sprite.OAM[i].Attr1 = hw_sprite.Attr1(511)
	}
}

// CopyOAM coppies all the current sprites into the OAM buffer so they can be drawn by the hardware
//
// NOTE: only 128 sprtes can be drawn by the hardware in a given frame. Sprites will be randomly
// shuffled so that if more than 128 sprites are being used they will be flickered
func CopyOAM() {
	for i := range buffer {
		// only copy up to 128 sprites into OAM data as that's the maximum allowed by the hardware
		if i > hw_sprite.MaxAttrs {
			return
		}

		spr := buffer[i]
		var attrs hw_sprite.Attrs

		// convert a sprite into it's hardware representation
		attrs.Attr0 = hw_sprite.Attr0(spr.Y) & hw_sprite.YMask
		attrs.Attr0 |= hw_sprite.Attr0(spr.Shape)
		if spr.Hide {
			attrs.Attr0 |= hw_sprite.Hide
		}
		if spr.Mosaic {
			attrs.Attr0 |= hw_sprite.Mosaic
		}
		if spr.Color256 {
			attrs.Attr0 |= hw_sprite.Color256
		}

		attrs.Attr1 = hw_sprite.Attr1(spr.X) & hw_sprite.XMask
		attrs.Attr1 = hw_sprite.Attr1(spr.Size)
		if spr.HFlip {
			attrs.Attr1 |= hw_sprite.HMirrior
		}
		if spr.VFlip {
			attrs.Attr1 |= hw_sprite.VMirrior
		}

		attrs.Attr2 = hw_sprite.Attr2(spr.TileIndex) & hw_sprite.IndexMask
		attrs.Attr2 |= hw_sprite.Attr2(spr.Priority)
		attrs.Attr2 |= hw_sprite.Attr2(spr.PaletteBank) << hw_sprite.PalShift

		// copy it into the OAM data
		hw_sprite.OAM[i] = attrs
	}
}

// Point is a 2d point it space
type Point struct {
	X, Y int
}

// Priority is the draw priority of the sprite
type Priority uint16

var (
	// Priority0 is the highest draw priority, sprites with priority 0 are drawn above all other sprites
	Priority0 = Priority(hw_sprite.Priority0)

	// Priority1 is the first priority level, sprites with priority 1 will be drawn above priority 2 & 3 sprites and below priority level 0
	Priority1 = Priority(hw_sprite.Priority1)

	// Priority2 is the second priority level, sprites with priority 2 will be drawn above priority 3 sprites and below priority levels 0 & 1
	Priority2 = Priority(hw_sprite.Priority2)

	// Priority3 is the lowest draw priority, sprites with priority 3 are drawn below all other sprites
	Priority3 = Priority(hw_sprite.Priority3)
)

// Shape is the shape of the sprite
type Shape uint16

var (
	// Square is a square sprite
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Square = Shape(hw_sprite.Square)

	// Wide is a sprite that is wider than it is tall
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Wide = Shape(hw_sprite.Wide)

	// Tall is a sprite that is taller than it is wide
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Tall = Shape(hw_sprite.Tall)
)

// Size is the size of the sprite
type Size uint16

var (
	// Small is a small sprite
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Small = Size(hw_sprite.Small)

	// Medium is a medium sized sprite
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Medium = Size(hw_sprite.Medium)

	// Large is a large sprite
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	Large = Size(hw_sprite.Large)

	// XL is an extra large sprite
	//
	// Sprite shape and size in pixels are determinded by both their size and shape attributes
	//   Square/ Small - 8  x 8
	//   Wide  / Small - 16 x 8
	//   Tall  / Small - 8  x 16
	//
	//   Square/ Med - 16 x 16
	//   Wide  / Med - 32 x 8
	//   Tall  / Med - 8  x 32
	//
	//   Square/ Large - 32 x 32
	//   Wide  / Large - 32 x 16
	//   Tall  / Large - 16 x 32
	//
	//   Square/ XL - 64 x 64
	//   Wide  / XL - 64 x 32
	//   Tall  / XL - 32 x 64
	XL = Size(hw_sprite.XL)
)

// Instance is an instance of a sprite
type Instance struct {
	Point

	id uint32

	Hide       bool
	Blend      bool
	WindowMask bool
	Mosaic     bool

	Color256 bool

	Shape Shape
	Size  Size

	HFlip, VFlip bool

	Priority    Priority
	PaletteBank uint

	TileIndex int
}

// NewInstance crates a new instance of a sprite at the X, Y position on the screen and the
// pixel data from the specified tile index.
//
// NOTE: that only 128 sprites can be shown by the GBA in a given frame.
// If more than 128 sprites are crated then the sprites will begin to flicker. You can prevent
// this by using Delete() on un-needed sprites
func NewInstance(x, y, TileIndex int) *Instance {
	i := &Instance{
		Point: Point{
			X: x,
			Y: y,
		},

		id:        rand.Uint32(),
		TileIndex: TileIndex,
	}

	buffer[i.id] = i

	return i
}

// Delete deletes a sprite. After being deleted a sprte instance should be discarded as
// it will no longer be drawn to the screen
func Delete(sprite *Instance) {
	delete(buffer, sprite.id)
}

// ID returns the imutable ID of the sprite instance
func (s *Instance) ID() uint32 {
	return s.id
}
