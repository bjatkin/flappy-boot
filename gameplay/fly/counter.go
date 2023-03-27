package fly

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/fix"
	"github.com/bjatkin/flappy_boot/internal/game"
)

// counter is a score counter used for displaying a score on the screen
type counter struct {
	score   [4]int
	convert []int
	digits  [4]*game.Sprite
	x, y    int
}

// newCounter creates a new counter struct
func newCounter(x, y int, e *game.Engine) *counter {
	c := &counter{
		convert: []int{16, 12, 28, 32, 24, 4, 8, 0, 36, 20},
		x:       x,
		y:       y,
	}
	for i := range c.digits {
		c.digits[i] = e.NewSprite(assets.NumbersTileSet)
		c.digits[i].X = fix.New(x+(i*11), 0)
		c.digits[i].Y = fix.New(y, 0)
	}

	return c
}

// Add adds 1 to the counters internal score
func (c *counter) Add() {
	c.score[3]++
	if c.score[3] > 9 {
		c.score[3] = 0
		c.score[2]++
	}
	if c.score[2] > 9 {
		c.score[2] = 0
		c.score[1]++
	}
	if c.score[1] > 9 {
		c.score[1] = 0
		c.score[0]++
	}
}

// Draw draw the graphics associated with the counter
func (c *counter) Draw() {
	var draw bool
	start := -1
	for i := range c.score {
		if i == len(c.score)-1 || c.score[i] > 0 {
			draw = true
		}
		if !draw {
			continue
		}
		if start == -1 {
			start = i * 6
		}

		c.digits[i].TileIndex = c.convert[c.score[i]]
		c.digits[i].X = fix.New(c.x-start+i*11, 0)
		c.digits[i].Add()
	}
}
