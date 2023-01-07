package dma

// TODO: update this to be in line with new standards
type Register uint

const (
	// RegDMA1SAD is the DMA Source Address for the DMA1 transfer channel
	// it is a 32 bit register
	RegDMA1SAD Register = 0x0400_00BE

	// RegDMA1DAD is the DMA Destination address for the DMA1 transfer channel
	// it is a 32 bit register
	RegDMA1DAD Register = 0x0400_00C0

	// RegDMA1Cnt is the controll register for the DMA1 transfer channel.
	// It can be used to start the DMA trasfer and has the following layout.
	//
	// [0h - Fh] Count - the number of transfers to make. This is an unsigned 16 bit value
	//
	// [15h - 16h] Destination Address Controll - controlls how the destination address changes after each transfer
	//     * DestAddrInc - sets the destination address to increment as each data segment is coppied
	//     * DestAddrDec - sets the destination address to decrement as each data segment is coppied
	//     * DestAddrFix - fixes the destination address. It will not change
	//     * DestAddrRe  - sets the destination address to increment as each data segment is coppied
	//                     resets the address each time the DMA repeats
	//
	// [17h - 18h] Source Address Controll - controlls how the source address changes after each transfer
	//     * SrcAddrInc - sets the source address to increment as each data segment is coppied
	//     * SrcAddrDec - sets the source address to decrement as each data segment is coppied
	//
	// [19h] Repeat Mode - changes how the Enable bit behaves after the DMA transfer has completed
	//     * RepeatOn - set the DMA transfer to repeat at each HBlank, VBlank, or empty FIFO buffer
	// [1Ah] Transfer Size - tels the DMA how many bits of data to transfer at a time
	//     * Transfer16 - transfers 16 bits of data at a time
	//     * Transfer32 - transfers 32 bits of data at a time
	// [1Dh - Ch] Start Mode - sets the start time for the DMA transfer
	//     * StartNow - start transfering data immediately
	//     * StartVBlank - start transfering data at the next VBlank
	//     * StartHBlank - start transfering data at the next HBlank
	//     * StartFIFO - start transfering data when the configured FIFO buffer is emptied
	//                   transfer count should be set to 1 and Transfer Size should be 32 bits
	// [1Eh] DMA Interrupt - can be use to enable DMA interrupts when the DMA transfer is finished
	//     * IRQEnable - raise an interrupt when finished
	//     * IRQDisable - do not rais and interrupt when finished
	// [1Fh] DMA Enable - turns on/ off the DMA channel
	//     * DMAOn - enable the DMA transfer
	//     * DMAOff - disable the DMA transfer
	RegDMA1Cnt Register = 0x0400_00C4
)

const (
	// SystemClock is the exact number of CPU ticks per cycle (16.78MHz)
	SystemClock = 16_777_216

	// ScreenRefresh is the exact number of cycles between each screen refresh
	ScreenRefresh = 280_896

	// DestAddrInc sets the destination address to increase after transfering each data segment
	// should be used with the RegDMA1CntH register
	DestAddrInc = 0x0000

	// DestAddrDec sets the destination address to decrease after transfering each data segment
	// should be used with the RegDMA1CntH register
	DestAddrDec = 0x0020

	// DestAddrFix sets the destination address to be fixed
	// should be used with the RegDMA1CntH register
	DestAddrFix = 0x0040

	// DestAddrRe sets the destination address to increase after transfering each data segment
	// and then reset each time the DMA repeats so each copy starts at the same destination
	// should be used with the RegDMA1CntH register
	DestAddrRe = 0x0060

	// SrcAddrInc sets the source address to increase after transfering each data segment
	// should be used with the RegDMA1CntH register
	SrcAddrInc = 0x0000

	// SrcAddrDec sets the source address to decrease after transfering each data segment
	// should be used with the RegDMA1CntH register
	SrcAddrDec = 0x0080

	// RepeatOn sets the DMA to repeat copying data at each VBlank HBlank or FIFO buffer empty
	// depending on the start time
	// should be used with the RegDMA1CntH register
	RepeatOn = 0x0200

	// Transfer16 sets the DMA to transfer 16 bits at a time
	// should be used with the RegDMA1CntH register
	Transfer16 = 0x0000

	// Transfer32 sets the DMA to transfer 32 bits at a time
	// should be used with the RegDMA1CntH register
	Transfer32 = 0x0400

	// StartNow sest the DMA transfer to start immediately
	// should be used with the RegDMA1CntH register
	StartNow = 0x0000

	// StartVBlank sets the DMA transfer to start on the next vertical blank
	// should be used with the RegDMA1CntH register
	StartVBlank = 0x1000

	// StartHBlank sets the DMA transfer to start horizontal blank
	// should be used with the RegDMA1CntH register
	StartHBlank = 0x2000

	// StartFIFO sest the DMA transfer to start when the configured FIFO buffer is empty
	// should be used with the RegDMA1CntH register
	StartFIFO = 0x3000

	// IRQEnable turns on interrupts when the DMA tranfer is complete
	// should be used with the RegDMA1CntH register
	IRQEnable = 0x4000

	// IRQDisable turns on interrupts when the DMA transfer is complete
	// should be used with the RegDMA1CntH register
	IRQDisable = 0x0000

	// start the DMA transfer
	// should be used with the RegDMA1CntH register
	DMAOn = 0x8000

	// stop the DMA transfer
	// should be used with the RegDMA1CntH register
	DMAOff = 0x0000
)
