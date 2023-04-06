package game

import (
	"errors"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type MetaSprite struct {
	// engine is a reference to the sprites parent engine
	engine *Engine

	pos     math.V2
	sprites []*Sprite
	offsets []math.V2
}

// NewSprite returns a new Sprite
func (e *Engine) NewMetaSprite(offset []math.V2, indexes []int, asset *assets.TileSet) (*MetaSprite, error) {
	if len(offset) != len(indexes) {
		return nil, errors.New("offsets and indexes do not match")
	}

	sprites := make([]*Sprite, len(offset))
	for i := range indexes {
		sprite := e.NewSprite(asset)
		sprite.TileIndex = indexes[i]
		sprites[i] = sprite
	}

	return &MetaSprite{
		engine:  e,
		offsets: offset,
		sprites: sprites,
	}, nil
}

// Set sets the x and y position of the meta sprite
func (s *MetaSprite) Set(pos math.V2) {
	s.pos = pos
	for i := range s.sprites {
		s.sprites[i].Pos = math.AddV2(s.pos, s.offsets[i])
	}
}

// Move moves the meta sprite by dx and dy
func (s *MetaSprite) Move(delta math.V2) {
	s.pos = math.AddV2(s.pos, delta)
	for i := range s.sprites {
		s.sprites[i].Pos = math.AddV2(s.pos, s.offsets[i])
	}
}

// Show adds the meta sprite's component sprites to the list of active sprites.
// if the sprites associated assets have not been loaded yet, Show will automatically attempt to load them.
// all active sprites are drawn every frame, if more than 128 sprites are active at a time all active
// sprites will be randomly flickered to ensure all sprites continue to be drawn
func (s *MetaSprite) Show() error {
	for i := range s.sprites {
		err := s.sprites[i].Show()
		if err != nil {
			return err
		}
	}
	return nil
}

// Hide removes the sprites from the list of active sprites.
// removing a sprites does not unload it's loaded assets from VRAM. To do that you must call Unload
func (s *MetaSprite) Hide() {
	for i := range s.sprites {
		s.sprites[i].Hide()
	}
}

// Load loads a sprites graphics data into memory
// if there is not enough free VRAM to accomodate the sprite an error will be returned
func (s *MetaSprite) Load() error {
	for i := range s.sprites {
		err := s.sprites[i].Load()
		if err != nil {
			return err
		}
	}
	return nil
}

// Unload removes a sprites graphics data from memory
func (s *MetaSprite) UnLoad() {
	for i := range s.sprites {
		s.sprites[i].UnLoad()
	}
}
