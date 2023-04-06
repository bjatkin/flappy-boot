package gameover

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/pillar"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/gameplay/state"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/lut"
	"github.com/bjatkin/flappy_boot/internal/math"
)

const (
	easeIn    = state.A
	main      = state.B
	confirmed = state.C
	fadeOut   = state.D
	done      = state.E
)

var sceneFrames = map[state.State]int{
	easeIn:    10,
	confirmed: 60,
	fadeOut:   30,
}

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

	state *state.Tracker

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

		state: &state.Tracker{
			SceneFrames: sceneFrames,
		},

		gravity:   math.FixQuarter,
		deathJump: -math.FixOne * 6,
	}, nil
}

func (s *Scene) Init(e *game.Engine) error {
	s.state.Init()
	s.Restart = false
	s.Quit = false

	s.player.Dead()
	s.player.Update(s.gravity, s.deathJump)
	s.menu.Reset()

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
	s.state.Update()
	s.player.Update(s.gravity, 0)
	s.score.Draw()
	s.highScore.Draw()

	if s.state.Is(main) {
		// this lerps the score and best banners in from off screen. It also uses the lut.Sin function to make the banners bob slightly
		t := math.Fix8(s.state.Frame() * 4)
		lerpT := math.Clamp(t*2, 0, math.FixOne)
		y := math.FixOne * 87
		ε := math.FixEighth
		s.scoreBanner.Set(y, math.Lerp(math.FixOne*-16, math.FixOne*8, lerpT)+lut.Sin(t)+ε)
		s.bestBanner.Set(y, math.Lerp(math.FixOne*-16, math.FixOne*48, lerpT)+lut.Sin(t+math.FixThird)+ε)

		s.menu.Update(e, s.state.Current())
		if s.menu.quit || s.menu.restart {
			s.state.Next()
		}

		return nil
	}

	if s.state.Is(confirmed | fadeOut) {
		s.menu.Update(e, s.state.Current())

	}

	if s.state.Is(fadeOut) {
		e.PalFade(game.White, s.state.Frac())

		return nil
	}

	if s.state.Is(done) {
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
	x, y          math.Fix8
	arrow         *game.Sprite
	bg            *game.Background
	restart, quit bool
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
		x:     x,
		y:     y,
		arrow: arrow,
		bg:    bg,
	}, nil
}

// Update updates the menu state each frame
func (m *menu) Update(e *game.Engine, s state.State) {
	m.arrow.Update()

	if s == easeIn || s == main {
		if e.KeyJustPressed(key.Down) {
			m.arrow.Y = m.y + math.FixOne*12
		}
		if e.KeyJustPressed(key.Up) {
			m.arrow.Y = m.y
		}
	}

	if s == main {
		if e.KeyJustPressed(key.A) && m.arrow.Y == m.y {
			m.restart = true
			m.arrow.PlayAnimation(arrowBlinkAnim)
		}
		if e.KeyJustPressed(key.A) && m.arrow.Y > m.y {
			m.quit = true
			m.arrow.PlayAnimation(arrowBlinkAnim)
		}
	}
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
func (m *menu) Reset() {
	m.arrow.TileIndex = 2
	m.arrow.PlayAnimation(arrowSpinAnim)
	m.arrow.Y = m.y

	m.bg.VScroll = 0
	m.restart = false
	m.quit = false
}
