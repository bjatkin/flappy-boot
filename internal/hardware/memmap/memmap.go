package memmap

import (
	"fmt"
	"io"
	"io/fs"
	"unsafe"
)

const (
	// KByte is the size of a kilobyte
	KByte = 0x0400

	// HalfKByte is the size of a kilobyte in uint16's
	HalfKByte = 0x0200

	// CharBlockOffset is the size of a vram charblock in HalfKBytes
	CharBlockOffset = 16 * HalfKByte

	// ScreenBlockOffset is the size of a vram screen block in HalfKBytes
	ScreenBlockOffset = 2 * HalfKByte

	// PaletteOffset is the size of a 16 color palette in 16 bit chunks
	PaletteOffset = 16

	// TileOffset4 is the size of a 4 bit per pixel tile in uint16's
	TileOffset4 = 16
)

// PaletteValue represents a valid color palette value
type PaletteValue uint16

// Palette is the system palette data, it consistes of 1kb and holds 16 bit color entries
// for both the background and sprite palettes
// the gba has 2, 256 color palettes. PaletteValues are uint16 which is why these values are in HalfKBytes
var paletteStart = (*PaletteValue)(unsafe.Pointer(PaletteAddr))
var Palette = unsafe.Slice(paletteStart, HalfKByte)

// VRAMValue represents a valid VRAM value
type VRAMValue uint16

// VRAM is the system vram data, there are 96kb and depending on the mode
// this data can be used to achieve different effect, such as drawing data to the screen and storing sprite gfx.
// the gba has 96 KByte of VRAM, VRAMValues are uint16 which is why these values are in HalfKBytes
var vramStart = (*VRAMValue)(unsafe.Pointer(VRAMAddr)) // vramStart is needed to prevent tinygo from failing
var VRAM = unsafe.Slice(vramStart, 96*HalfKByte)

// OAMValue represents a valid OAM value
type OAMValue uint16

// OAM is the object attribute data in the GBA hardware
// the gba has 128 normal sprite attributes and 32 affine attributes. These attributes
// are interlaced resulting in 1kb of data. OAMValues are uint16 which is why these
// values are in HalfKBytes
var oamStart = (*OAMValue)(unsafe.Pointer(OAMAddr)) // oamStart is needed to prevent tinygo from failing
var OAM = unsafe.Slice(oamStart, HalfKByte)

// values is a composit type of all the core memory value types
type values interface {
	PaletteValue | VRAMValue | OAMValue
}

// Copy16 coppies data from the source to the destination in 16 bit chunks
func Copy16[T values](dest []T, src []byte) {
	ptr := (*T)(unsafe.Pointer(&src[0]))
	src16 := unsafe.Slice(ptr, len(src)/2)

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
	ptr := (*T)(unsafe.Pointer(&buffer[0]))
	buffer16 := unsafe.Slice(ptr, len(buffer)/2)

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
