package audio

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

const (
	// FIFOA is the 4 bit FIFO register for DMA sound data. It is used for direct sound channel A
	// The data is treated as 8bit samples and is played in FIFO order with the least significant
	// byte played first
	FIFOA uintptr = 0x0400_00A0

	// FIFOB is the 4 bit FIFO register for DMA sound data. It is used for direct sound channel B
	// The data is treated as 8bit samples and is played in FIFO order with the least significant
	// byte played first
	FIFOB uintptr = 0x0400_00A4
)

// StatReg is the type used for the audio stats register. See Stat for more information on using this type
type StatReg uint16

// Stat is the controll register for enabling master sound. It also shows the status of the DMG
// channels. Notes that bits 0 - 3 are read only
//
// [0 - 3] Channel Activity - shows the stataus of the 4 DMG sound channels. These bits are read only.
//
// [7] Master Sound Enable - this is the master controll for soundon the GBA. If this is not enabled
//
//	no sound will play if this bit is not enabled. Also, note the when this bit is cleared all
//	registers in the range 0x0400_0060 to 0x0400_0081 are reset to zero and can not be written to
var Stat = (*StatReg)(unsafe.Pointer(memmap.IOAddr + 0x0084))

const (
	// MasterSoundEnable turns on master sound for the system
	MasterSoundEnable StatReg = 0x0080

	// MasterSoundDisable turns off master sound for the system
	MasterSoundDisable StatReg = 0x0000
)

// DSControllReg is the type used for the direct sound controll register. See DSControll for more infromation on using this type
type DSControllReg uint16

// DSControll is the controll register for Direct Sound. It also contains some bits related to DMG sounds channels.
// It has the following layout
//
// [0 - 1] DMG Volume - sets the volume ratio for the DMG sound channels
//   - Dmg25 - DMG volmue ratio 25%
//   - Dmg50 - DMG volume ratio 50%
//   - Dmg100 - DMG volume raito 100%
//
// [2] Direct Sound A Volume - sets the volume ratio for direct sound A
//   - A50 - DirectSound A volume ratio 50%
//   - A100 - DirectSound A volume ratio 100%
//
// [3] Direct Sound B Volume - sets the volume ratio for direct sound B
//   - B50 - DirectSound B volume ratio 50%
//   - B100 - DirectSound B volume ratio 100%
//
// [8 - 9] Direct Sound A Enable - turns on direct sound A for the left and right speakers
//   - AREnable - enable Direct Sound A on the right speaker
//   - ALEnable - enable Direct Sound A on the left speakers
//
// [A] Direct Sound A Timer Controll - set direct sound A to use timer 0 or 1
//   - ATimer0 - use timer 0 for DirectSound A
//   - ATimer1 - use timer 1 for DirectSound A
//
// [B] Direct Sound A FIFO Reset - set direct sound A to reset the FIFO buffer
//   - AReset - FIFO reset for DirectSound A. When using DMA for direct sound, this will cause the DMA to reset the FIFO buffer after it is used
//
// [C - D] Direct Sound B Enable - turns on direct sound B for the left and right speakers
//   - BREnable - enable direct sound B on the right speaker
//   - BLEnable - enable direct sound B on the left speaker
//
// [E] Direct Sound B Timer Controll - sets direct sound B to use timer 0 or 1
//   - BTimer0 - use timer 0 for DirectSound B                                               |
//   - BTimer1 - use timer 1 for DirectSound B                                               |
//
// [F] Direct Sound B FIFO Reset - sets direct sound B to reset the FIFO buffer
//   - BReset - FIFO reset for direct sound B. When using DMA for DS, this will cause DMA to reset the FIFO buffer after it is used
var DSControll = (*DSControllReg)(unsafe.Pointer(memmap.IOAddr + 0x00082))

const (
	// Dmg25 sets the DMG channels volume ratio to 25%.
	// it should be used with the RegSndDSCnt register
	Dmg25 DSControllReg = 0x0000

	// Dmg50 sets the DMG channels volume ratio to 50%
	// it should be used with the RegSndDSCnt register
	Dmg50 DSControllReg = 0x0001

	// Dmg100 sets the DMG channels volume ratio to 100%
	// it should be used with the RegSndDSCnt register
	Dmg100 DSControllReg = 0x0002

	// A50 sets the Direct Sound A volume ratio to 50%
	// it should be used with the RegSndDSCnt register
	A50 DSControllReg = 0x0000

	// A100 sets the Direct Sound A volume ratio to 100%
	// it should be used with the RegSndDSCnt register
	A100 DSControllReg = 0x0004

	// B50 sets the Direct Sound B volume ratio to 50%
	// it should be used with the RegSndDSCnt register
	B50 DSControllReg = 0x0000

	// B100 sets the Direct Sound B volume ratio to 100%
	// it should be used with the RegSndDSCnt register
	B100 DSControllReg = 0x0008

	// AREnable enables the Direct Sound A on the right speaker
	// it should be used with the RegSndDSCnt register
	AREnable DSControllReg = 0x0100

	// ALEnable enables the Direct Sound A on the left speaker
	// it should be used with the RegSndDSCnt register
	ALEnable DSControllReg = 0x0200

	// ATimer0 sets Direct Sound A to use timer 0
	// it should be used with the RegSndDSCnt register
	ATimer0 DSControllReg = 0x0000

	// ATimer1 sets Direct Sound A to use timer 1
	// it should be used with the RegSndDSCnt register
	ATimer1 DSControllReg = 0x0400

	// AReset is the FIFO reset for Direct Sound A. When using DMA for direct sound, this will cause
	// DMA to reset the FIFO buffer after it's used
	// it should be used with the RegSndDSCnt register
	AReset DSControllReg = 0x0800

	// BREnable enables the Direct Sound B on the right speaker
	// it should be used with the RegSndDSCnt register
	BREnable DSControllReg = 0x1000

	// BLEnable enables the Direct Sound B on the left speaker
	// it should be used with the RegSndDSCnt register
	BLEnable DSControllReg = 0x2000

	// BTimer0 sets Direct Sound B to use timer 0
	// it should be used with the RegSndDSCnt register
	BTimer0 DSControllReg = 0x0000

	// BTimer1 sets Direct Sound B to use timer 1
	// it should be used with the RegSndDSCnt register
	BTimer1 DSControllReg = 0x4000

	// BReset is the FIFO reset for Direct Sound B. When using DMA for direct sound, this will cause
	// DMA to reset the FIFO buffer after it's used
	// it should be used with the RegSndDSCnt register
	BReset DSControllReg = 0x8000
)
