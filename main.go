package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/key"
)

//go:embed assets/gba
var assetFS embed.FS

func main() {
	engine := game.NewEngine()
	engine.Run(&Test{})
}

type Test struct {
	sky, grass *game.Background

	player actor
}

type actor struct {
	*game.Sprite
	Dx, Dy fix.P8
}

var (
	gravity = fix.New(0, 64)

	jump = fix.New(-3, 0)

	groundY = fix.New(131, 0)

	scrollSpeed = 1
)

func (t *Test) Init(e *game.Engine) error {
	t.grass = e.NewBackground(assets.BackgroundTileMap, hw_display.Priority2)
	err := t.grass.Add()
	if err != nil {
		return err
	}

	t.sky = e.NewBackground(assets.SkyTileMap, hw_display.Priority3)
	err = t.sky.Add()
	if err != nil {
		return err
	}

	t.player.Sprite = e.NewSprite(assets.Player)
	t.player.X = fix.New(40, 0)
	t.player.Y = fix.New(72, 0)
	err = t.player.Add()
	if err != nil {
		return err
	}

	return nil
}

func (t *Test) Update(e *game.Engine, frame int) error {
	if key.JustPressed(key.A) {
		t.player.Dy = jump
	}

	t.player.Dy += gravity
	t.player.Y += t.player.Dy

	if t.player.Y > groundY {
		t.player.Y = groundY
		t.player.Dy = 0
	}

	if t.player.Y < 0 {
		t.player.Y = 0
		t.player.Dy = 0
	}

	t.grass.Scroll(scrollSpeed, 0)
	if key.JustPressed(key.B) {
		for i := 0; i < 2; i++ {
			for ii := 0; ii < 2; ii++ {
				t.sky.SetTile(i, ii, 7)
			}
		}
		// BOTTOM PILLAR
		t.grass.SetTile(12, 14, 10)
		t.grass.SetTile(13, 14, 22)
		t.grass.SetTile(14, 14, 20)
		t.grass.SetTile(15, 14, 11)

		// TOP PILLAR
		t.grass.SetTile(12, 2, 10)
		t.grass.SetTile(13, 2, 22)
		t.grass.SetTile(14, 2, 20)
		t.grass.SetTile(15, 2, 11)
	}

	return nil
}

func (t *Test) Draw(e *game.Engine) error {
	return nil
}

func (t *Test) Done() (game.Runable, bool) {
	return nil, false
}
