package memmap

type Reg interface {
	AudioStat | DSControll |
		DisplayControll | DisplayStat | BGControll |
		Input | InputControll
}

// AudioStat is the type used for the audio stats register. See Stat for more information on using this type
type AudioStat uint16

// DSControll is the type used for the direct sound controll register. See DSControll for more infromation on using this type
type DSControll uint16

// DisplayControll is the type used for the display controll register, see Controll for more information on useing this type
type DisplayControll uint16

// DisplayStatReg is the type used for the display stats register, see Stat for more information on using this type
type DisplayStat uint16

// BGControll is the type used for the background controll registers, see BG#Controll for more information on using this type
type BGControll uint16

// Input is the type used for the input register, see Input for more information on using this type
type Input uint16

// InputControll is the type used for the input register, see Controll for more information on using this type
type InputControll uint16
