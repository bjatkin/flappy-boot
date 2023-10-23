package memmap

// TODO: should this be moved into it's own register package?

type reg interface {
	uint16 |
		AudioStat | DSControll |
		DisplayControll | DisplayStat | BGControll | DisplayVCount |
		Input | InputControll |
		WaitControll
}

// AudioStat is the type used for the audio stats register. See audio.Stat for more information on using this type
type AudioStat uint16

// DSControll is the type used for the direct sound controll register. See audio.DSControll for more infromation on using this type
type DSControll uint16

// DisplayControll is the type used for the display controll register, see display.Controll for more information on useing this type
type DisplayControll uint16

// DisplayStatReg is the type used for the display stats register, see display.Stat for more information on using this type
type DisplayStat uint16

// DisplayVCount is the type used for the display vertical line count register, see display.VCount for more information on using this type
type DisplayVCount uint16

// BGControll is the type used for the background controll registers, see display.BG#Controll for more information on using this type
type BGControll uint16

// Input is the type used for the input register, see key.Input for more information on using this type
type Input uint16

// InputControll is the type used for the input register, see key.Controll for more information on using this type
type InputControll uint16

// WaitControll is the type used for the system controll wait state register
type WaitControll uint16
