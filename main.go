package main

import (
	"embed"

	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/display"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/mode0"
	"github.com/bjatkin/flappy_boot/internal/mode3"
)

//go:embed assets/gba
var assetFS embed.FS

func main() {
	//game.Run(gameplay.NewDemo(assetFS))

	mode0.Enable(mode0.With1DSprites())

	engine := game.NewEngine()
	grassBG := engine.NewBackground(assets.BackgroundTileSet, &assets.BackgroundTileMap)
	err := grassBG.Load()
	if err != nil {
		exit(err)
	}

	err = grassBG.Add()
	if err != nil {
		exit(err)
	}
}

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
