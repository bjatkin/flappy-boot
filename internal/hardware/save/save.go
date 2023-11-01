package save

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// SRAMValue is a valid sram value. SRAM can only be written to one byte at a time
type SRAMValue byte

// SRAM is persistent storage that exists inside some GBA cartriages.
// it can be either SRAM which is battery powered or FRAM which is solid state memory
var sramStart = (*SRAMValue)(unsafe.Pointer(memmap.SRAMAddr))
var SRAM = unsafe.Slice(sramStart, 0x7FFF)

// WaitControll is the register is used to configure game pak access timings
// Game ROM is mirrored to three addresses in memory
// * 0x0800_0000 [wait state 0]
// * 0x0A00_0000 [wait state 1]
// * 0x0C00_0000 [wait state 2]
//
// This is useful for cases where there are several ROM chips with different access tims
// It is also needed for configuring SRAM/FRAM/EEPROM data
//
// [0 - 1] is the SRAM wait controll (4, 3, 2 or 8 cycles)
//
// [2 - 3] wait state 0 first access (4, 3, 2 or 8 cycles)
//
// [4] wait state 0 second access (2, 1 cycles)
//
// [5 - 6] wait state 1 first access (4, 3, 2, 8 cycles)
//
// [7] wait state 1 second access (4, 1 cycles)
//
// [8 - 9] wait state 2 first access (4, 3, 2, 8 cycles)
//
// [A] wait state 2 second access (8, 1 cycles)
//
// [B - C] PHI terminal output (Disabled, 4.19MHz, 8.38MHz, 16.78MHz)
//
// [E] Game ROM Prefetch Buffer (0=Disable, 1=Enable)
//
// [F] Game ROM type flag (Read Only) (0=GBA, 1=GBC)
var WaitControll = (*memmap.WaitControll)(unsafe.Pointer(memmap.IOAddr + 0x0204))

const (
	// SRAM4 sets the SRAM wait controll register to 4 cycles
	SRAM4 memmap.WaitControll = 0x0000

	// SRAM3 sets the SRAM wait controll register to 3 cycles
	SRAM3 memmap.WaitControll = 0x0001

	// SRAM2 sets the SRAM wait controll register to 2 cycles
	SRAM2 memmap.WaitControll = 0x0002

	// SRAM8 sets the SRAM wait controll register to 8 cycles
	SRAM8 memmap.WaitControll = 0x0003

	// W0First4 sets the 0 wait state controll to 4 cycles
	W0First4 memmap.WaitControll = 0x0000

	// W0First3 sets the 0 wait state controll to 3 cycles
	W0First3 memmap.WaitControll = 0x0004

	// W0First2 sets the 0 wait state controll to 2 cycles
	W0First2 memmap.WaitControll = 0x0008

	// W0First8 sets the 0 wait state controll to 8 cycles
	W0First8 memmap.WaitControll = 0x000C

	// W0Second2 sets the 0 wait state controll to 2 cycles for the second access
	W0Second2 memmap.WaitControll = 0x0000

	// W0Second1 sets the 0 wait state controll to 2 cycles for the second access
	W0Second1 memmap.WaitControll = 0x0010

	// W1First4 sets the 1 wait state controll to 4 cycles
	W1First4 memmap.WaitControll = 0x0000

	// W1First3 sets the 1 wait state controll to 3 cycles
	W1First3 memmap.WaitControll = 0x0020

	// W1First2 sets the 1 wait state controll to 2 cycles
	W1First2 memmap.WaitControll = 0x0040

	// W1First8 sets the 1 wait state contorll to 8 cycles
	W1First8 memmap.WaitControll = 0x0060

	// W1Second4 sets the 1 wait state controll to 4 cycles for the second access
	W1Second4 memmap.WaitControll = 0x0000

	// W1Second1 sets the 1 wait state controll to 1 cycles for the second access
	W1Second1 memmap.WaitControll = 0x0080

	// W2First4 sets the 2 wait state controll to 4 cycles
	W2First4 memmap.WaitControll = 0x0000

	// W2First3 sets the 2 wait state controll to 3 cycles
	W2First3 memmap.WaitControll = 0x0100

	// W2First2 sets the 2 wait state controll to 2 cycles
	W2First2 memmap.WaitControll = 0x0200

	// W2First8 sets the 2 wait state controll to 8 cycles
	W2First8 memmap.WaitControll = 0x0300

	// W2Second8 sets the 2 wait state controll to 8 cycles for the second access
	W2Second8 memmap.WaitControll = 0x0000

	// W2Second1 sets the 2 wait state controll to 1 cycle for the second access
	W2Second1 memmap.WaitControll = 0x0400

	// PHIOutDisable disables the PHI terminal output
	PHIOutDisable memmap.WaitControll = 0x0000

	// PHIOut4_19 setst the PHI terminal output to 4.19MHz
	PHIOut4_19 memmap.WaitControll = 0x0800

	// PHIOut8_38 sets the PHI terminal output to 8.38MHz
	PHIOut8_38 memmap.WaitControll = 0x1000

	// PHIOut16_78 sets the PHI terminal output to 16.78MHz
	PHIOut16_78 memmap.WaitControll = 0x1800

	// GameROMPrefetchDisable disables GBA cartriage ROM prefetch
	GameROMPrefetchDisable memmap.WaitControll = 0x0000

	// GameROMPrefetchEnable enables GBA cartriage ROM prefetch
	GameROMPrefetchEnable memmap.WaitControll = 0x4000
)
