package registers

import (
	"github.com/bjatkin/flappy_boot/internal/hardware/audio"
	"github.com/bjatkin/flappy_boot/internal/hardware/display"
)

// TODO: make a genric type for all the register types
type Reg interface {
	audio.StatReg |
		display.BGControllReg | display.ControllReg | display.StatReg
}
