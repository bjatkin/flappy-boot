package ppu

import (
	"image"
	"image/color"

	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/hajimehoshi/ebiten/v2"
)

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
	Image         *image.RGBA
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
			color := palColorToRGBA(palData, indexes[y*8+x])

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
			b.Image.Set(px, py, color)
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

				Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
			},
			{
				controll:    display.BG1Controll,
				bgHOffset:   display.BG1HOffset,
				bgVOffset:   display.BG1VOffset,
				enableCheck: display.BG1,

				Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
			},
			{
				controll:    display.BG2Controll,
				bgHOffset:   display.BG2HOffset,
				bgVOffset:   display.BG2VOffset,
				enableCheck: display.BG2,

				Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
			},
			{
				controll:    display.BG3Controll,
				bgHOffset:   display.BG3HOffset,
				bgVOffset:   display.BG3VOffset,
				enableCheck: display.BG3,

				Image: image.NewRGBA(image.Rect(0, 0, 512, 512)),
			},
		},
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

// palCache makes converting palette values faster by caching past colors that have been converted
var palCache = make(map[memmap.PaletteValue]color.RGBA, 1024)

// palColorToRGBA converts a palette's color into an RGBA color
func palColorToRGBA(palette []memmap.PaletteValue, index int) color.RGBA {
	if index == 0 {
		return color.RGBA{}
	}
	c := palette[index]

	if color, ok := palCache[c]; ok {
		return color
	}

	b := (c & 0b01111100_00000000) >> 0xA
	g := (c & 0b00000011_11100000) >> 0x5
	r := c & 0b00000000_00011111

	ret := color.RGBA{
		// R: uint8(float64(r) * 8.2258),
		// G: uint8(float64(g) * 8.2258),
		// B: uint8(float64(b) * 8.2258),
		// multiplying by 8 is less accurate but way faster
		R: uint8(r * 8),
		G: uint8(g * 8),
		B: uint8(b * 8),
		A: 255,
	}

	palCache[c] = ret
	return ret
}
