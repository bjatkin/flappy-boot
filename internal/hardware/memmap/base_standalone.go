//go:build standalone

package memmap

import (
	"unsafe"
)

var (
	IORegBlock   = [KByte]byte{}
	PaletteBlock = [2 * KByte]byte{}
	VRAMBlock    = [96 * KByte]byte{}
	OAMBlock     = [KByte]byte{}
	SRAMBlock    = [4 * KByte]byte{}
)

// These values reflect the values in base.go, the difference here is they are dynamically
// allocated by go, allowing all code that references them to run successfully on normal
// systems where memory is not maped to the specific GBA addresses

var (
	// IO is the base memory address for the LCD I/O Registers
	IOAddr = uintptr(unsafe.Pointer(&IORegBlock))

	// Keypad is the base memory address for the keypad input registers
	KeypadAddr = IOAddr + 0x0130

	// Palette is the base memory address for the BG and Sprite palettes (1 Kbyte)
	PaletteAddr = uintptr(unsafe.Pointer(&PaletteBlock))

	// VRAM is the base address for video RAM (96 KBytes)
	VRAMAddr = uintptr(unsafe.Pointer(&VRAMBlock))

	// OAM is the base addres of all the object (sprite) attributes (1 Kbyte)
	OAMAddr = uintptr(unsafe.Pointer(&OAMBlock))

	// SRAMAddr is the base memory address for SRAM in the gba pack memory
	SRAMAddr = uintptr(unsafe.Pointer(&SRAMBlock))
)

// GetReg replaces GetReg from base.go, this is because durring emulation,
// volitile memory access is not nessiary
func GetReg[T reg](reg *T) T {
	return *reg
}

// SetReg replaces SetReg form base.go this is because durring emulation,
// volitile memory access is not nessisary
func SetReg[T reg](reg *T, value T) {
	*reg = value
}
