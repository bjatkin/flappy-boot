package sprite

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// MaxAttrs is the maximum number of OAM attrs that can be stored in memory at one time
const MaxAttrs = 128

var (
	// Block0 is the first block of sprite data, it contains sprites 0 - 255 if the sprites use the
	// 256 color palette, and sprites 0 - 511 if the sprites use the 16 color palettes
	Block0 = memmap.VRAM[4*memmap.CharBlockOffset : 5*memmap.CharBlockOffset]

	// Block1 is the second block of sprite data, it contains sprites 256 - 511 if the sprites use the
	// 256 color palette, and sprites 512 - 1023 if the sprites use the 16 color palettes
	Block1 = memmap.VRAM[5*memmap.CharBlockOffset:]

	// Palette is the 256 color palette memory for sprites
	Palette = memmap.Palette[0x0100:]

	// Palette0 is the zero-ith 16 color palette for sprites
	Palette0 = memmap.Palette[0x0100:0x0110]

	// Palette1 is the first 16 color palette for sprites
	Palette1 = memmap.Palette[0x0110:0x0120]

	// Palette2 is the second 16 color palette for sprites
	Palette2 = memmap.Palette[0x0120:0x0130]

	// Palette3 is the third 16 color palette for sprites
	Palette3 = memmap.Palette[0x0130:0x0140]

	// Palette4 is the fourth 16 color palette for sprites
	Palette4 = memmap.Palette[0x0140:0x0150]

	// Palette5 is the fifth 16 color palette for sprites
	Palette5 = memmap.Palette[0x0150:0x0160]

	// Palette6 is the sixth 16 color palette for sprites
	Palette6 = memmap.Palette[0x0160:0x0170]

	// Palette7 is the seventh 16 color palette for sprites
	Palette7 = memmap.Palette[0x0170:0x0180]

	// Palette8 is the eighth 16 color palette for sprites
	Palette8 = memmap.Palette[0x0180:0x0190]

	// Palette9 is the nineth 16 color palette for sprites
	Palette9 = memmap.Palette[0x0190:0x01A0]

	// PaletteA is the tenth 16 color palette for sprites
	PaletteA = memmap.Palette[0x01A0:0x01B0]

	// PaletteB is the eleventh 16 color palette for sprites
	PaletteB = memmap.Palette[0x01B0:0x01C0]

	// PaletteC is the twelfth 16 color palette for sprites
	PaletteC = memmap.Palette[0x01C0:0x01D0]

	// PaletteD is the thirteenth 16 color palette for sprites
	PaletteD = memmap.Palette[0x01D0:0x01E0]

	// PaletteE is the fourteenth 16 color palette for sprites
	PaletteE = memmap.Palette[0x01E0:0x01F0]

	// PaletteF is the fifteenth 16 color palette for sprites
	PaletteF = memmap.Palette[0x01F0:0x0200]
)

// OAM contains all the regular sprite data, it can hold up to 128 sprites, note that only 96 sprites can be drawn on a given horizontal line
var OAM = *((*[]Attrs)(unsafe.Pointer(&memmap.OAM)))

// AffineOAM contains all the affine sprite data, it can hold up to 32 affine sprite attributes, note that the affine sprite index must be
// set using the regular sprite data
var AffineOAM = *((*[]AffineAttrs)(unsafe.Pointer(&memmap.OAM)))

type (
	// Attr0 is the type of the first attribute in the Attrs struct
	Attr0 memmap.OAMValue

	// Attr1 is the type of the second attribute in the Attrs struct
	Attr1 memmap.OAMValue

	// Attr2 is the type of the third attribute in the Attrs struct
	Attr2 memmap.OAMValue
)

// Attrs is the structure of the sprite OAM attribute, it includes the 3 seperate attributes
// used for controlling a sprite, including size, location, color mode and others
type Attrs struct {
	// Attr0 has the following format
	//
	// [0 - 7] Y Position - y position of the top left corner of the sprite (0 - 255)
	//
	// [8 - 9] Sprite Mode - Set the draw mode of the sprite
	//   - Normal - Sprite is rendered normally (default)
	//   - Affine - Sprite is affine and rendered using the affine matrix
	//   - Hide - The Sprite is not drawn
	//   - AffineDBL - Affine sprite using double rendering area
	//
	// [A - B] Sprite Effect - Set the sprite draw effect
	//   - Normal - Sprite is rendered normally (default)
	//   - Blend - Sprite is rendered with alpha blending
	//   - Window - Sprite is used as a mask
	//
	// [C] Mosiac - Enables the Mosiac graphical effect
	//   - Mosaic - Sprite is rendered using the mosaic effect
	//
	// [D] Color Mode - Sets the sprite color mode
	//   - Color16 - use one of the 16, 16 color palettes when rendering (default)
	//   - Color256 - use the 256 color palette when rendering
	//
	// [E - F] Sprite Shape - Sets the shape of the sprite, combined with the sprite size to get the final sprite size
	//   - Square - the sprite is a square sprite
	//   - Wide - the sprite is wider than it is tall
	//   - Tall - this sprite is taller than it is wide
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
	Attr0 Attr0

	// Attr1 has the following format
	//
	// [0 - 8] X - The position of the top left corner of the sprite (0 - 511)
	//
	// [9 - D] Affine Index - The index for affine sprite data (0 - 32), only used if Attr0 is set to use affine attributes
	//
	// [C] Horizontal Mirrior - if set, the sprite is mirriored horizontally
	//
	// [D] Vertical Mirrior - if set, the sprite is mirriored vertically
	//
	// [E - F] Sprite Size - the size of the sprite, combined with the sprite size to get the final sprite size
	//   - Small - a small sprite, 8px to 16px in width and height
	//   - Medium - a medium sprite, 8px to 32px in width and height
	//   - Large - a large sprite, 16px to 32px in width and height
	//   - XL - an extra large sprite, 32px to 64px in width and height
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
	Attr1 Attr1

	// Attr2 has the following format
	//
	// [0 - 9] Tile Index - the index of the base tile for the sprite, starts at 512 in bit map modes(0 - 1024)
	//
	// [A - B] Priority - sets the priority/ layer of the sprites
	//   - priority0 - the highest sprite priorty, will be drawn above all other sprites
	//   - Priority1 - priority 1, will be drawn above priority 2 & 3 and below prioritiy 0
	//   - Priority2 - priority 2, will be drawn above priorities 3 and below priorities 0 & 1
	//   - Priority3 - the lowest priority, will be drawing below all other sprites
	//
	// [C - F] Palette Bank - the index of the 16 bit palette to use for the sprite, this will be ignored if the sprite is using 256 colors (0 - 16)
	Attr2 Attr2

	// _ is used for struct spacing because regular and affine OAM data is interlaced
	_ memmap.OAMValue
}

const (

	// Regular sets a sprite to be rendered normally
	Normal Attr0 = 0x0000

	// Affine sets a sprite to use the affine tranformation matrix when rendering
	Affine Attr0 = 0x0100

	// Hide prevents the sprite from being rendered (i.e. it is hidden)
	Hide Attr0 = 0x0200

	// AffineDBL sets a sprite to use double affine rendering
	AffineDBL Attr0 = 0x0300

	// Blend sets the sprite to use alpha blending when rendering
	Blend Attr0 = 0x0400

	// Win turns the sprite into a clipping mask for other sprites and backgrounds
	Win Attr0 = 0x0800

	// Mosaic sets the sprite to be rendered using mosaic mode
	Mosaic Attr0 = 0x1000

	// Color16 sets the sprite to use on of the 16, 16 color palettes
	Color16 Attr0 = 0x0000

	// Color256 sets the sprite to use the 256 color palettes
	Color256 Attr0 = 0x2000

	// Sqare sets the sprites shape to square
	Square Attr0 = 0x0000

	// Wide sets the sprites shape to be wider than it is tall
	Wide Attr0 = 0x4000

	// Tall sets the sprites shape to be taller than it is wide
	Tall Attr0 = 0x8000

	// YMask masks out all the bits from Attr0 that are not part of the y position attribute
	YMask Attr0 = 0x00FF

	// ShapMask masks out all the bits from Attr0 that are not part of the sprite shape
	ShapeMask Attr0 = 0xC000

	// SpriteModeMask masks out all the bits from Attr0 that are not part of the sprite mode
	SpriteModeMask Attr0 = 0x0300
)

const (
	// HMirrior mirriors the sprite horizontally
	HMirrior Attr1 = 0x1000

	// VMirrior mirriors the sprite vertically
	VMirrior Attr1 = 0x2000

	// Small sets the sprites size as small
	Small Attr1 = 0x0000

	// Medium sets the sprites size as medium
	Medium Attr1 = 0x4000

	// Large sets the sprites size as large
	Large Attr1 = 0x8000

	// XL sets the sprites size as extra large
	XL Attr1 = 0xC000

	// XMask masks out all the bits from Attr1 that are not part of the x position attribute
	XMask Attr1 = 0x01FF

	// AffineIndexMask masks out all the bits from Attr1 that are not part of the affine index
	AffineIndexMask Attr1 = 0x3E00

	// SizeMask masks out all the bits from Attr1 that are not part of the sprite size
	SizeMask Attr1 = 0xC000

	// HMirriorMask masks out all the bits from Attr1 that are not part of the sprite horizontal mirrior
	HMirriorMask Attr1 = 0x1000

	// VMirriorMask masks out all the bits from Attr1 that are not part of the sprite vertical mirrior
	VMirriorMask Attr1 = 0x2000
)

const (
	// Priority0 is the highest draw priority, sprites with priority 0 are drawn above all other sprites
	Priority0 Attr2 = 0x0000

	// Priority1 is the first priority level, sprites with priority 1 will be drawn above priority 2 & 3 sprites and below priority level 0
	Priority1 Attr2 = 0x0400

	// Priority2 is the second priority level, sprites with priority 2 will be drawn above priority 3 sprites and below priority levels 0 & 1
	Priority2 Attr2 = 0x0800

	// Priority3 is the lowest draw priority, sprites with priority 3 are drawn below all other sprites
	Priority3 Attr2 = 0x0C00

	// IndexMask masks out all the bits from Attr2 that are not part of the sprites index
	IndexMask Attr2 = 0x03FF

	// PalMask masks out all the bits from Attr2 that are not part of the pallet bank
	PalMask Attr2 = 0xF000

	// PalShift is the offset of the palette bank in Attr2
	PalShift Attr2 = 0x000C

	// PriorityMask masks out all the bits from Attr2 that are not part of the sprite priority
	PriorityMask Attr2 = 0x0C00
)

// AffineAttrs is the structure sprite OAM affine attributes, it maps the sprites pixels from screen space to the sprites pixel space
//
// [ Pa, Pb ] = [ Sx*cos(alpha),  Sy*sin(alpha) ]
// [ Pc, Pd ]   [ -Sx*sin(alpha), Sy*cos(alpha) ]
type AffineAttrs struct {
	// _ is for interlaced spacing, do not use this value
	_ [3]memmap.OAMValue

	// Pa is the 0,0 value of the matrix
	Pa memmap.OAMValue

	// _ is for interlaced spacing, do not use this value
	_ [3]memmap.OAMValue

	// Pb is the 0,1 value of the matrix
	Pb memmap.OAMValue

	// _ is for interlaced spacing, do not use this value
	_ [3]memmap.OAMValue

	// Pc is the 1,0 value of the matrix
	Pc memmap.OAMValue

	// _ is for interlaced spacing, do not use this value
	_ [3]memmap.OAMValue

	// Pd is the 1,1 value of the matrix
	Pd memmap.OAMValue
}
