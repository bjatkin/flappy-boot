package ppu

import (
	"image"
	"image/color"

	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

// transparent is the transparent color to use for RGB15 images
const transparent memmap.PaletteValue = 0x7C1F

// RGB15 is an image that holds memmap.PaletteValue colors directly
type RGB15 struct {
	colors []memmap.PaletteValue
	width  int
	height int
}

// NewRGB15 creates a new RGB15 image with the given width and height
func NewRGB15(width, height int) *RGB15 {
	return &RGB15{
		colors: make([]memmap.PaletteValue, width*height),
		width:  width,
		height: height,
	}
}

// DrawImage draws the given RGB15 image onto the target RGB15 image
func (i *RGB15) DrawImage(x, y int, draw *RGB15) {
	for dy := 0; dy < draw.height; dy++ {
		for dx := 0; dx < draw.width; dx++ {
			if y+dy >= i.height ||
				y+dy < 0 ||
				x+dx >= i.width ||
				x+dx < 0 {
				continue
			}

			c := draw.At(dx, dy)
			if c == transparent {
				continue
			}

			i.Set(x+dx, y+dy, c)
		}
	}
}

// At gets the memmap.PaletteValue from the image at location x,y
func (i *RGB15) At(x, y int) memmap.PaletteValue {
	return i.colors[y*i.width+x]
}

// ColAt gets the memmap.PaletteValue from the image at location x,y and converts that into
// a valid color.RGBA value
func (i *RGB15) ColAt(x, y int) color.RGBA {
	c := i.At(x, y)
	if c == transparent {
		return color.RGBA{}
	}

	b := int(c&0b01111100_00000000) >> 0xA
	g := int(c&0b00000011_11100000) >> 0x5
	r := int(c & 0b00000000_00011111)

	ret := color.RGBA{
		R: uint8(float64(r) * 8.2258),
		G: uint8(float64(g) * 8.2258),
		B: uint8(float64(b) * 8.2258),
		A: 255,
	}

	return ret
}

// Set sets the color of the pixel at x,y in the image to the given value
func (i *RGB15) Set(x, y int, c memmap.PaletteValue) {
	if y >= i.height ||
		y < 0 ||
		x >= i.width ||
		x < 0 {
		return
	}

	i.colors[y*i.width+x] = c
}

// v2 is a simple vector 2
type v2 struct {
	X int
	Y int
}

// tileData combines all the data related to a tile into a single structure
type tileData struct {
	mapOffset int
	gfxOffset int
	hflip     bool
	vflip     bool
}

// Background contains all the data for a GBA ppu background
type Background struct {
	controll    *memmap.BGControll
	bgHOffset   *uint16
	bgVOffset   *uint16
	enableCheck memmap.DisplayControll

	Enabled       bool
	Pos           v2
	Size          v2
	Priority      int
	Image         *RGB15
	SkipGFXUpdate bool
}

// update updates the backgrounds individual fields and gfx data
func (b *Background) update(palDirty bool) {
	b.Enabled = *display.Controll&b.enableCheck > 0
	if !b.Enabled {
		return
	}

	switch *b.controll & display.BGSizeMask {
	case display.BGSizeLarge:
		b.Size = v2{X: 2, Y: 2}
	case display.BGSizeTall:
		b.Size = v2{X: 1, Y: 2}
	case display.BGSizeWide:
		b.Size = v2{X: 2, Y: 1}
	case display.BGSizeSmall:
		b.Size = v2{X: 1, Y: 1}
	}

	b.Pos = v2{
		X: int(*b.bgHOffset),
		Y: int(*b.bgVOffset),
	}
	b.Priority = int(*b.controll & display.PriorityMask)

	if b.SkipGFXUpdate && !palDirty {
		return
	}

	screenBlock := (*b.controll & display.SBBMask) >> display.SBBShift
	vramOffset := int(screenBlock * memmap.ScreenBlockOffset)
	tileMap := memmap.VRAM[vramOffset : vramOffset+memmap.ScreenBlockOffset*b.Size.X*b.Size.Y]

	charBlock := (*b.controll & display.CBBMask) >> display.CBBShift
	gfxData := memmap.VRAM[charBlock*memmap.CharBlockOffset : (charBlock+1)*memmap.CharBlockOffset]

	for i := range tileMap {
		// palette data can differ across bg tiles so we need to calculate it inside the loop
		palette := int(tileMap[i]&0xF000) >> 0xc
		palOffset := memmap.PaletteOffset * palette
		palData := memmap.Palette[palOffset : palOffset+16]
		b.setTile(
			gfxData,
			palData,
			tileData{
				mapOffset: i,
				hflip:     (tileMap[i] & 0x0400) > 0,
				vflip:     (tileMap[i] & 0x0800) > 0,
				gfxOffset: int(tileMap[i] & 0x01FF),
			},
		)
	}
}

// setTile draws the tiles pixels onto the background image
func (b *Background) setTile(gfxData []memmap.VRAMValue, palData []memmap.PaletteValue, data tileData) {
	var indexes [16 * 4]int
	for i := 0; i < 16; i++ {
		quartet := getIndexQuartet(i, gfxData[data.gfxOffset*16:])
		indexes[i*4] = quartet[0]
		indexes[i*4+1] = quartet[1]
		indexes[i*4+2] = quartet[2]
		indexes[i*4+3] = quartet[3]
	}

	view := b.getTileView(data.mapOffset)
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			var px, py int
			switch {
			case data.hflip && data.vflip:
				px = (7 - x) + view.Bounds().Min.X
				py = (7 - y) + view.Bounds().Min.Y
			case data.hflip:
				px = (7 - x) + view.Bounds().Min.X
				py = y + view.Bounds().Min.Y
			case data.vflip:
				px = x + view.Bounds().Min.X
				py = (7 - y) + view.Bounds().Min.Y
			default:
				px = x + view.Bounds().Min.X
				py = y + view.Bounds().Min.Y
			}

			// TODO: instead of setting each pixel it would be better to copy data directly between images
			i := indexes[y*8+x]
			if i == 0 {
				b.Image.Set(px, py, transparent)
			} else {
				b.Image.Set(px, py, palData[i])
			}
		}
	}
}

// getTileView gets a specific rectangle in the background
func (b *Background) getTileView(tile int) image.Rectangle {
	screen := tile / 1024
	screenX := (screen % b.Size.X) * 256
	screenY := (screen / b.Size.X) * 256

	x := screenX + (((tile - screen*1024) % 32) * 8)
	y := screenY + (((tile - screen*1024) / 32) * 8)
	return image.Rect(x, y, x+8, y+8)
}

// Sprite is a PPU sprite
type Sprite struct {
	attrs *sprite.Attrs

	Enabled  bool
	Size     v2
	Pos      v2
	HFlip    bool
	VFlip    bool
	Priority int
	Image    *ebiten.Image
}

// update updates s sprites fields and gfx
func (s *Sprite) update() {
	s.Enabled = (s.attrs.Attr0 & sprite.SpriteModeMask) != sprite.Hide
	// TODO: clear out the image here (base it on the last size to improve performance)
	if !s.Enabled {
		return
	}

	s.Pos = v2{
		X: int(s.attrs.Attr1 & sprite.XMask),
		Y: int(s.attrs.Attr0 & sprite.YMask),
	}
	s.Size = s.sizeAsV2()
	// TODO: skip anything off screen (use the % calculation se we don't kill sprite wrapping)

	s.Priority = int(s.attrs.Attr2&sprite.PriorityMask) >> sprite.PriorityShift

	s.VFlip = (s.attrs.Attr1 & sprite.VMirriorMask) > 0
	s.HFlip = (s.attrs.Attr1 & sprite.HMirriorMask) > 0

	tile := (s.attrs.Attr2 & sprite.IndexMask)
	vramOffset := memmap.CharBlockOffset*4 + tile*16
	gfxData := memmap.VRAM[vramOffset:]

	pal := (s.attrs.Attr2 & sprite.PalMask) >> sprite.PalShift
	palOffset := memmap.PaletteOffset * (pal + 16)
	palData := memmap.Palette[palOffset : palOffset+16]

	s.updateImage(gfxData, palData)
}

// updateImage updates the sprites image data
func (s *Sprite) updateImage(gfxData []memmap.VRAMValue, palData []memmap.PaletteValue) {
	var indexes []int
	for i := 0; i < (s.Size.X/4)*s.Size.Y; i++ {
		quartet := getIndexQuartet(i, gfxData)
		indexes = append(indexes, quartet[:]...)
	}

	// fill all sprites with a transparent color to start because sprites will randomly shift around in memory
	// due to the fact that the underlying implementation uses a map
	s.Image.Fill(color.RGBA{})

	for i := 0; i < (s.Size.X*s.Size.Y)/64; i++ {
		for y := 0; y < 8; y++ {
			for x := 0; x < 8; x++ {
				color := palColorToRGBA(palData, indexes[i*64+y*8+x])
				tileX := (i * 8) % s.Size.X
				tileY := ((i * 8) / s.Size.X) * 8
				s.Image.Set(x+tileX, y+tileY, color)
			}
		}

	}
}

// sizeAsV2 converts the sprites size to a v2
func (s *Sprite) sizeAsV2() v2 {
	ref := [...]v2{
		{X: 8, Y: 8},
		{X: 16, Y: 8},
		{X: 8, Y: 16},

		{X: 16, Y: 16},
		{X: 32, Y: 8},
		{X: 8, Y: 32},

		{X: 32, Y: 32},
		{X: 32, Y: 16},
		{X: 16, Y: 32},

		{X: 64, Y: 64},
		{X: 64, Y: 32},
		{X: 32, Y: 64},
	}

	shape := int(s.attrs.Attr0&sprite.ShapeMask) >> 0xE
	size := int(s.attrs.Attr1&sprite.SizeMask) >> 0xE

	delta := ref[size*3+shape]

	return delta
}

// PPU is a emulated GBA PPU. It only emulates the pieces of the PPU used by flappy boot
type PPU struct {
	Sprites     [128]Sprite
	Backgrounds [4]Background
	lastPal     [512]memmap.PaletteValue
	palDirty    bool
	backBuffer  *RGB15
	Screen      *image.RGBA
}

// New creates a new PPU struct
func New() *PPU {
	ppu := &PPU{
		Backgrounds: [4]Background{
			{
				controll:    display.BG0Controll,
				bgHOffset:   display.BG0HOffset,
				bgVOffset:   display.BG0VOffset,
				enableCheck: display.BG0,

				//Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
				Image: NewRGB15(512, 512),
			},
			{
				controll:    display.BG1Controll,
				bgHOffset:   display.BG1HOffset,
				bgVOffset:   display.BG1VOffset,
				enableCheck: display.BG1,

				// Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
				Image: NewRGB15(512, 512),
			},
			{
				controll:    display.BG2Controll,
				bgHOffset:   display.BG2HOffset,
				bgVOffset:   display.BG2VOffset,
				enableCheck: display.BG2,

				// Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
				Image: NewRGB15(512, 512),
			},
			{
				controll:    display.BG3Controll,
				bgHOffset:   display.BG3HOffset,
				bgVOffset:   display.BG3VOffset,
				enableCheck: display.BG3,

				// Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
				Image: NewRGB15(512, 512),
			},
		},
		backBuffer: NewRGB15(240, 160),
		Screen:     image.NewRGBA(image.Rect(0, 0, 240, 160)),
	}

	for i := range ppu.Sprites {
		ppu.Sprites[i].Image = ebiten.NewImage(64, 64)
		ppu.Sprites[i].attrs = &sprite.OAM[i]
	}

	return ppu
}

// Update updates the ppu resources
func (p *PPU) Update() {
	// update palette cache to prevent updating background graphics unessiarily
	for i := range memmap.Palette {
		if p.lastPal[i] != memmap.Palette[i] {
			p.palDirty = true
			break
		}
	}
	if p.palDirty {
		copy(p.lastPal[:], memmap.Palette)
	}

	for i := range p.Backgrounds {
		p.Backgrounds[i].update(p.palDirty)
	}

	for i := range p.Sprites[:60] {
		p.Sprites[i].update()
	}

	for i := 3; i >= 0; i-- {
		for _, bg := range p.Backgrounds {
			if bg.Priority != i {
				continue
			}
			if !bg.Enabled {
				continue
			}

			x := -(bg.Pos.X % (bg.Size.X * 256))
			y := -bg.Pos.Y
			p.backBuffer.DrawImage(x, y, bg.Image)
			// duplicate the BG for horizontal scrolling
			p.backBuffer.DrawImage(x+(bg.Size.X*256), y, bg.Image)
		}

		// TOOD: we should also draw the sprites to the back buffer here but my first attempt ended
		// up killing performance so I'll need to make another atempt in the future.
	}

	for y := 0; y < 160; y++ {
		for x := 0; x < 240; x++ {
			p.Screen.Set(x, y, p.backBuffer.ColAt(x, y))
		}
	}

	p.palDirty = false
}

// getIndexQuartet converts a VRAMValue into the 4 palette indexes
func getIndexQuartet(i int, gfxData []memmap.VRAMValue) [4]int {
	return [4]int{
		int(gfxData[i] & 0x000F),
		int(gfxData[i]&0x00F0) >> 0x4,
		int(gfxData[i]&0x0F00) >> 0x8,
		int(gfxData[i]&0xF000) >> 0xC,
	}
}

// palColorToRGBA converts a palette's color into an RGBA color
func palColorToRGBA(palette []memmap.PaletteValue, index int) color.RGBA {
	if index == 0 {
		return color.RGBA{}
	}
	c := palette[index]

	b := (c & 0b01111100_00000000) >> 0xA
	g := (c & 0b00000011_11100000) >> 0x5
	r := c & 0b00000000_00011111

	ret := color.RGBA{
		R: uint8(float64(r) * 8.2258),
		G: uint8(float64(g) * 8.2258),
		B: uint8(float64(b) * 8.2258),
		A: 255,
	}

	return ret
}
