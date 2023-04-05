package game

import (
	hw_key "github.com/bjatkin/flappy_boot/internal/hardware/key"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
	"github.com/bjatkin/flappy_boot/internal/key"
)

// keyPoll reads they key input register and the current key state
func (e *Engine) keyPoll() {
	e.previousKeys = e.currentKeys
	e.currentKeys = memmap.GetReg(hw_key.Input)
}

// StillPressed returns true if the key is being held down
// it will not return true if the key was first pressed durring this KeyPoll
func (e *Engine) KeyPressed(key key.Key) bool {
	return (^e.previousKeys & ^e.currentKeys & memmap.Input(key) == memmap.Input(key))
}

// IsReleased returns true if the key is not being held down
// it will not return true if the key was released durring this KeyPoll
func (e *Engine) KeyReleased(key key.Key) bool {
	return (^e.previousKeys & ^e.currentKeys & memmap.Input(key)) == 0
}

// JustPressed returns true if the key was first pressed down durring this KeyPoll
func (e *Engine) KeyJustPressed(key key.Key) bool {
	return (e.previousKeys & ^e.currentKeys & memmap.Input(key)) == memmap.Input(key)
}

// JustReleased returns true if the key was released durring this KeyPoll
func (e *Engine) KeyJustReleased(key key.Key) bool {
	return (^e.previousKeys & e.currentKeys & memmap.Input(key)) == memmap.Input(key)
}
