//go:build !standalone

package memmap

// #include "memmap.h"
import "C"

import "unsafe"

const (
	// IO is the base memory address for the LCD I/O Registers
	IOAddr uintptr = 0x0400_0000

	// Keypad is the base memory address for the keypad input registers
	KeypadAddr uintptr = 0x0400_0130

	// Palette is the base memory address for the BG and Sprite palettes (1 Kbyte)
	PaletteAddr uintptr = 0x0500_0000

	// VRAM is the base address for video RAM (96 KBytes)
	VRAMAddr uintptr = 0x0600_0000

	// OAM is the base addres of all the object (sprite) attributes (1 Kbyte)
	OAMAddr uintptr = 0x0700_0000

	// SRAMAddr is the base memory address for SRAM in the gba pack memory
	SRAMAddr uintptr = 0x0E00_0000
)

// GetReg returns the volatile value of a 16 bit regiter
func GetReg[T reg](reg *T) T {
	v := C.GetReg((*C.ushort)(unsafe.Pointer(reg)))
	return T(v)
}

// SetReg sets the value of a 16 bit volitile register
func SetReg[T reg](reg *T, value T) {
	C.SetReg((*C.ushort)(unsafe.Pointer(reg)), C.ushort(value))
}
