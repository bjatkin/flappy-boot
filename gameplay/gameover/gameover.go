package gameover

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/pillar"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/lut"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	sky       *game.Background
	clouds    *game.Background
	player    *actor.Player
	score     *score.Counter
	highScore *score.Counter
	pillars   *pillar.BG

	scoreBanner *game.MetaSprite
	bestBanner  *game.MetaSprite
	menu        *menu

	gravity   math.Fix8
	deathJump math.Fix8
	t         math.Fix8
	palFade   math.Fix8

	Restart, Quit bool
}

func NewScene(e *game.Engine, sky, clouds *game.Background, pillars *pillar.BG, player *actor.Player, roundScore, highScore *score.Counter) (*Scene, error) {
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
		pillars:   pillars,

		scoreBanner: scoreBanner,
		bestBanner:  bestBanner,
		menu:        menu,

		gravity:   math.FixQuarter,
		deathJump: -math.FixOne * 6,
	}, nil
}

func (s *Scene) Init(e *game.Engine) error {
	s.t = 0
	s.palFade = 0

	s.player.Dead()
	s.player.Update(s.gravity, s.deathJump)
	s.menu.Reset(math.FixOne*87, math.FixOne*102)

	s.scoreBanner.Set(math.FixOne*87, math.FixOne*-16)
	err := s.scoreBanner.Add()
	if err != nil {
		return err
	}

	s.bestBanner.Set(math.FixOne*87, math.FixOne*-16)
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

	s.score.Draw()
	return nil
}

func (s *Scene) Update(e *game.Engine) error {
	// this lerps the score and best banners in from off screen. It also uses the lut.Sin function to make the banners bob slightly
	s.t += 4
	s.scoreBanner.Set(math.FixOne*87, math.Lerp(math.FixOne*-16, math.FixOne*8, math.Clamp(s.t*2, 0, math.FixOne))+lut.Sin(s.t)+math.FixEighth)
	s.bestBanner.Set(math.FixOne*87, math.Lerp(math.FixOne*-16, math.FixOne*48, math.Clamp(s.t*2, 0, math.FixOne))+lut.Sin(s.t+math.FixThird)+math.FixEighth)

	s.player.Update(s.gravity, 0)
	s.score.Draw()
	s.highScore.Draw()
	s.menu.Update(e)
	if s.menu.selectCountDown > 0 && s.menu.selectCountDown > 10 {
		s.palFade += math.FixSixteenth
	}
	s.palFade = math.Clamp(s.palFade, 0, math.FixOne)
	e.PalFade(game.White, s.palFade)

	if s.menu.selectCountDown <= 0 {
		s.Restart = s.menu.restart
		s.Quit = s.menu.quit
	}
	return nil
}

// Hide hides all the assets associated with the scene
func (s *Scene) Hide() {
	s.menu.Hide()
	s.pillars.Hide()
	s.bestBanner.Remove()
	s.scoreBanner.Remove()
	s.highScore.Remove()
	s.score.Remove()
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

var (
	arrowSpinAnim = []game.Frame{
		{Index: 2, Len: 30},
		{Index: 1, Len: 10},
		{Index: 0, Len: 10},
		{Index: 0, HFlip: true, Len: 10},
		{Index: 1, HFlip: true, Len: 10},
		{Index: 2, Len: 10},
		{Index: 1, Len: 10},
		{Index: 0, Len: 10},
		{Index: 0, HFlip: true, Len: 10},
		{Index: 1, HFlip: true, Len: 10},
	}

	arrowBlinkAnim = []game.Frame{
		{Index: 2, Len: 7},
		{Index: 3, Len: 7},
	}
)

// newMenu creates a new game over menu
func newMenu(x, y math.Fix8, e *game.Engine) (*menu, error) {
	arrow := e.NewSprite(assets.SelectTileSet)
	arrow.X = x
	arrow.Y = y
	arrow.TileIndex = 2
	arrow.PlayAnimation(arrowSpinAnim)

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
func (m *menu) Update(e *game.Engine) {
	if m.restart || m.quit {
		m.selectCountDown--
		m.arrow.Update()
		return
	}

	if e.KeyJustPressed(key.Down) {
		m.arrow.Y = m.y + math.FixOne*12
	}
	if e.KeyJustPressed(key.Up) {
		m.arrow.Y = m.y
	}

	if m.selectStart > 0 {
		m.selectStart--
		m.arrow.Update()
		return
	}

	if e.KeyJustPressed(key.A) && m.arrow.Y == m.y {
		m.restart = true
	}
	if e.KeyJustPressed(key.A) && m.arrow.Y > m.y {
		m.quit = true
	}

	if m.restart || m.quit {
		m.arrow.PlayAnimation(arrowBlinkAnim)
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

// Hide hides the graphics associated with the menu
func (m *menu) Hide() {
	m.bg.Remove()
	m.arrow.Remove()
}

// Reset resets the menu back to it's initial state
func (m *menu) Reset(x, y math.Fix8) {
	m.arrow.TileIndex = 2
	m.arrow.PlayAnimation(arrowSpinAnim)

	m.x = x
	m.y = y
	m.selectStart = 30
	m.bg.VScroll = 0
	m.arrow.Y = y
	m.restart = false
	m.quit = false
}
