package memmap

// #include "memmap.h"
import "C"

import (
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"unsafe"
)

const (
	// HalfKByte is the size of a kilobyte in uint16's
	HalfKByte = 0x200

	// CharBlockOffset is the size of a vram charblock in HalfKBytes
	CharBlockOffset = 8 * HalfKByte

	// ScreenBlockOffset is the size of a vram screen block in HalfKBytes
	ScreenBlockOffset = HalfKByte
)

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
)

// PaletteValue represents a valid color palette value
type PaletteValue uint16

// Palette is the system palette data, it consistes of 1kb and holds 16 bit color entries
// for both the background and sprite palettes
var Palette = *((*[]PaletteValue)(unsafe.Pointer(&reflect.SliceHeader{
	Data: PaletteAddr,

	// the gba has 2, 256 color palettes. PaletteValues are uint16 which is why these values are in HalfKBytes
	Cap: HalfKByte,
	Len: HalfKByte,
})))

// VRAMValue represents a valid VRAM value
type VRAMValue uint16

// VRAM is the system vram data, there are 96kb and depending on the mode
// this data can be used to achieve different effect.
var VRAM = *((*[]VRAMValue)(unsafe.Pointer(&reflect.SliceHeader{
	Data: VRAMAddr,

	// the gba has 96 KByte of VRAM, VRAMValues are uint16 which is why these values are in HalfKBytes
	Cap: 96 * HalfKByte,
	Len: 96 * HalfKByte,
})))

// OAMValue represents a valid OAM value
type OAMValue uint16

var OAM = *((*[]OAMValue)(unsafe.Pointer(&reflect.SliceHeader{
	Data: OAMAddr,

	// the gba has 128 normal sprite attributes and 32 affine attributes. These attributes
	// are interlaced resulting in 1kb of data. OAMValues are uint16 which is why these
	// values are in HalfKBytes
	Cap: HalfKByte,
	Len: HalfKByte,
})))

// values is a composit type of all the core memory value types
type values interface {
	PaletteValue | VRAMValue | OAMValue
}

// GetReg returns the volatile value of a 16 bit regiter
func GetReg[T any](reg *T) T {
	v := C.GetReg(C.uint(any(reg).(uintptr)))
	return any(v).(T)
}

// SetReg sets the value of a 16 bit volitile register
func SetReg[T any](reg *T, value T) {
	C.SetReg(C.uint(any(reg).(uintptr)), C.ushort(any(value).(uint16)))
}

// Copy16 coppies data from the source to the destination in 16 bit chunks
func Copy16[T values](dest []T, src []byte) {
	src16 := *((*[]T)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&src[0])),

		// these need to be uintptrs since the len and cap need to match the size of a pointer
		// on the GBA
		Cap: uintptr(len(src) / 2),
		Len: uintptr(len(src) / 2),
	})))

	for i := range src16 {
		if len(dest) <= i {
			return
		}
		dest[i] = src16[i]
	}
}

// Loads16 loads data from an embedded file into memory using the provided buffer
// the buffer should be less than 256kb to prevent compilation failures due to overflowing
// internal ram. Sizes less thean 32kb may lead to faster loading times as the buffer will fit
// into internal work ram. Ultimately it is up to the compiler if this happens however.
func Load16[T values](dest []T, src fs.File, buffer []byte) error {
	buffer16 := *((*[]T)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&buffer[0])),

		// these need to be uintptrs since the len and cap need to match the size of a pointer
		// on the GBA
		Cap: uintptr(len(buffer) / 2),
		Len: uintptr(len(buffer) / 2),
	})))

	var offset int
	for {
		offset++

		// this loads data into the buffer which shares the same memory location as buffer32
		_, err := src.Read(buffer)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("failed to load data into internal memory: %s", err)
		}

		for i := range dest {
			dest[i+offset] = buffer16[i]
		}
	}
}
