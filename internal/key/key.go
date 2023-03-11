package key

import (
	hw_key "github.com/bjatkin/flappy_boot/internal/hardware/key"
	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

type Key memmap.Input

// these constants remap the hardware register values to the Key type to make them eaiser to use
// in gameplay code
const (
	// A is the a button
	A = Key(hw_key.AMask)

	// B is the b button
	B = Key(hw_key.BMask)

	// Select is the select button
	Select = Key(hw_key.SelectMask)

	// Start is the start button
	Start = Key(hw_key.StartMask)

	// Right is the right directional button
	Right = Key(hw_key.RightMask)

	// Left is the left directional button
	Left = Key(hw_key.LeftMask)

	// Up is the up directional button
	Up = Key(hw_key.UpMask)

	// Down is the down directional button
	Down = Key(hw_key.DownMask)

	// L is the left shoulder button
	L = Key(hw_key.LMask)

	// R is the right shoulder button
	R = Key(hw_key.RMask)
)

var (
	// previous holds the state of the hardware key input register durring the last KeyPoll
	// it is used to check key transition states
	previous memmap.Input

	// current holds the state of the hardware key input register the current KeyPoll
	// it is used to check key transition states
	current memmap.Input
)

// KeyPoll reads they key input register and the current key state
func KeyPoll() {
	previous = current
	current = memmap.GetReg(hw_key.Input)
}

// Combo conbines multiple keys into a single key, it can be used to check if multiple keys are
// being held down simultaniously. It should not be used to check for key state changes as it is
// unlikely that all keys will be pressed down on the exact same KeyPoll
func Combo(keys ...Key) Key {
	var key Key

	for _, k := range keys {
		key |= k
	}

	return key
}

// StillPressed returns true if the key is being held down
// it will not return true if the key was first pressed durring this KeyPoll
func StillPressed(key Key) bool {
	return (^previous & ^current & memmap.Input(key) == memmap.Input(key))
}

// IsUp returns true if the key is not being held down
// it will not return true if the key was released durring this KeyPoll
func IsReleased(key Key) bool {
	return (^previous & ^current & memmap.Input(key)) == 0
}

// JustPressed returns true if the key was first pressed down durring this KeyPoll
func JustPressed(key Key) bool {
	return (previous & ^current & memmap.Input(key)) == memmap.Input(key)
}

// JustReleased returns true if the key was released durring this KeyPoll
func JustReleased(key Key) bool {
	return (^previous & current & memmap.Input(key)) == memmap.Input(key)
}

// StillPressed returns true if the key was first press down durring this KeyPoll
// or if it is currently being held down
func IsPressed(key Key) bool {
	return (^current & memmap.Input(key)) == memmap.Input(key)
}
