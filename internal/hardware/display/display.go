package display

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

const (
	// Width is the width of the LCD display in pixels
	Width = 240

	// Height is the width of the LCD display in pixels
	Height = 160
)

// Controll is the LCD controll register it can be used to set up and configure the display and other video data.
// The register is R/W with the exception of bit 3 and has the following bit layout.
//
// [0 - 2] Display Mode - Used to change the GBA's video mode
// [3] GBC - set to true if the game is a game boy color game (read only)
// [4] Active Page - Used to swap the active video page in mode 4 and 5
//   - PageA - The default video page in memory
//   - PageB - The secondary video page in memory
//
// [5] OAM HBlank Update Setting - Change wether OAM data can be updated durring an HBlank period
//   - OAMHBlankEnable - enable updating OAM data durring an hblank period
//   - OAMHBlankDisable - disable updating OAM data durring the hblank period
//
// [6] Sprite Tile Mapping Mode - the mapping mode for tiles in VRAM memory
//   - Sprite2D - use 2d object mapping for sprite data (default)
//   - Sprite1D - use 1d object mapping for sprite data
//
// [7] Force Blank - Force the screen to blank
//   - ForceBlank - force the screen to blank
//
// [8 - B] Background Enable - Enables background 0 - 3
//   - BG0 - Enable background 0
//   - BG1 - Enable background 1
//   - BG2 - Enable background 2
//   - BG3 - Enable background 3
//
// [C] Sprite Enable -Enables Sprites
//   - Sprite - Enables sprites
//
// [D - E] Window Enable - Used to enable windows for screen masking
//   - Window0 - Enable masking window 0
//   - Window1 -  Enable masking window 1
//
// [F] Sprite Window Enable - Use the sprite window for screen masking
//   - SprWindow - Enable the sprite masking window
var Controll = (*memmap.DisplayControll)(unsafe.Pointer(memmap.IOAddr + 0x0000))

const (
	// Mode0 is the first tile display mode
	//
	// Affine Backgrounds:  None
	// Regular Backgrounds: BG0, BG1, BG2, BG3
	Mode0 memmap.DisplayControll = 0x0000

	// Mode1 is the second tile display mode
	//
	// Affine Backgrounds:  BG2
	// Regular Backgrounds: BG0, BG1, BG2
	Mode1 memmap.DisplayControll = 0x0001

	// Mode2 is the third tile display mode
	//
	// Affine Backgrounds:  BG2, BG3
	// Regular Backgrounds: None
	Mode2 memmap.DisplayControll = 0x0002

	// Mode3 is the first bitmap mode and uses background 2
	//
	// Dimensions:     [240 x 160]
	// Bits per pixel: 16
	// Page Flipping:  No
	Mode3 memmap.DisplayControll = 0x0003

	// Mode4 is the second bitmap mode and uses background 2
	//
	// Dimensions:     [240 x 160]
	// Bits per pixel: 8
	// Page Flipping:  Yes
	Mode4 memmap.DisplayControll = 0x0004

	// Mode5 is the third bitmap mode and uses background 2
	//
	// Dimensions:     [160 x 128]
	// Bits per pixel: 16
	// Page Flipping:  Yes
	Mode5 memmap.DisplayControll = 0x0005

	// PageA is the default video page used by Mode4 and Mode5
	PageA memmap.DisplayControll = 0x0000

	// PageB is the alternative video page used by Mode4 and Mode5
	PageB memmap.DisplayControll = 0x0010

	// OAMHBlank allows the OAM(object attribute memory or sprite memory) to be updated
	// durring an HBlank, the normal behavior of the GBA prevents any updates to this
	// section of memory unless the screen is in the VBlank period.
	// Using this setting will reduce the number of sprites drawn per line
	OAMHBlank memmap.DisplayControll = 0x0020

	// Sprite1D sets the sprite memory layout to a single linear array of pixel data
	Sprite1D memmap.DisplayControll = 0x0040

	// Sprite2D is the default sprite mapping behavior and sets the sprite memory layout
	// to a 2d matrix consisting of 32x32 sprite tiles.
	Sprite2D memmap.DisplayControll = 0x0000

	// BG0 enables background 0
	BG0 memmap.DisplayControll = 0x0100

	// BG1 enables background 1
	BG1 memmap.DisplayControll = 0x0200

	// BG2 enables background 2
	BG2 memmap.DisplayControll = 0x0400

	// BG3 enables background 3
	BG3 memmap.DisplayControll = 0x0800

	// Sprites enables the sprite layer (somtimes refered to as the object layer)
	Sprites memmap.DisplayControll = 0x1000

	// Win1 enables the first background masking window
	Win1 memmap.DisplayControll = 0x2000

	// Win2 enables the second background masking window
	Win2 memmap.DisplayControll = 0x4000

	// WinSpr enables the sprite masking window
	WinSpr memmap.DisplayControll = 0x8000
)

// Stat is the LCD status controll register it can be use to read the display stats and controll
// line interrupts. It is R/W with the exception of bits 0-3 which are read only.
//
// [0] VBlank Flag - The hardware will set this to 1 when the screen is in a VBlank (read-only)
// [1] HBlank Flag - The hardware will set this to 1 when the screen is in an HBlank (read-only)
// [2] VCount Flag - The hardware will set this to 1 if the VCount setting matches the current vertical scan line
// [3] VBlank Interrupt - Enables the VBlank interrupt
//   - VBlankIRQ - Enables the VBlank interrupt
//
// [4] HBlank Interrupt - Enables the HBlank interrupt
//   - HBlankIRQ - Enables the HBlank interrupt
//
// [5] VCounter Interrupt - Enables the VCounter Interrupt
//   - VCounterIRQ - Enalbes the VCounter Interrupt when the VCount settings matches the current vertical scan line
//
// [9 - F] VCount Setting - The vertical scan line to match for VCount interrupts (0 - 227)
var Stat = (*memmap.DisplayStat)(unsafe.Pointer(memmap.IOAddr + 0x0004))

const (
	// VBlankMask masks out every bit that is not the VBlank flag bit
	VBlankMask memmap.DisplayStat = 0x0001

	// HBlankMask masks out every bit that is not the HBlank flag bit
	HBlankMask memmap.DisplayStat = 0x0002

	// VCountMask masks out every bit that is not the VCounter flag bit
	VCountMask memmap.DisplayStat = 0x0004

	// VBlankIRQ enables VBlank hardware interrupts
	VBlankIRQ memmap.DisplayStat = 0x0008

	// HBlankIRQ enables HBlank hardware interrupts
	HBlankIRQ memmap.DisplayStat = 0x0010

	// VCounterIRQ enables VCounter hardware interrupts
	VCounterIRQ memmap.DisplayStat = 0x0020

	// VCountShift shifts the value of of the Stat registers so only the vertical scan line setting remains
	VCountShift memmap.DisplayStat = 0x0009
)

var (
	// BG0Controll is the background 0 controll registers, it can be used to control various aspects
	// of background layer 0, background 0 can only be a regular background
	//
	// [0 - 1] Background Priority - Sets the Background 0 priority. 0 is the top priority and 3 is the lowest
	// [2 - 3] Character Base Block - Sets the character base block for tiled background modes
	// [6] Mosaic - Enables Mosaic Mode for background 0
	//   - Mosaic - Enable the mosaic mode for background 0
	//
	// [7] Color Mode - Sets the color mode for background 0
	//   - Color16 - Use 16 bit color for background 0 (default)
	//   - Color256 - Use 256 bit color for background 0
	//
	// [8 - C] Screen Base Block - Sets the screen base block for tiled background modes (0 - 31)
	// [E - F] Background Size - Sets the size for background 0
	//   - BGSizeSmall - 256 x 256 pixels
	//   - BGSizeWide - 512 x 256 pixels
	//   - BGSizeTall - 256 x 512 pixels
	//   - BGSizeLarge - 512 x 512 pixels
	BG0Controll = (*memmap.BGControll)(unsafe.Pointer(memmap.IOAddr + 0x0008))

	// BG1Controll is the background 1 controll registers, it can be used to control various aspects
	// of background layer 1, background 1 can only be a regular background
	//
	// [0 - 1] Background Priority - Sets the Background 1 priority. 0 is the top priority and 3 is the lowest
	// [2 - 3] Character Base Block - Sets the character base block for tiled background modes
	// [6] Mosaic - Enables Mosaic Mode for background 1
	//   - Mosaic - Enable the mosaic mode for background 1
	//
	// [7] Color Mode - Sets the color mode for background 1
	//   - Color16 - Use 16 bit color for background 1 (default)
	//   - Color256 - Use 256 bit color for background 1
	//
	// [8 - C] Screen Base Block - Sets the screen base block for tiled background modes (0 - 31)
	// [E - F] Background Size - Sets the size for background 1
	//   - BGSizeSmall - 256 x 256 pixels
	//   - BGSizeWide - 512 x 256 pixels
	//   - BGSizeTall - 256 x 512 pixels
	//   - BGSizeLarge - 512 x 512 pixels
	BG1Controll = (*memmap.BGControll)(unsafe.Pointer(memmap.IOAddr + 0x000A))

	// BG2Controll is the background 2 controll registers, it can be used to control various aspects
	// of the background layer 2, background 2 can be either an affine or a regular background, it is also
	// layer used by the non-tiled display modes
	//
	// [0 - 1] Background Priority - Sets the Background 2 priority. 0 is the top priority and 3 is the lowest
	// [2 - 3] Character Base Block - Sets the character base block for tiled background modes
	// [6] Mosaic - Enables Mosaic Mode for background 2
	//   - Mosaic - Enable the mosaic mode for background 2
	//
	// [7] Color Mode - Sets the color mode for background 2
	//   - Color16 - Use 16 bit color for background 2 (default)
	//   - Color256 - Use 256 bit color for background 2
	//
	// [8 - C] Screen Base Block - Sets the screen base block for tiled background modes (0 - 31)
	// [E - F] Background Size - Sets the size for background 2
	//   - BGSizeSmall - 256 x 256 pixels
	//   - BGSizeWide - 512 x 256 pixels
	//   - BGSizeTall - 256 x 512 pixels
	//   - BGSizeLarge - 512 x 512 pixels
	BG2Controll = (*memmap.BGControll)(unsafe.Pointer(memmap.IOAddr + 0x000C))

	// BG3Controll is the background 3 controll registers, it can be used to control various assets
	// of background layer 3, background 3 can be either an affine or a regular background
	//
	// [0 - 1] Background Priority - Sets the Background 3 priority. 0 is the top priority and 3 is the lowest
	// [2 - 3] Character Base Block - Sets the character base block for tiled background modes
	// [6] Mosaic - Enables Mosaic Mode for background 3
	//   - Mosaic - Enable the mosaic mode for background 3
	//
	// [7] Color Mode - Sets the color mode for background 3
	//   - Color16 - Use 16 bit color for background 3 (default)
	//   - Color256 - Use 256 bit color for background 3
	//
	// [8 - C] Screen Base Block - Sets the screen base block for tiled background modes (0 - 31)
	// [E - F] Background Size - Sets the size for background 3
	//   - BGSizeSmall - 256 x 256 pixels
	//   - BGSizeWide - 512 x 256 pixels
	//   - BGSizeTall - 256 x 512 pixels
	//   - BGSizeLarge - 512 x 512 pixels
	BG3Controll = (*memmap.BGControll)(unsafe.Pointer(memmap.IOAddr + 0x000E))
)

const (
	// Priority0 is the top priority for backgrounds, it will be drawn above all other
	// backgrounds and below sprites with priority 0, but above sprites with lower priorities
	Priority0 memmap.BGControll = 0x0000

	// Priority1 is priority 1 for backgrounds, it will be drawn above backgrounds with priorities 2 and 3
	// backgrounds and below sprites with priority 1, but above sprites with lower priorities
	Priority1 memmap.BGControll = 0x0001

	// Priority2 is priority 2 for backgrounds, it will be drawn above bakcgrouds with priority 3
	// backgrounds and below sprites with priority 2, but above sprites with lower priorities
	Priority2 memmap.BGControll = 0x0002

	// Priority3 is the lowest priority for backgrounds, it will be drawn below all other backgrounds
	// backgrounds and all below sprites
	Priority3 memmap.BGControll = 0x0003

	// Mosaic enables the mosaic background effect
	Mosaic memmap.BGControll = 0x0020

	// Color16 sets the background to use 16 x 16 color palettes
	Color16 memmap.BGControll = 0x0000

	// Color256 sets the background to use 1 x256 color palettes
	Color256 memmap.BGControll = 0x0040

	// BGSizeSmall sets the background size to 256 x 256 pixels
	BGSizeSmall memmap.BGControll = 0x0000

	// BGSizeWide sets the background size to 512 x 256 pixels
	BGSizeWide memmap.BGControll = 0x0000

	// BGSizeTall sets the background size to 256 x 512 pixels
	BGSizeTall memmap.BGControll = 0x0000

	// BGSizeLarge sets the background size to 512 x 512 pixels
	BGSizeLarge memmap.BGControll = 0x0000

	// SBBShift shifts a number into the correct bits to set the screen base block
	SBBShift memmap.BGControll = 0x0008
)
