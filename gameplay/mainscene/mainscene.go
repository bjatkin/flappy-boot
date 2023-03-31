package mainscene

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/fly"
	"github.com/bjatkin/flappy_boot/gameplay/gameover"
	"github.com/bjatkin/flappy_boot/gameplay/pillar"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/math"
)

type Scene struct {
	sky     *game.Background
	clouds  *game.Background
	pillars *pillar.BG
	player  *actor.Player

	fly      *fly.Scene
	gameOver *gameover.Scene

	activeScene game.Runable
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) Init(e *game.Engine) error {
	sky := e.NewBackground(assets.SkyTileMap, display.Priority3)
	clouds := e.NewBackground(assets.CloudsTileMap, display.Priority2)
	player := actor.NewPlayer(math.FixOne*40, math.FixOne*62, e.NewSprite(assets.PlayerTileSet))
	score := score.NewCounter(97, 24, e)

	fly := fly.NewScene(sky, clouds, player, score)
	err := fly.Init(e)
	if err != nil {
		return err
	}

	s.fly = fly
	s.activeScene = fly

	s.gameOver = gameover.NewScene(sky, clouds, player, score)
	return nil
}

func (s *Scene) Update(e *game.Engine, frame int) error {
	err := s.activeScene.Update(e, frame)
	if err != nil {
		return err
	}

	if s.activeScene == s.fly && s.fly.GameOver {
		s.activeScene = s.gameOver
		err = s.gameOver.Init(e)
		if err != nil {
			return err
		}
	}

	return nil
}
