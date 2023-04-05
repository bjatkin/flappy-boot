package game

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	hw_sprite "github.com/bjatkin/flappy_boot/internal/hardware/sprite"
	"github.com/bjatkin/flappy_boot/internal/math"
)

// Frame is a single frame of sprite animation data
type Frame struct {
	Index  int
	HFlip  bool
	VFlip  bool
	Offset math.V2
	Len    int
}

// Sprite is a game engine sprite
type Sprite struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	Y         math.Fix8
	X         math.Fix8
	Offset    math.V2
	TileIndex int
	Hide      bool
	HFlip     bool
	VFlip     bool
	Priority  int
	size      hw_sprite.Attr1
	shape     hw_sprite.Attr0

	animation  []Frame
	aniFrame   int
	aniCounter int

	tileSet *assets.TileSet
	hwAttrs *hw_sprite.Attrs
}

// NewSprite returns a new Sprite
func (e *Engine) NewSprite(tileSet *assets.TileSet) *Sprite {
	return &Sprite{
		engine:  e,
		tileSet: tileSet,
		size:    tileSet.Size(),
		shape:   tileSet.Shape(),
		hwAttrs: &hw_sprite.Attrs{},
	}
}

func (s *Sprite) attrs() *hw_sprite.Attrs {
	var hideAttr hw_sprite.Attr0
	if s.Hide {
		hideAttr = hw_sprite.Hide
	}

	var vFlipAttr hw_sprite.Attr1
	if s.VFlip {
		vFlipAttr = hw_sprite.HMirrior
	}
	var hFlipAttr hw_sprite.Attr1
	if s.HFlip {
		hFlipAttr = hw_sprite.VMirrior
	}

	var priorityAttr hw_sprite.Attr2
	switch s.Priority {
	case 0:
		priorityAttr = hw_sprite.Priority0
	case 1:
		priorityAttr = hw_sprite.Priority1
	case 2:
		priorityAttr = hw_sprite.Priority2
	case 3:
		priorityAttr = hw_sprite.Priority3
	}

	x := s.X.Int() + s.Offset.X.Int()
	if x < 0 {
		x += 512
	}
	y := s.Y.Int() + s.Offset.Y.Int()
	if y < 0 {
		y += 256
	}
	s.hwAttrs.Attr0 = hw_sprite.Attr0(y%256) | s.shape | hideAttr
	s.hwAttrs.Attr1 = hw_sprite.Attr1(x%512) | vFlipAttr | hFlipAttr | s.size
	s.hwAttrs.Attr2 = (hw_sprite.Attr2(s.TileIndex) + s.tileSet.Offset()) |
		priorityAttr |
		s.tileSet.SprPalette()

	return s.hwAttrs
}

// PlayAnimation sets the animation data for the sprite
func (s *Sprite) PlayAnimation(frames []Frame) {
	s.aniCounter = 0
	s.aniFrame = 0
	s.animation = frames
}

// StopAnimation removes the animation data from the sprite
func (s *Sprite) StopAnimation() {
	s.animation = nil
}

// Add adds the sprite to the list of active sprites.
// if the sprites associated assets have not been loaded yet, Add will automatically attempt to load them.
// all active sprites are drawn every frame, if more than 128 sprites are active at a time all active
// sprites will be randomly flickered to ensure all sprites continue to be drawn
func (s *Sprite) Add() error {
	err := s.Load()
	if err != nil {
		return err
	}

	s.engine.addSprite(s)
	return nil
}

// Remove removes the sprites from the list of active sprites.
// removing a sprites does not unload it's loaded assets from VRAM. To do that you must call Unload
func (s *Sprite) Remove() {
	s.engine.removeSprite(s)
}

// Load loads a sprites graphics data into memory
// if there is not enough free VRAM to accomodate the sprite an error will be returned
func (s *Sprite) Load() error {
	return s.tileSet.Load(s.engine.sprTileAlloc, s.engine.sprPalAlloc)
}

// Unload removes a sprites graphics data from memory
func (s *Sprite) UnLoad() {
	s.tileSet.Free(s.engine.sprTileAlloc, s.engine.sprPalAlloc)
}

// Update updates the sprites graphics
func (s *Sprite) Update() {
	if s.animation == nil {
		return
	}

	s.aniCounter++
	if s.aniCounter < s.animation[s.aniFrame].Len {
		return
	}
	s.aniCounter = 0
	s.aniFrame++
	s.aniFrame %= len(s.animation)

	s.TileIndex = s.animation[s.aniFrame].Index
	s.HFlip = s.animation[s.aniFrame].HFlip
	s.VFlip = s.animation[s.aniFrame].VFlip
	s.Offset = s.animation[s.aniFrame].Offset
}
