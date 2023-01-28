#define REG(reg) *((volatile unsigned short*) (reg))

// GetReg returns the volitile avlue of a 16 bit register
volatile unsigned short GetReg(unsigned short* reg) {
    return REG(reg);
}

// SetReg sets the value of a 16 bit volitile register
void SetReg(unsigned short* reg, unsigned short value) {
    REG(reg) = value;
}
