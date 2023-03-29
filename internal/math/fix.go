package math

// Fix8 is a fixed point number with a fix point at bit 8
type Fix8 int32

// New creates a new P8 with the given integer and fractional parts
func NewFix8(i int, frac byte) Fix8 {
	return Fix8(i<<8) | Fix8(frac)
}

// Int converts the P8 into a standard int
func (p Fix8) Int() int {
	return int(p >> 8)
}

// Int16 converts the P8 into an int16
func (p Fix8) Uint16() uint16 {
	return uint16(p >> 8)
}

const (
	// FixOne is the value 1 in P8 format
	FixOne Fix8 = 0x0100

	// FixHalf is the value of 1/2 in P8 format
	FixHalf Fix8 = 0x0080

	// FixQuarter is the value of 1/4 in P8 format
	FixQuarter Fix8 = 0x0040

	// FixEighth is the value of 1/8 in P8 format
	FixEighth Fix8 = 0x0020

	// Sixteenth is the value of 1/16 in P8 format
	Sixteenth Fix8 = 0x0010

	// Third is the value fo 1/3 in P8 format
	Third Fix8 = 0x0056
)
