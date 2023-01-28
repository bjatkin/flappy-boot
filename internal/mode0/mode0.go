package mode0

import (
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// Enable sets the LCD to mode 0 and configures the VRAM using the provided options
// the following options are available
//
// - WithBG can be used to enable background layers 0 - 3
// - With1DSprites enables sprites and sets the VRAM to use 1D sprite tile mapping
// - With2DSprites enables sprites and sets the VRAM to use 2D sprite tile mapping
// - WithWindows can be used to enable the 2 background masking window and the 1 sprite masking window
func Enable(options ...Option) {
	controll := display.Mode0
	for _, opt := range options {
		controll = opt(controll)
	}

	memmap.SetReg(display.Controll, controll)
}

// Option is a functional option that is used to define the parameters for
// the mode 0 graphics mode
type Option func(memmap.DisplayControll) memmap.DisplayControll

// WithBG can be used to enable backgrounds 0 - 3
func WithBG(BG0, BG1, BG2, BG3 bool) Option {
	return func(r memmap.DisplayControll) memmap.DisplayControll {
		val := r
		if BG0 {
			val |= display.BG0
		}
		if BG1 {
			val |= display.BG1
		}
		if BG2 {
			val |= display.BG2
		}
		if BG3 {
			val |= display.BG3
		}
		return val
	}
}

// With1DSprites enables sprites in mode 0 and will set the sprite memory mapping
// mode to be 1D rather than the default which is 2D sprite mapping
func With1DSprites() Option {
	return func(r memmap.DisplayControll) memmap.DisplayControll {
		return r | display.Sprites | display.Sprite1D
	}
}

// With2DSprites enables sprites in mode 0 and will set the sprite memory mapping
// mode to be 2D which is the default
func With2DSprites() Option {
	return func(r memmap.DisplayControll) memmap.DisplayControll {
		return r | display.Sprites
	}
}

// WithWindows can enable windows 1 and 2 as well as the sprite window
func WithWindows(Win1, Win2, WinSpr bool) Option {
	return func(r memmap.DisplayControll) memmap.DisplayControll {
		val := r
		if Win1 {
			val |= display.Win1
		}
		if Win2 {
			val |= display.Win2
		}
		if WinSpr {
			val |= display.WinSpr
		}
		return val
	}
}
