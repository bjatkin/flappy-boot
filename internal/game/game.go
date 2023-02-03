package game

import (
	"github.com/bjatkin/flappy_boot/internal/display"
	"github.com/bjatkin/flappy_boot/internal/key"
	"github.com/bjatkin/flappy_boot/internal/mode3"
	"github.com/bjatkin/flappy_boot/internal/sprite"
)

// Node is an interface that includes an Init, Update, and Draw function
type Node interface {
	Init() error
	Update(uint) (Node, error)
	Draw() error
	Unload() error
}

// Updater is a simplified version of a node that only includes an update function
type Updater interface {
	Update(uint) error
}

// Run starts a game by running the Node interface
func Run(node Node) {
	sprite.Reset()

	for {
		err := node.Init()
		if err != nil {
			exit(err)
		}

		var frame uint
		for {
			key.KeyPoll()
			next := (Node)(nil)
			next, err := node.Update(frame)
			if err != nil {
				exit(err)
			}

			frame++

			// block until it's safe to update the screen
			vSyncWait()

			// copy OAM buffer to OAM data proper here before we start drawing.
			// sprite data should be updated before the draw function so doing this
			// here ensures that only drawing code takes place in the draw function.
			// this also ensures that sprite data will never be updated mid-frame
			// which would cause screen tearing
			// sprite.CopyOAM()

			err = node.Draw()
			if err != nil {
				exit(err)
			}

			// unload this Node and then load the next one
			if next != nil {
				sprite.Reset()
				node.Unload()
				node = next
				break
			}
		}
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

// vSyncWait blocks while it waits for the screen to enter the vertical blank and then returns
func vSyncWait() {
	// TODO: leverage hardware interrupts rather than spinning the GBA cpu like this.
	// my guess is this is going to lead to pretty high power usage for no real benefit

	// wait till VDraw
	for display.VCount() >= 160 {
	}

	// wailt tile VBlank
	for display.VCount() < 160 {
	}
}
