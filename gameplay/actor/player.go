package actor

import (
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/math"
)

// Player is a struct representing a player
type Player struct {
	Sprite *game.Sprite
	dy     math.Fix8
	maxDy  math.Fix8

	dead    bool
	started bool
}

// NewPlayer creates a new player struct
func NewPlayer(pos math.V2, sprite *game.Sprite) *Player {
	sprite.TileIndex = 16
	sprite.PlayAnimation(glideAni)
	sprite.Pos = pos

	return &Player{
		Sprite: sprite,
		maxDy:  math.FixOne * 8,
	}
}

// Start indicates the the game has started and the player should start applying physics
func (p *Player) Start() {
	p.started = true
}

// Dead sets the player state to dead
func (p *Player) Dead() {
	p.dead = true
}

// Init resets all the players properties to be the same as they were on creation. It also move the sprite to the specified location
func (p *Player) Init(pos math.V2) {
	p.dy = 0
	p.started = false
	p.dead = false
	p.Sprite.Pos = pos
	p.Sprite.HFlip = false
	p.Sprite.TileIndex = 16
	p.Sprite.PlayAnimation(glideAni)
}

// Rect returns the hitbox of the player as a math.Rect
func (p *Player) Rect() math.Rect {
	return math.Rect{
		X1: p.Sprite.Pos.X.Int() + 12,
		Y1: p.Sprite.Pos.Y.Int() + 2,
		X2: p.Sprite.Pos.X.Int() + 22,
		Y2: p.Sprite.Pos.Y.Int() + 12,
	}
}

// Show whos the player sprite
func (p *Player) Show() error {
	err := p.Sprite.Show()
	if err != nil {
		return err
	}

	return nil
}

// Hide hides the player sprite
func (p *Player) Hide() {
	p.Sprite.Hide()
}

var (
	jumpAni = []game.Frame{
		{Index: 16, Len: 3},
		{Index: 32, Len: 4},
		{Index: 0, Len: 7},
		{Index: 8, Len: 8},
		{Index: 24, Len: 2},

		{Index: 16, Len: 40},
		{Index: 24, Len: 40, Offset: math.V2{X: 0, Y: math.FixOne}},

		{Index: 16, Len: 40},
		{Index: 24, Len: 40, Offset: math.V2{X: 0, Y: math.FixOne}},

		{Index: 16, Len: 40},
		{Index: 24, Len: 40, Offset: math.V2{X: 0, Y: math.FixOne}},
	}

	glideAni = []game.Frame{
		{Index: 16, Len: 40},
		{Index: 24, Len: 40, Offset: math.V2{X: 0, Y: math.FixOne}},
	}
)

// Update updates the players physics and interal properites
func (p *Player) Update(gravity, jump math.Fix8) {
	p.Sprite.Update()

	if !p.started {
		// don't update physics if the game has not started yet
		return
	}

	p.dy += gravity
	if p.dy > p.maxDy {
		p.dy = p.maxDy
	}

	if jump != 0 {
		p.Sprite.PlayAnimation(jumpAni)
		p.dy = jump
	}

	if p.dead {
		p.Sprite.StopAnimation()
		p.Sprite.HFlip = true
		p.Sprite.TileIndex = 0
	}

	p.Sprite.Pos.Y += p.dy
	if p.Sprite.Pos.Y > math.FixOne*200 {
		p.Sprite.Pos.Y = math.FixOne * 200
	}

	if p.Sprite.Pos.Y < -math.FixOne*16 {
		p.Sprite.Pos.Y = -math.FixOne * 16
	}
}
