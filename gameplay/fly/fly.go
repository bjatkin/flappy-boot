package fly

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/pillar"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	GameOver bool

	sky         *game.Background
	clouds      *game.Background
	pillars     *pillar.BG
	player      *actor.Player
	score       *score.Counter
	scrollSpeed math.Fix8

	gravity    math.Fix8
	ground     math.Fix8
	jumpHeight math.Fix8
}

func NewScene(e *game.Engine, sky, clouds *game.Background, pillars *pillar.BG, player *actor.Player, score *score.Counter) *Scene {
	return &Scene{
		scrollSpeed: math.FixOne + math.FixEighth,
		gravity:     math.FixQuarter,
		ground:      math.FixOne * 147,
		jumpHeight:  -math.FixOne * 3,

		pillars: pillars,
		sky:     sky,
		clouds:  clouds,
		player:  player,
		score:   score,
	}
}

func (s *Scene) Init(e *game.Engine) error {
	s.GameOver = false
	s.player.Reset(math.FixOne*40, math.FixOne*62)
	s.pillars.Reset()
	s.score.Set(0)

	err := s.pillars.Show()
	if err != nil {
		return err
	}

	err = s.sky.Add()
	if err != nil {
		return err
	}

	err = s.clouds.Add()
	if err != nil {
		return err
	}
	s.clouds.VScroll = 20

	err = s.player.Show()
	if err != nil {
		return err
	}

	s.score.Draw()

	return nil
}

func (s *Scene) Update(e *game.Engine, frame int) error {
	var jump math.Fix8
	if key.JustPressed(key.A) {
		s.pillars.Start()
		s.player.Start()
		jump = -math.FixOne * 3
	}

	s.player.Update(s.gravity, jump)
	if s.player.Rect().Y2 >= s.ground.Int() {
		s.GameOver = true
	}

	s.sky.HScroll += s.scrollSpeed / 3
	err := s.sky.Add()
	if err != nil {
		return err
	}

	s.clouds.HScroll += s.scrollSpeed / 2
	err = s.clouds.Add()
	if err != nil {
		return err
	}

	s.pillars.Update(s.scrollSpeed)
	err = s.pillars.Show()
	if err != nil {
		return err
	}

	if s.pillars.CheckPoint(s.player.Rect()) {
		s.score.Add()
	}

	s.score.Draw()

	if s.pillars.CollisionCheck(s.player.Rect()) {
		s.GameOver = true
	}

	return nil
}
