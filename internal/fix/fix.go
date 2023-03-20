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

// One is the value 1 in P8 format
const One P8 = 0x0100
