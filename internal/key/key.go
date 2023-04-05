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
