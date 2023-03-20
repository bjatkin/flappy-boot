package game

import (
	"github.com/bjatkin/flappy_boot/internal/alloc"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

type Runable interface {
	Init(*Engine) error
	Update(*Engine, int) error
	Next() (Runable, bool)
}

// Engine is the core game engine
type Engine struct {
	// activeSprites are the sprites tha need to be drawn each frame
	activeSprites map[*Sprite]struct{}

	// activeBackgrounds are the backgrounds that need to be drawn each frame
	activeBackgrounds [4]*Background

	// Allocators
	bgPalAlloc   *alloc.Pal
	sprPalAlloc  *alloc.Pal
	bgTileAlloc  *alloc.VRAM
	sprTileAlloc *alloc.VRAM
	mapAlloc     *alloc.VRAM
}

// NewEngine creates a new instances of a game engine
func NewEngine() *Engine {
	memmap.SetReg(display.Controll, display.Sprite1D|display.ForceBlank)
	return &Engine{
		activeSprites: make(map[*Sprite]struct{}, 128),
		bgPalAlloc:    alloc.NewPal(memmap.Palette[:256]),
		sprPalAlloc:   alloc.NewPal(memmap.Palette[256:]),

		// the first tile is left transparent and can be shared by all tile maps
		bgTileAlloc:  alloc.NewVRAM(memmap.VRAM[memmap.TileOffset4:memmap.CharBlockOffset*2], 16),
		sprTileAlloc: alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*4:], 16),
		mapAlloc:     alloc.NewVRAM(memmap.VRAM[memmap.CharBlockOffset*2:], memmap.HalfKByte*2),
	}
}

func (e *Engine) Run(run Runable) error {
	// enable sprites
	memmap.SetReg(display.Controll, *display.Controll|display.Sprites)
	sprite.Reset()

	for {
		err := run.Init(e)
		if err != nil {
			exit(err)
		}

		var frame int
		for {
			key.KeyPoll()
			err := run.Update(e, frame)
			if err != nil {
				exit(err)
			}

			frame++

			vSyncWait()

			// copy active sprite data into OAM memory
			e.drawSprites()

			// copy active background data into the background registers
			e.drawBackgrounds()

			if next, ok := run.Next(); ok {
				sprite.Reset()
				run = next
				break
			}
		}
	}
}

func (e *Engine) drawSprites() {
	var i int
	for s := range e.activeSprites {
		hw_sprite.OAM[i] = *s.attrs()
		i++
	}
}

func (e *Engine) drawBackgrounds() {
	// TODO: only update the backgrouds if something has changed?
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			continue
		}

		controll := e.activeBackgrounds[i].controll()
		switch i {
		case 0:
			memmap.SetReg(display.BG0Controll, controll)
			memmap.SetReg(display.Controll, *display.Controll|display.BG0)
			memmap.SetReg(display.BG0HOffset, e.activeBackgrounds[0].hScroll)
			memmap.SetReg(display.BG0VOffset, e.activeBackgrounds[0].vScroll)
		case 1:
			memmap.SetReg(display.BG1Controll, controll)
			memmap.SetReg(display.Controll, *display.Controll|display.BG1)
			memmap.SetReg(display.BG1HOffset, e.activeBackgrounds[1].hScroll)
			memmap.SetReg(display.BG1VOffset, e.activeBackgrounds[1].vScroll)
		case 2:
			memmap.SetReg(display.BG2Controll, controll)
			memmap.SetReg(display.Controll, *display.Controll|display.BG2)
			memmap.SetReg(display.BG2HOffset, e.activeBackgrounds[2].hScroll)
			memmap.SetReg(display.BG2VOffset, e.activeBackgrounds[2].vScroll)
		case 3:
			memmap.SetReg(display.BG3Controll, controll)
			memmap.SetReg(display.Controll, *display.Controll|display.BG3)
			memmap.SetReg(display.BG3HOffset, e.activeBackgrounds[3].hScroll)
			memmap.SetReg(display.BG3VOffset, e.activeBackgrounds[3].vScroll)
		}
	}
}

// // PalAllocation is an allocation for a palette in GBA palette memory
// type PalAllocation struct {
// 	t      PalAllocationType
// 	Pal    memmap.VRAMValue
// 	Memory []memmap.PaletteValue
// }

// // PalAllocationType indicates if a palette is a background or sprite palette
// type PalAllocationType int

// const (
// 	// BGPalette is a background palette in GBA palette memory
// 	BGPalette PalAllocationType = iota

// 	// SpritePalette is a sprite palette in GBA palette memory
// 	SpritePalette
// )

// // AllocPal allocates a palette in GBA palette memory
// func (e *Engine) AllocPal(allocType PalAllocationType) (*PalAllocation, error) {
// 	var start int
// 	if allocType == SpritePalette {
// 		start = 8
// 	}

// 	for i := 0; i < 8; i++ {
// 		if !e.palettes[start+i] {
// 			return &PalAllocation{
// 				t:      allocType,
// 				Pal:    memmap.VRAMValue(i) << display.PaletteShift,
// 				Memory: memmap.Palette[memmap.PaletteOffset*(i+16) : memmap.PaletteOffset*(i+17)],
// 			}, nil
// 		}
// 	}
// 	return nil, ErrOOM
// }

// // FreePal frees an allocated palette from the GBAs palette memory
// func (e *Engine) FreePal(alloc *PalAllocation) {
// 	i := alloc.Pal >> display.PaletteShift
// 	if alloc.t == BGPalette {
// 		i += 8
// 	}

// 	e.palettes[i] = false
// }

// // TileAllocation is an allocation for tile data in the GBAs VRAM
// type TileAllocation struct {
// 	Memory []memmap.VRAMValue
// 	Index  memmap.VRAMValue
// }

// // AllocBGTiles allocates an array of background tiles in GBA VRAM
// func (e *Engine) AllocBGTiles(tiles int) (*TileAllocation, error) {
// 	i, err := e.bgTiles.Alloc(tiles)
// 	if err != nil {
// 		return nil, err
// 	}

// 	start := i * 16
// 	end := start + (tiles * 16)
// 	return &TileAllocation{
// 		Memory: memmap.VRAM[start:end],
// 		Index:  memmap.VRAMValue(i),
// 	}, nil
// }

// // FreeBGTiles frees an allocated array of background tiles in GBA VRAM
// func (e *Engine) FreeBGTiles(alloc *TileAllocation) {
// 	e.bgTiles.Free(int(alloc.Index))
// }

// // TileMapAllocation is an allocation for map data in the GBAs VRAM
// type TileMapAllocation struct {
// 	Memory         []memmap.VRAMValue
// 	CharacterBlock memmap.BGControll
// 	ScreenBlock    memmap.BGControll
// }

// // AllocTileMap allocates an array of screen blocks in the GBAs VRAM
// func (e *Engine) AllocTileMap(screens int) (*TileMapAllocation, error) {
// 	i, err := e.maps.Alloc(screens)
// 	if err != nil {
// 		return nil, err
// 	}

// 	characterBlock := memmap.BGControll(2) << display.CBBShift
// 	if i > 8 {
// 		characterBlock = 3 << display.CBBShift
// 		i -= 8
// 	}

// 	start := (memmap.CharBlockOffset * 2)
// 	end := start + (screens * memmap.ScreenBlockOffset)
// 	return &TileMapAllocation{
// 		Memory:         memmap.VRAM[start:end],
// 		CharacterBlock: characterBlock,
// 		ScreenBlock:    memmap.BGControll(i) << display.SBBShift,
// 	}, nil
// }

// // FreeTileMap frees an allocated array of screen blocks in GBA VRAM
// func (e *Engine) FreeTileMap(alloc *TileMapAllocation) {
// 	i := alloc.ScreenBlock >> display.SBBShift
// 	if alloc.CharacterBlock == 3<<display.CBBShift {
// 		i += 8
// 	}

// 	e.maps.Free(int(i))
// }

// // SprAllocation is an allocated array of sprite tiles in GBA VRAM
// type SprAllocation struct {
// 	Memory []memmap.VRAMValue
// 	Offset hw_sprite.Attr2
// }

// // AllocSprs allocates an array of sprite tiles in GBA VRAM
// func (e *Engine) AllocSprs(sprites int) (SprAllocation, error) {
// 	i, err := e.sprTiles.Alloc(sprites)
// 	if err != nil {
// 		return SprAllocation{}, ErrOOM
// 	}

// 	start := 4*memmap.CharBlockOffset + i*16
// 	end := start + sprites*16
// 	return SprAllocation{
// 		Memory: memmap.VRAM[start:end],
// 		Offset: hw_sprite.Attr2(i),
// 	}, nil
// }

// // FreeSprs frees an allocated array of sprite tiles in GBA VRAM
// func (e *Engine) FreeSprs(alloc SprAllocation) {
// 	e.sprTiles.Free(int(alloc.Offset))
// }

// func (e *Engine) loadBGPal(data []memmap.PaletteValue) int {
// 	palID := e.bgPalPtr
// 	for i := range data {
// 		memmap.Palette[i+(memmap.PaletteOffset*palID)] = data[i]
// 	}

// 	e.bgPalPtr += 1

// 	return palID
// }

// func (e *Engine) loadSprPal(data []memmap.PaletteValue) int {
// 	palID := e.spritePalPtr
// 	for i := range data {
// 		memmap.Palette[i+memmap.PaletteOffset*palID] = data[i]
// 	}
// 	e.spritePalPtr++

// 	return palID - 16
// }

// func (e *Engine) loadCB(data []memmap.VRAMValue) (int, int) {
// 	tileOffset := e.bgPtr
// 	for i := range data {
// 		memmap.VRAM[i+tileOffset] = data[i]
// 	}

// 	e.bgPtr += len(data)
// 	return 0, tileOffset / 16
// }

// func (e *Engine) loadSB(data []memmap.VRAMValue, offset, palID int) int {
// 	screenID := e.screenBlockPtr
// 	for i := range data {
// 		memmap.VRAM[i+memmap.ScreenBlockOffset*screenID] = (data[i] + memmap.VRAMValue(offset)) | memmap.VRAMValue(palID)<<0x0C
// 	}

// 	e.screenBlockPtr += len(data) / memmap.ScreenBlockOffset
// 	return screenID
// }

// func (e *Engine) updateSB(screenBlock, offset int, tile memmap.VRAMValue) {
// 	memmap.VRAM[offset+memmap.ScreenBlockOffset*screenBlock] = tile
// }

// func (e *Engine) loadSprite(data []memmap.VRAMValue) int {
// 	tileOffset := e.spritePtr
// 	for i := range data {
// 		memmap.VRAM[i+tileOffset] = data[i]
// 	}

// 	e.spritePtr += len(data)
// 	return tileOffset - memmap.CharBlockOffset*4
// }

// addBackground adds a new background to the list of active backgrounds.
func (e *Engine) addBackground(bg *Background) error {
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == nil {
			e.activeBackgrounds[i] = bg
			return nil
		}
	}

	return alloc.ErrOOM
}

// removeBackground removes the current background from the list of active backgrounds. It will not unload the
// background from memory so you must do that yourself if the background is no longer needed
func (e *Engine) removeBackground(bg *Background) {
	for i := range e.activeBackgrounds {
		if e.activeBackgrounds[i] == bg {
			e.activeBackgrounds[i] = nil
			return
		}
	}
}

// addSprite adds a new sprite to the list of active sprites.
func (e *Engine) addSprite(sprite *Sprite) {
	e.activeSprites[sprite] = struct{}{}
}

// removeSprite removes a sprite from the list of active sprites. It will not unload the sprites assets
// from memory so you must do that yourself if the sprite is no longer needed
func (e *Engine) removeSprite(sprite *Sprite) {
	delete(e.activeSprites, sprite)
}

// NewSprite returns a new Sprite
func (e *Engine) NewSprite(tileSet *assets.TileSet) *Sprite {
	return &Sprite{
		engine:  e,
		tileSet: tileSet,
		size:    hw_sprite.Medium,
		shape:   hw_sprite.Square,
	}
}

/*
// TileSet is a set of 8x8 tiles that can be loaded into VRAM for use by a background or sprite
type TileSet struct {
	// loaded is true if this tileset has been loaded into VRAM
	loaded bool

	// the number of 8x8 tiles in this tileset
	Count uint32

	// tiles is the pixel data for the tile set
	Tiles     []memmap.VRAMValue
	TileIndex int

	charBaseBlock int
	Sprite        bool

	// palette is the palette data for the tileset
	Palette      Palette
	PaletteIndex int
}

func (ts *TileSet) Load(e *Engine) error {
	if ts.loaded {
		return nil
	}

	if ts.Sprite {
		ts.PaletteIndex = e.loadSprPal(ts.Palette)

		ts.TileIndex = e.loadSprite(ts.Tiles)
	} else {
		ts.PaletteIndex = e.loadBGPal(ts.Palette)

		ts.charBaseBlock, ts.TileIndex = e.loadCB(ts.Tiles)
	}

	// // TODO: get rid of these magic numbers
	// if palBase > 32 {
	// 	return fmt.Errorf("OOM: invalid palette base %d", palBase)
	// }

	// // TODO: get rid of these magic numbers
	// if tileBase > 512*6 {
	// 	return fmt.Errorf("OOM: invalid tile base %d", tileBase)
	// }
	ts.loaded = true

	return nil
}

// TileMap is tilemap data for a background
type TileMap struct {
	// loaded is true if this tilemap has been loaded into VRAM
	loaded bool

	screenBaseBlock int
	ScreenSize      memmap.BGControll
	Data            []memmap.VRAMValue
	TileSet         *TileSet
}

func (t *TileMap) Load(e *Engine) {
	if t.loaded {
		return
	}

	// make sure the tile set is loaded into memory
	t.TileSet.Load(e)

	t.screenBaseBlock = e.loadSB(t.Data, t.TileSet.TileIndex, t.TileSet.PaletteIndex)
	t.loaded = true
}

// Paletet is a 16 color palette
type Palette []memmap.PaletteValue
*/
