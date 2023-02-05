package fix

// P8 is a fixed point number with a fix point at bit 8
type P8 int32

func New(i int, frac byte) P8 {
	return P8(i<<8) | P8(frac)
}

func (p P8) Int() int {
	return int(p >> 8)
}
