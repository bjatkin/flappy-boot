package gameplay

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/fly"
	"github.com/bjatkin/flappy_boot/gameplay/gameover"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Manager struct {
	sky        *game.Background
	clouds     *game.Background
	player     *actor.Player
	roundScore *score.Counter
	highScore  *score.Counter

	fly      *fly.Scene
	gameOver *gameover.Scene

	activeScene game.Runable
	initErr     error
}

func NewManager(e *game.Engine) *Manager {
	sky := e.NewBackground(assets.SkyTileMap, display.Priority3)
	clouds := e.NewBackground(assets.CloudsTileMap, display.Priority2)
	player := actor.NewPlayer(math.FixOne*40, math.FixOne*62, e.NewSprite(assets.PlayerTileSet))
	roundScore := score.NewCounter(97, 28, e)
	highScore := score.NewCounter(240, 0, e)
	over, err := gameover.NewScene(e, sky, clouds, player, roundScore, highScore)

	return &Manager{
		sky:        sky,
		clouds:     clouds,
		player:     player,
		roundScore: roundScore,
		highScore:  highScore,

		fly:      fly.NewScene(e, sky, clouds, player, roundScore),
		gameOver: over,
		initErr:  err,
	}
}

func (s *Manager) Init(e *game.Engine) error {
	if s.initErr != nil {
		return s.initErr
	}

	err := s.fly.Init(e)
	if err != nil {
		return err
	}

	s.activeScene = s.fly
	return nil
}

func (s *Manager) Update(e *game.Engine, frame int) error {
	err := s.activeScene.Update(e, frame)
	if err != nil {
		return err
	}

	switch s.activeScene {
	case s.fly:
		if s.fly.GameOver {
			s.activeScene = s.gameOver
			if err = s.gameOver.Init(e); err != nil {
				return err
			}
			if s.roundScore.Score() > s.highScore.Score() {
				s.highScore.Set(s.roundScore.Score())
			}
		}
	case s.gameOver:
		if s.gameOver.Restart {
			s.activeScene = s.fly
			if err = s.fly.Init(e); err != nil {
				return err
			}
			s.gameOver.Hide()
		}
	}

	return nil
}
