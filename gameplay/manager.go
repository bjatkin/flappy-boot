package gameplay

import (
	"github.com/bjatkin/flappy_boot/gameplay/actor"
	"github.com/bjatkin/flappy_boot/gameplay/fly"
	"github.com/bjatkin/flappy_boot/gameplay/gameover"
	"github.com/bjatkin/flappy_boot/gameplay/pillar"
	"github.com/bjatkin/flappy_boot/gameplay/score"
	"github.com/bjatkin/flappy_boot/gameplay/titlescreen"
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

	fly         *fly.Scene
	gameOver    *gameover.Scene
	titleScreen *titlescreen.Scene

	activeScene game.Runable
	initErr     error
}

// scoreToSaveData converts a score into a byte slice that can be saved in FRAM
func scoreToSaveData(score int) []byte {
	return []byte{byte(score >> 8), byte(score)}
}

// saveDataToScore converts data from FRAM into a valid score
func saveDataToScore(data []byte) int {
	return int(data[0])<<8 | int(data[1])
}

func NewManager(e *game.Engine) *Manager {
	sky := e.NewBackground(assets.SkyTileMap, display.Priority3)
	clouds := e.NewBackground(assets.CloudsTileMap, display.Priority2)
	player := actor.NewPlayer(math.V2{X: math.FixOne * 32, Y: math.FixOne * 62}, e.NewSprite(assets.PlayerAnimTileSet))
	pillars := pillar.NewBG(100, e.NewBackground(assets.PillarsTileMap, display.Priority1))
	roundScore := score.NewCounter(97, 28, e)

	highScore := score.NewCounter(240, 0, e)
	// load the highScore from FRAM
	highScore.Set(saveDataToScore(e.LoadData(2)))

	var initErr error
	over, err := gameover.NewScene(e, sky, clouds, pillars, player, roundScore, highScore)
	if err != nil {
		initErr = err
	}

	title, err := titlescreen.NewScene(e, sky, clouds, player)
	if err != nil {
		initErr = err
	}

	return &Manager{
		sky:        sky,
		clouds:     clouds,
		player:     player,
		roundScore: roundScore,
		highScore:  highScore,

		fly:         fly.NewScene(e, sky, clouds, pillars, player, roundScore),
		gameOver:    over,
		titleScreen: title,
		initErr:     initErr,
	}
}

func (s *Manager) Init(e *game.Engine) error {
	if s.initErr != nil {
		return s.initErr
	}

	err := s.titleScreen.Init(e)
	if err != nil {
		return err
	}

	s.activeScene = s.titleScreen
	return nil
}

func (s *Manager) Update(e *game.Engine) error {
	err := s.activeScene.Update(e)
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
				// save the new high score in FRAM
				e.SaveData(scoreToSaveData(s.roundScore.Score()))
			}
		}
	case s.gameOver:
		if s.gameOver.Restart {
			s.gameOver.Hide()
			s.activeScene = s.fly
			if err = s.fly.Init(e); err != nil {
				return err
			}
		}
		if s.gameOver.Quit {
			s.gameOver.Hide()
			s.activeScene = s.titleScreen
			if err = s.titleScreen.Init(e); err != nil {
				return err
			}
		}
	case s.titleScreen:
		if s.titleScreen.Done {
			s.titleScreen.Hide()
			s.activeScene = s.fly
			err := s.fly.Init(e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
