package mode3

import (
	"fmt"
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/display"
	hw_display "github.com/bjatkin/flappy_boot/internal/hardware/display"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/hardware/sprite"
)

var (
	// screenData points into the portion of VRAM tham maps to the screen when in mode 3
	// it is private to the package as the underlying data is exposed through ScreenData
	// with a more useful datatype
	screenData = memmap.VRAM[:hw_display.Width*hw_display.Height]

	// ScreenData is the color data for the screen, it uses 15 bit color
	//
	// WARNING: when copying large blocks of data do not use copy to update the screen
	// data this will cause graphical issues and the data will not be coppied correctly.
	// instead you shouls use CopyScreenData()
	ScreenData = *((*[]display.Color)(unsafe.Pointer(&screenData)))
)

// CopyDispData copies raw byte data to the screen. Copying a full screen of data
// will take more than a single frame to complete as this function does not use any
// DMA transfer channels
func CopyDispBytes(src []byte) {
	memmap.Copy16(screenData, src)
}

// CopyDispColors copies display.Color colors to the screen
// will take more than a single frame to complete as this function does not use any
// DMA transfer channels
func CopyDispColors(src []display.Color) {
	for i := range ScreenData {
		ScreenData[i] = src[i]
	}
}

// Enable sets the LCD to mode 3 and configures the VRAM using the provided options
// the following options are avilable
//
// * With1DSprites enables sprites and sets the VRAM to use 1D sprite tile mapping
// * With2DSprites enables sprites and sets the VRAM to use 2D sprite tile mapping
func Enable(options ...Option) {
	controll := hw_display.Mode3 | hw_display.BG2
	var bgControl hw_display.BGControllReg
	for _, opt := range options {
		controll, bgControl = opt(controll, bgControl)
	}

	memmap.SetReg(hw_display.Controll, controll)
	memmap.SetReg(hw_display.BG2Controll, bgControl)
}

// Option is a functional option that is used to define the parameters for
// the mode 3 graphics mode
type Option func(hw_display.ControllReg, hw_display.BGControllReg) (hw_display.ControllReg, hw_display.BGControllReg)

// With1DSprites enables sprites in mode 3 and will set the sprite memory mapping
// mode to be 1D rather than the default which is 2D sprite mapping
func With1DSprites() Option {
	return func(r hw_display.ControllReg, c hw_display.BGControllReg) (hw_display.ControllReg, hw_display.BGControllReg) {
		return r | hw_display.Sprites | hw_display.Sprite1D, c
	}
}

// With2DSprites enables sprites in mode 3 and will set the sprite memory mapping
// mode to be 2D which is the default
func With2DSprites() Option {
	return func(r hw_display.ControllReg, c hw_display.BGControllReg) (hw_display.ControllReg, hw_display.BGControllReg) {
		return r | hw_display.Sprites, c
	}
}

// WithWindows can enable windows 1 and 2 as well as the sprite window
func WithWindows(Win1, Win2, WinSpr bool) Option {
	return func(r hw_display.ControllReg, c hw_display.BGControllReg) (hw_display.ControllReg, hw_display.BGControllReg) {
		val := r
		if Win1 {
			val |= hw_display.Win1
		}
		if Win2 {
			val |= hw_display.Win2
		}
		if WinSpr {
			val |= hw_display.WinSpr
		}
		return val, c
	}
}

// WithPriority sets the priority for the background, it can be 0, 1, 2, or 3
// with priority 3 being rendered below 2, 2 rendered below 1 and 1 rendered below 0
func WithPriority(priority uint) Option {
	return func(r hw_display.ControllReg, c hw_display.BGControllReg) (hw_display.ControllReg, hw_display.BGControllReg) {
		if priority > 3 {
			priority = 3
		}
		return r, hw_display.BGControllReg(priority)
	}
}

// CopySprites loads sprite data into the video memory, if there is more sprite
// data and memory it will return an error and the data will not be loaded
func CopySprites(sprites []byte) error {
	if len(sprites) > len(sprite.Block1) {
		return fmt.Errorf(
			"sprite data is to large to fit into BlockB, got %d values but BlockB is only %d values in size",
			len(sprites),
			len(sprite.Block1),
		)
	}

	memmap.Copy16(sprite.Block1, sprites)
	return nil
}
