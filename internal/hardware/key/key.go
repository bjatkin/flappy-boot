package key

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// Input is the register that is updated based on controller input, the bits in theses registers
// are LOW-ACTIVE meaning their value is CLEARED when a button is press and not the reverse as one
// you might expect. This register is read only and has the following layout.
//
// [0] A
// [1] B
// [2] Select
// [3] Start
// [4] Right
// [5] Left
// [6] Up
// [7] Down
// [8] R
// [9] L
var Input = (*memmap.Input)(unsafe.Pointer(memmap.KeypadAddr + 0x0000))

const (
	// AMask masks out every bit that is not the A button
	AMask memmap.Input = 0x0001

	// BMask masks out every bit that is not the B button
	BMask memmap.Input = 0x0002

	// SelectMask masks out every bit that is not the select button
	SelectMask memmap.Input = 0x0004

	// StartMask masks out every bit that is not the start button
	StartMask memmap.Input = 0x0008

	// RightMask masks out every bit that is not the right directional button
	RightMask memmap.Input = 0x0010

	// LeftMask masks out every bit that is not the left directional button
	LeftMask memmap.Input = 0x0020

	// UpMask masks out every bit that is not the up directional button
	UpMask memmap.Input = 0x0040

	// DownMask masks out every bit that is not the down directional button
	DownMask memmap.Input = 0x0080

	// LMask masks out every bit that is not the left shoulder button
	LMask memmap.Input = 0x0100

	// RMask masks out every bit that is not the right shoulder button
	RMask memmap.Input = 0x0200
)

// Controll is the register that sets interrupt controll data for the keypad
// It is a R/W register and has the following layout
//
// [0 - 9] Interrupt Keys - used to set the keys that will trigger an interrupt
// [E] Interrupt Enable - enables keypad interrupts
// [F] Interrupt Mode - sets the mode to use when triggering interupts
//   - And - triggers the interrupt only if all the interrupt keys are pressed
//   - Or - triggers the interrupt if any of the interrupt keys are pressed
var Controll = (*memmap.InputControll)(unsafe.Pointer(memmap.KeypadAddr + 0x0002))

const (
	// A is the a button
	A memmap.InputControll = 0x0001

	// B is the b button
	B memmap.InputControll = 0x0002

	// Select is the select button
	Select memmap.InputControll = 0x0004

	// Start is the start button
	Start memmap.InputControll = 0x0008

	// Right is the right directional button
	Right memmap.InputControll = 0x0010

	// Left is the left directional button
	Left memmap.InputControll = 0x0020

	// Up is the up directional button
	Up memmap.InputControll = 0x0040

	// Down is the down directional button
	Down memmap.InputControll = 0x0080

	// L is the left shoulder button
	L memmap.InputControll = 0x0100

	// R is the right shoulder button
	R memmap.InputControll = 0x0200

	// Interrupt turns on keypad interrupts
	Interrupt memmap.InputControll = 0x4000

	// And sets keypad interrupts to AND mode where all the specified keys must be pressed for
	// an interrupt to be triggered
	And memmap.InputControll = 0x8000

	// Or sets keypad interrupts to OR mode where an iterrupt will be fired if any of the
	// specified keys are pressed
	Or memmap.InputControll = 0x0000
)
