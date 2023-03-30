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

	started bool
}

// NewPlayer creates a new player struct
func NewPlayer(x, y math.Fix8, sprite *game.Sprite) *Player {
	p := &Player{
		Sprite: sprite,
		maxDy:  math.FixOne * 8,
	}

	p.Sprite.X = x
	p.Sprite.Y = y
	return p
}

func (p *Player) Start() {
	p.started = true
}

func (p *Player) Rect() math.Rect {
	return math.Rect{
		X1: p.Sprite.X.Int() + 2,
		Y1: p.Sprite.Y.Int() + 2,
		X2: p.Sprite.X.Int() + 12,
		Y2: p.Sprite.Y.Int() + 12,
	}
}

func (p *Player) Show() error {
	err := p.Sprite.Add()
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) Hide() {
	p.Sprite.Remove()
}

func (p *Player) Update(gravity, jump math.Fix8) {
	if !p.started {
		// don't update physics if the game has not started yet
		return
	}

	p.dy += gravity
	if p.dy > p.maxDy {
		p.dy = p.maxDy
	}

	if jump != 0 {
		p.dy = jump
	}

	p.Sprite.Y += p.dy
	if p.Sprite.Y > math.FixOne*200 {
		p.Sprite.Y = math.FixOne * 200
	}

	if p.Sprite.Y < -math.FixOne*16 {
		p.Sprite.Y = -math.FixOne * 16
	}
}
