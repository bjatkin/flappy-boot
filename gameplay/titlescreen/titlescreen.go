package titlescreen

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
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

	startPressed int
	blinkOn      bool
	palFade      math.Fix8

	Done bool
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

		palFade: math.FixOne,
	}, nil
}

func (s *Scene) Init(e *game.Engine) error {
	s.Done = false
	s.startPressed = 0
	s.blinkOn = false
	s.palFade = math.FixOne

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

	s.player.Sprite.X = math.FixOne * 104
	s.player.Sprite.Y = math.FixOne * 124
	s.player.Sprite.TileIndex = 16
	s.player.Sprite.HFlip = false
	if err := s.player.Show(); err != nil {
		return err
	}

	return nil
}

func (s *Scene) Update(e *game.Engine) error {
	if s.startPressed > 0 && (e.Frame()-s.startPressed) > 50 {
		s.palFade += math.FixSixteenth
	}
	if s.startPressed == 0 {
		s.palFade -= math.FixSixteenth
	}
	s.palFade = math.Clamp(s.palFade, 0, math.FixOne)
	e.PalFade(game.White, s.palFade)

	s.clouds.HScroll += math.FixEighth
	if e.KeyJustPressed(key.Start) && s.startPressed == 0 {
		s.startPressed = e.Frame()
	}

	if s.startPressed > 0 && (e.Frame()-s.startPressed)%10 == 0 {
		if s.blinkOn {
			s.press.Set(math.FixOne*72, math.FixOne*74)
			s.start.Set(math.FixOne*128, math.FixOne*74)
		} else {
			s.press.Set(math.FixOne*240, 0)
			s.start.Set(math.FixOne*240, 0)
		}
		s.blinkOn = !s.blinkOn
	}

	s.Done = s.startPressed > 0 && (e.Frame()-s.startPressed) > 90

	if e.KeyJustPressed(key.B) {
		if err := s.alter.Add(); err != nil {
			return err
		}
	}

	return nil
}

func (s *Scene) Hide() {
	s.alter.Remove()
	s.press.Remove()
	s.start.Remove()
	s.logo.Remove()
	s.advance.Remove()
}
