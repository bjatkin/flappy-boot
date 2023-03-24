package fix

// P8 is a fixed point number with a fix point at bit 8
type P8 int32

// New creates a new P8 with the given integer and fractional parts
func New(i int, frac byte) P8 {
	return P8(i<<8) | P8(frac)
}

// Int converts the P8 into a standard int
func (p P8) Int() int {
	return int(p >> 8)
}

const (
	// One is the value 1 in P8 format
	One P8 = 0x0100

	// Half is the value of 1/2 in P8 format
	Half P8 = 0x0080

	// Quarter is the value of 1/4 in P8 format
	Quarter P8 = 0x0040

	// Eighth is the value of 1/8 in P8 format
	Eighth P8 = 0x0020

	// Sixteenth is the value of 1/16 in P8 format
	Sixteenth P8 = 0x0010

	// Third is the value fo 1/3 in P8 format
	Third P8 = 0x0056
)
