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

	menu, err := newMenu(math.V2{X: math.FixOne * 87, Y: math.FixOne * 102}, e)
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
	s.menu.Init()

	s.scoreBanner.Set(math.V2{X: math.FixOne * 87, Y: math.FixOne * -16})
	err := s.scoreBanner.Show()
	if err != nil {
		return err
	}

	s.bestBanner.Set(math.V2{X: math.FixOne * 87, Y: math.FixOne * -16})
	err = s.bestBanner.Show()
	if err != nil {
		return err
	}

	s.highScore.X = 97
	s.highScore.Y = 70

	err = s.menu.Show()
	if err != nil {
		return err
	}

	s.score.Update()
	return nil
}

func (s *Scene) Update(e *game.Engine) error {
	s.state.Update()
	s.player.Update(s.gravity, 0)
	s.score.Update()
	s.highScore.Update()

	if s.state.Is(main) {
		// this lerps the score and best banners in from off screen. It also uses the lut.Sin function to make the banners bob slightly
		t := math.Fix8(s.state.Frame() * 4)
		lerpT := math.Clamp(t*2, 0, math.FixOne)
		x := math.FixOne * 87
		ε := math.FixEighth
		s.scoreBanner.Set(math.V2{
			X: x,
			Y: math.Lerp(math.FixOne*-16, math.FixOne*8, lerpT) + lut.Sin(t) + ε,
		})
		s.bestBanner.Set(math.V2{
			X: x,
			Y: math.Lerp(math.FixOne*-16, math.FixOne*48, lerpT) + lut.Sin(t+math.FixThird) + ε,
		})

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
	s.bestBanner.Hide()
	s.scoreBanner.Hide()
	s.highScore.Hide()
	s.score.Hide()
}

// menu is a simple game over menu
type menu struct {
	pos           math.V2
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
func newMenu(pos math.V2, e *game.Engine) (*menu, error) {
	arrow := e.NewSprite(assets.SelectTileSet)
	arrow.Pos = pos
	arrow.TileIndex = 2
	arrow.PlayAnimation(arrowSpinAnim)

	bg := e.NewBackground(assets.BluebgTileMap, display.Priority0)

	return &menu{
		pos:   pos,
		arrow: arrow,
		bg:    bg,
	}, nil
}

// Init resets the menu back to it's initial state
func (m *menu) Init() {
	m.arrow.TileIndex = 2
	m.arrow.PlayAnimation(arrowSpinAnim)
	m.arrow.Pos.Y = m.pos.Y

	m.bg.VScroll = 0
	m.restart = false
	m.quit = false
}

// Update updates the menu state each frame
func (m *menu) Update(e *game.Engine, s state.State) {
	m.arrow.Update()

	if s == easeIn || s == main {
		if e.KeyJustPressed(key.Down) {
			m.arrow.Pos.Y = m.pos.Y + math.FixOne*12
		}
		if e.KeyJustPressed(key.Up) {
			m.arrow.Pos.Y = m.pos.Y
		}
	}

	if s == main {
		if e.KeyJustPressed(key.A) && m.arrow.Pos.Y == m.pos.Y {
			m.restart = true
			m.arrow.PlayAnimation(arrowBlinkAnim)
		}
		if e.KeyJustPressed(key.A) && m.arrow.Pos.Y > m.pos.Y {
			m.quit = true
			m.arrow.PlayAnimation(arrowBlinkAnim)
		}
	}
}

// Show adds menu sprites and backgrounds into active engine memory
func (m *menu) Show() error {
	err := m.bg.Show()
	if err != nil {
		return err
	}

	err = m.arrow.Show()
	if err != nil {
		return err
	}

	return nil
}

// Hide hides the graphics associated with the menu
func (m *menu) Hide() {
	m.bg.Hide()
	m.arrow.Hide()
}
