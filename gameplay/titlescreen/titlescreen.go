package titlescreen

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/state"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/math"
)

const (
	fadeIn    = state.A
	main      = state.B
	confirmed = state.C
	fadeOut   = state.D
	done      = state.E
)

var sceneFrames = map[state.State]int{
	fadeIn:    30,
	confirmed: 60,
	fadeOut:   30,
}

// Scene is the intro scene for the game. It contains the title and allows the player to start the game
type Scene struct {
	sky, clouds *game.Background
	alter       *game.Background
	player      *actor.Player

	logo    *game.MetaSprite
	advance *game.MetaSprite
	press   *game.MetaSprite
	start   *game.MetaSprite

	state *state.Tracker

	Done bool
}

// NewScene creates a title screen scene
func NewScene(e *game.Engine, sky, clouds *game.Background, player *actor.Player) (*Scene, error) {
	logo, err := e.NewMetaSprite(
		[]math.V2{{X: 0}, {X: math.FixOne * 32}, {X: math.FixOne * 64}},
		[]int{16, 32, 0},
		assets.LogoTileSet,
	)
	if err != nil {
		return nil, err
	}

	advance, err := e.NewMetaSprite(
		[]math.V2{{X: 0}, {X: math.FixOne * 16}, {X: math.FixOne * 32}},
		[]int{8, 4, 0},
		assets.AdvanceTileSet,
	)
	if err != nil {
		return nil, err
	}

	press, err := e.NewMetaSprite(
		[]math.V2{{X: 0}, {X: math.FixOne * 16}, {X: math.FixOne * 32}},
		[]int{6, 2, 0},
		assets.StartTileSet,
	)
	if err != nil {
		return nil, err
	}

	start, err := e.NewMetaSprite(
		[]math.V2{{X: 0}, {X: math.FixOne * 16}, {X: math.FixOne * 32}},
		[]int{10, 8, 4},
		assets.StartTileSet,
	)
	if err != nil {
		return nil, err
	}

	return &Scene{
		sky:    sky,
		clouds: clouds,
		alter:  e.NewBackground(assets.MainmenuTileMap, display.Priority1),
		player: player,

		logo:    logo,
		advance: advance,
		press:   press,
		start:   start,

		state: &state.Tracker{
			SceneFrames: sceneFrames,
		},
	}, nil
}

// Init sets all the values to their initial steate for the Scene, it is safe to call repetedly
func (s *Scene) Init(e *game.Engine) error {
	s.state.Init()
	s.Done = false

	s.logo.Set(math.V2{X: math.FixOne * 72, Y: math.FixOne * 20})
	if err := s.logo.Show(); err != nil {
		return err
	}

	s.advance.Set(math.V2{X: math.FixOne * 128, Y: math.FixOne * 40})
	if err := s.advance.Show(); err != nil {
		return err
	}

	s.press.Set(math.V2{X: math.FixOne * 72, Y: math.FixOne * 74})
	if err := s.press.Show(); err != nil {
		return err
	}

	s.start.Set(math.V2{X: math.FixOne * 128, Y: math.FixOne * 74})
	if err := s.start.Show(); err != nil {
		return err
	}

	if err := s.sky.Show(); err != nil {
		return err
	}

	if err := s.clouds.Show(); err != nil {
		return err
	}

	if err := s.alter.Show(); err != nil {
		return err
	}

	s.player.Sprite.Pos = math.V2{X: math.FixOne * 104, Y: math.FixOne * 124}
	s.player.Sprite.TileIndex = 16
	s.player.Sprite.HFlip = false
	if err := s.player.Show(); err != nil {
		return err
	}

	return nil
}

// Update draws the title screen, updates the background and waits for the player to press start
func (s *Scene) Update(e *game.Engine) error {
	s.state.Update()
	s.clouds.HScroll += math.FixEighth

	if s.state.Is(fadeIn) {
		e.PalFade(game.White, math.FixOne-s.state.Frac())
		return nil
	}

	if s.state.Is(main) {
		if e.KeyJustPressed(key.Start) {
			s.state.Next()
		}
		return nil
	}

	if s.state.Is(confirmed | fadeOut) {
		if s.state.Frame()>>3%2 == 0 {
			s.press.Set(math.V2{X: math.FixOne * 72, Y: math.FixOne * 74})
			s.start.Set(math.V2{X: math.FixOne * 128, Y: math.FixOne * 74})
		} else {
			s.press.Set(math.V2{X: math.FixOne * 240})
			s.start.Set(math.V2{X: math.FixOne * 240})
		}

	}

	if s.state.Is(fadeOut) {
		e.PalFade(game.White, s.state.Frac())
		return nil
	}

	if s.state.Is(done) {
		s.Done = true
	}

	return nil
}

// Hide removes the title screen from view
func (s *Scene) Hide() {
	s.alter.Hide()
	s.press.Hide()
	s.start.Hide()
	s.logo.Hide()
	s.advance.Hide()
}
