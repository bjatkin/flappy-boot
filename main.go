package main

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
)

func main() {
	engine := game.NewEngine()
	engine.Run(&Test{})
}

type Test struct {
	sky, pillars *game.Background
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

	return nil
}

func (t *Test) Update(e *game.Engine, frame int) error {
	return nil
}

func (t *Test) Next() (game.Runable, bool) {
	return nil, false
}
