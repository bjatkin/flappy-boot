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
	player    *actor.Player
	score     *score.Counter
	highScore *score.Counter

	scoreBanner *game.MetaSprite
	bestBanner  *game.MetaSprite
	menu        *menu

	gravity   math.Fix8
	deathJump math.Fix8

	Restart, Quit bool
}

func NewScene(e *game.Engine, sky, clouds *game.Background, player *actor.Player, roundScore, highScore *score.Counter) (*Scene, error) {
	scoreBanner, err := e.NewMetaSprite(
		[]math.V2{{X: 0, Y: 0}, {X: math.FixOne * 32, Y: 0}},
		[]int{24, 0},
		assets.BannersTileSet,
	)
	if err != nil {
		return nil, err
	}

	bestBanner, err := e.NewMetaSprite(
		[]math.V2{{X: 0, Y: 0}, {X: math.FixOne * 32, Y: 0}},
		[]int{16, 8},
		assets.BannersTileSet,
	)
	if err != nil {
		return nil, err
	}

	menu, err := newMenu(math.FixOne*87, math.FixOne*102, e)
	if err != nil {
		return nil, err
	}

	return &Scene{
		sky:       sky,
		clouds:    clouds,
		player:    player,
		score:     roundScore,
		highScore: highScore,

		scoreBanner: scoreBanner,
		bestBanner:  bestBanner,
		menu:        menu,

		gravity:   math.FixQuarter,
		deathJump: -math.FixOne * 6,
	}, nil
}

func (s *Scene) Init(e *game.Engine) error {
	s.player.Sprite.HFlip = true
	s.player.Update(s.gravity, s.deathJump)
	s.menu.Reset(math.FixOne*87, math.FixOne*102)

	s.scoreBanner.Set(math.FixOne*87, math.FixOne*8)
	err := s.scoreBanner.Add()
	if err != nil {
		return err
	}

	s.bestBanner.Set(math.FixOne*87, math.FixOne*48)
	err = s.bestBanner.Add()
	if err != nil {
		return err
	}

	s.highScore.X = 97
	s.highScore.Y = 70

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
	if s.menu.selectCountDown <= 0 {
		s.Restart = s.menu.restart
		s.Quit = s.menu.quit
	}
	return nil
}

func (s *Scene) Hide() {
	s.menu.Hide()
	s.bestBanner.Set(math.FixOne*240, 0)
	s.scoreBanner.Set(math.FixOne*240, 0)
	s.highScore.X = 240
	s.highScore.Draw()
}

// menu is a simple game over menu
type menu struct {
	x, y            math.Fix8
	arrow           *game.Sprite
	bg              *game.Background
	restart, quit   bool
	selectCountDown int
	selectStart     int
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
		x:           x,
		y:           y,
		arrow:       arrow,
		bg:          bg,
		selectStart: 30,
	}, nil
}

// Update updates the menu state each frame
func (m *menu) Update() {
	if m.selectStart > 0 {
		m.selectStart--
		m.arrow.Update()
		return
	}

	if m.restart || m.quit {
		m.selectCountDown--
		m.arrow.Update()
		return
	}

	if key.JustPressed(key.Down) {
		m.arrow.Y = m.y + math.FixOne*12
	}
	if key.JustPressed(key.Up) {
		m.arrow.Y = m.y
	}
	if key.JustPressed(key.A) && m.arrow.Y == m.y {
		m.restart = true
	}
	if key.JustPressed(key.A) && m.arrow.Y > m.y {
		m.quit = true
	}

	if m.restart || m.quit {
		m.arrow.SetAnimation(
			game.Frame{Index: 2, Len: 7},
			game.Frame{Index: 3, Len: 7},
		)
		m.selectCountDown = 60
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

func (m *menu) Hide() {
	m.bg.VScroll = math.FixOne * 160
	m.arrow.X = math.FixOne * 240
}

func (m *menu) Reset(x, y math.Fix8) {
	m.arrow.TileIndex = 2
	m.arrow.SetAnimation(
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

	m.x = x
	m.y = y
	m.selectStart = 30
	m.bg.VScroll = 0
	m.arrow.X = x
	m.restart = false
	m.quit = false
}
