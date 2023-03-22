package main

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
)

func main() {
	engine := game.NewEngine()
	engine.Run(&Test{})
}

type Test struct {
	sky, pillars *game.Background
	player       *game.Sprite
}

func (t *Test) Init(e *game.Engine) error {
	t.pillars = e.NewBackground(assets.PillarsTileMap, hw_display.Priority2)
	err := t.pillars.Add()
	if err != nil {
		return err
	}

	t.sky = e.NewBackground(assets.SkyTileMap, hw_display.Priority3)
	err = t.sky.Add()
	if err != nil {
		return err
	}

	t.player = e.NewSprite(assets.PlayerTileSet)
	t.player.X = fix.New(20, 0)
	t.player.Y = fix.New(40, 0)
	err = t.player.Add()
	if err != nil {
		return err
	}

	return nil
}

func (t *Test) Update(e *game.Engine, frame int) error {
	return nil
}

func (t *Test) Next() (game.Runable, bool) {
	return nil, false
}
