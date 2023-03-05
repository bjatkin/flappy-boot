package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
)

//go:embed assets/gba
var assetFS embed.FS

func main() {
	engine := game.NewEngine()
	engine.Run(&Test{})
}

type Test struct {
	sky, grass *game.Background
	player     *game.Sprite
}

func (t *Test) Init(e *game.Engine) error {
	t.sky = e.NewBackground(assets.BackgroundTileMap, hw_display.Priority2)
	err := t.sky.Add()
	if err != nil {
		return err
	}

	t.grass = e.NewBackground(assets.SkyTileMap, hw_display.Priority3)
	err = t.grass.Add()
	if err != nil {
		return err
	}

	t.player = e.NewSprite(assets.Player)
	t.player.X = fix.New(40, 0)
	t.player.Y = fix.New(60, 0)
	err = t.player.Add()
	if err != nil {
		return err
	}

	return nil
}

func (t *Test) Update(e *game.Engine, frame int) error {
	if t.player.Y < fix.New(130, 0) {
		t.player.Y += fix.One / 4
	}

	return nil
}

func (t *Test) Draw(e *game.Engine) error {
	return nil
}

func (t *Test) Done() (game.Runable, bool) {
	return nil, false
}
