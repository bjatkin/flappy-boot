package gameover

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	sky       *game.Background
	clouds    *game.Background
	pillars   *game.Background
	player    *actor.Player
	score     *score.Counter
	highScore *score.Counter

	scoreBanner *game.MetaSprite
	bestBanner  *game.MetaSprite
	menu        *menu

	gravity   math.Fix8
	deathJump math.Fix8
}

func NewScene(sky, clouds, pillars *game.Background, player *actor.Player, roundScore *score.Counter) *Scene {
	if roundScore.Score() > score.Best {
		score.Best = roundScore.Score()
	}

	return &Scene{
		sky:     sky,
		clouds:  clouds,
		pillars: pillars,
		player:  player,
		score:   roundScore,

		gravity:   math.FixQuarter,
		deathJump: -math.FixOne * 6,
	}
}

func (s *Scene) Init(e *game.Engine) error {
	s.player.Sprite.HFlip = true
	s.player.Update(s.gravity, s.deathJump)

	// TODO: should this lerp in?
	// shift the score down so it can sit uner the score banner
	s.score.Y = 28

	var err error
	s.scoreBanner, err = e.NewMetaSprite(
		[]math.V2{{X: 0, Y: 0}, {X: math.FixOne * 32, Y: 0}},
		[]int{0, 16},
		assets.BannersTileSet,
	)
	if err != nil {
		return err
	}

	s.scoreBanner.Set(math.FixOne*87, math.FixOne*8)
	err = s.scoreBanner.Add()
	if err != nil {
		return err
	}

	s.bestBanner, err = e.NewMetaSprite(
		[]math.V2{{X: 0, Y: 0}, {X: math.FixOne * 32, Y: 0}},
		[]int{8, 24},
		assets.BannersTileSet,
	)
	if err != nil {
		return err
	}

	s.bestBanner.Set(math.FixOne*87, math.FixOne*48)
	err = s.bestBanner.Add()
	if err != nil {
		return err
	}

	s.highScore = score.NewCounter(97, 70, e)
	s.highScore.Set(score.Best)

	s.menu, err = newMenu(math.FixOne*87, math.FixOne*102, e)
	if err != nil {
		return err
	}

	err = s.menu.Add()
	if err != nil {
		return err
	}

	return nil
}

func (s *Scene) Update(e *game.Engine, frame int) error {
	s.player.Update(s.gravity, 0)
	s.score.Draw()
	s.highScore.Draw()
	s.menu.Update()
	return nil
}

func (s *Scene) Next() (game.Runable, bool) {
	return nil, false
}

// menu is a simple game over menu
type menu struct {
	x, y  math.Fix8
	arrow *game.Sprite
	bg    *game.Background
}

// newMenu creates a new game over menu
func newMenu(x, y math.Fix8, e *game.Engine) (*menu, error) {
	arrow := e.NewSprite(assets.SelectTileSet)
	arrow.X = x
	arrow.Y = y
	arrow.TileIndex = 2
	arrow.SetAnimation(
		game.Frame{Index: 2, Len: 30},
		game.Frame{Index: 1, Len: 10},
		game.Frame{Index: 0, Len: 10},
		game.Frame{Index: 0, HFlip: true, Len: 10},
		game.Frame{Index: 1, HFlip: true, Len: 10},
		game.Frame{Index: 2, Len: 10},
		game.Frame{Index: 1, Len: 10},
		game.Frame{Index: 0, Len: 10},
		game.Frame{Index: 0, HFlip: true, Len: 10},
		game.Frame{Index: 1, HFlip: true, Len: 10},
	)

	bg := e.NewBackground(assets.BluebgTileMap, display.Priority0)

	return &menu{
		x:     x,
		y:     y,
		arrow: arrow,
		bg:    bg,
	}, nil
}

// Update updates the menu state each frame
func (m *menu) Update() {
	if key.JustPressed(key.Down) {
		m.arrow.Y = m.y + math.FixOne*12
	}
	if key.JustPressed(key.Up) {
		m.arrow.Y = m.y
	}

	m.arrow.Update()
}

// Add adds menu sprites and backgrounds into active engine memory
func (m *menu) Add() error {
	err := m.bg.Add()
	if err != nil {
		return err
	}

	err = m.arrow.Add()
	if err != nil {
		return err
	}

	return nil
}
