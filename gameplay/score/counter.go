package score

import (
	"github.com/bjatkin/flappy_boot/internal/assets"
	"github.com/bjatkin/flappy_boot/internal/game"
	"github.com/bjatkin/flappy_boot/internal/math"
)

// Best is the best score that has been scored in this play session
var Best int

// Counter is a score Counter used for displaying a score on the screen
type Counter struct {
	score   [4]int
	convert []int
	digits  [4]*game.Sprite
	X, Y    int
}

// NewCounter creates a new counter struct
func NewCounter(x, y int, e *game.Engine) *Counter {
	c := &Counter{
		convert: []int{16, 12, 28, 32, 24, 4, 8, 0, 36, 20},
		X:       x,
		Y:       y,
	}
	for i := range c.digits {
		c.digits[i] = e.NewSprite(assets.NumbersTileSet)
		c.digits[i].X = math.NewFix8(x+(i*11), 0)
		c.digits[i].Y = math.NewFix8(y, 0)
	}

	return c
}

// Set sets the current value of the counter to the provided score.
// score must be between 0 and 9999
func (c *Counter) Set(score int) {
	c.score[0] = score / 1000
	c.score[1] = (score - c.score[0]*1000) / 100
	c.score[2] = (score - c.score[0]*1000 - c.score[1]*100) / 10
	c.score[3] = score - c.score[0]*1000 - c.score[1]*100 - c.score[2]*10
}

// Add adds 1 to the counters internal score
func (c *Counter) Add() {
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
func (c *Counter) Draw() {
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
		c.digits[i].X = math.NewFix8(c.X-start+i*11, 0)
		c.digits[i].Y = math.NewFix8(c.Y, 0)
		c.digits[i].Add()
	}
}

// Score
func (c *Counter) Score() int {
	return c.score[0]*1000 +
		c.score[1]*100 +
		c.score[2]*10 +
		c.score[3]
}
