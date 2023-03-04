package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/display"
	"github.com/bjatkin/flappy_boot/internal/game"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/mode3"
)

//go:embed assets/gba
var assetFS embed.FS

func main() {
	// game.Run(gameplay.NewDemo(assetFS))
	// mode0.Enable(mode0.With1DSprites(), mode0.WithBG(true, true, true, true))

	// debug()
	// return

	var err error //remove
	engine := game.NewEngine()
	grassBG := engine.NewBackground(assets.BackgroundTileSet, assets.BackgroundTileMap, hw_display.Priority2)
	err = grassBG.Load()
	if err != nil {
		exit(err)
	}

	err = grassBG.Add()
	if err != nil {
		exit(err)
	}

	skyBG := engine.NewBackground(assets.SkyTileSet, assets.SkyTileMap, hw_display.Priority3)
	err = skyBG.Load()
	if err != nil {
		exit(err)
	}

	err = skyBG.Add()
	if err != nil {
		exit(err)
	}

	for {
	}
}

// func debug() {
// 	bg := assets.NewSky()
// 	fmt.Println("width:", bg.Height, "height:", bg.Height)
// }

// exit exits the game loop and draws error infromation to the screen
func exit(err error) {
	// TODO: this should be updated to use a 'system' font to write out the
	// error data to make debugging easier
	mode3.Enable()

	// Draw red to the screen so we can tell there was an error
	red := display.RGB15(0, 0, 31)
	for i := 10; i < 240*160; i++ {
		mode3.ScreenData[i] = red
	}

	// block forever
	for {
	}
}
