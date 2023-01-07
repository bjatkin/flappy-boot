#define REG(reg) *((volatile unsigned short*) (reg))

// GetReg returns the volitile avlue of a 16 bit register
volatile unsigned short GetReg(unsigned int reg) {
    return REG(reg);
}

// SetReg sets the value of a 16 bit volitile register
void SetReg(unsigned int reg, unsigned short value) {
    REG(reg) = value;
}
