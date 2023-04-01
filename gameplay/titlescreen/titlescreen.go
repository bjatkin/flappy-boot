package titlescreen

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	sky, clouds *game.Background
	alter       *game.Background
	player      *actor.Player

	logo    *game.MetaSprite
	advance *game.MetaSprite
	press   *game.MetaSprite
	start   *game.MetaSprite
}

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
	}, nil
}

func (s *Scene) Init(e *game.Engine) error {
	s.logo.Set(math.FixOne*72, math.FixOne*20)
	if err := s.logo.Add(); err != nil {
		return err
	}

	s.advance.Set(math.FixOne*128, math.FixOne*40)
	if err := s.advance.Add(); err != nil {
		return err
	}

	s.press.Set(math.FixOne*72, math.FixOne*74)
	if err := s.press.Add(); err != nil {
		return err
	}

	s.start.Set(math.FixOne*128, math.FixOne*74)
	if err := s.start.Add(); err != nil {
		return err
	}

	if err := s.sky.Add(); err != nil {
		return err
	}

	if err := s.clouds.Add(); err != nil {
		return err
	}

	if err := s.alter.Add(); err != nil {
		return err
	}

	s.player.Sprite.X = math.FixOne * 114
	s.player.Sprite.Y = math.FixOne * 122
	if err := s.player.Show(); err != nil {
		return err
	}

	return nil
}

func (s *Scene) Update(e *game.Engine, frame int) error {
	s.clouds.HScroll += math.FixEighth
	return nil
}
