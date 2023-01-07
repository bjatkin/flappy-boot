package timer

import (
	"unsafe"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

type Register uint

const (
	// RegTM0Count is the controll register that sets the <reload> value for timer 0. Setting this register does not changes the
	// current counter value. Rather, this value is loaded into timer 0 register when the timer starts or overflows
	// It is important to note reading from this register does NOT return the <reload> value.
	// instead it returns the current counter value (or recent/frozen counter value if timer 0 has been stopped)
	RegTM0Count Register = 0x0400_0100

	// RegTM0Cnt is the controll register used to controll timer 0. It can also be used to start or stop
	// the timer. It has the following layout
	//
	// [0 - 1] Increment Frequency - modifies how often the Timer0 ticks
	//   - Freq1 - sets the timer to increment once every cpu cycle (55.59 ns)
	//   - Freq64 - sets the timer to increment once every 64 cpu cycles (3.815 μs)
	//   - Freq256 - sets the timer to increment once every 256 cpu cycles (15.26 μs)
	//   - Freq1024 - sett the timer to increment once every 1024 cpu cycles (61.04 μs)
	//
	// [6] Timer Interrupt - can be used to enable timer interrupts when the timer value overflows
	//   - IRQEnable - enables hardware intrrupts when timer 0 overflows
	//   - IRQDisable - disables hardware interrupts when timer 0 overflows
	//
	// [7] Timer Start - used to start and stop the timer
	//   - TmrStart - start timer 0. The value of RegTM0CntL will be used as the starting point
	//   - TmrStop - stop/ freeze timmer 0
	RegTM0Cnt Register = 0x0400_0102
)

// CounterReg is the type used for all the timer counters, see Counter# for more information on using this type
type CounterReg uint16

var (
	// Counter0 is the controll register that sets the <reload> value for timer 0. Setting this register does not changes the
	// current counter value. Rather, this value is loaded into timer 0 register when the timer starts or overflows
	// It is important to note reading from this register does NOT return the <reload> value.
	// instead it returns the current counter value (or recent/frozen counter value if timer 0 has been stopped)
	Counter0 = (*CounterReg)(unsafe.Pointer(memmap.IOAddr + 0x0100))

	// Counter1 is the controll register that sets the <reload> value for timer 1. Setting this register does not changes the
	// current counter value. Rather, this value is loaded into timer 1 register when the timer starts or overflows
	// It is important to note reading from this register does NOT return the <reload> value.
	// instead it returns the current counter value (or recent/frozen counter value if timer 1 has been stopped)
	Counter1 = (*CounterReg)(unsafe.Pointer(memmap.IOAddr + 0x0104))

	// Counter2 is the controll register that sets the <reload> value for timer 2. Setting this register does not changes the
	// current counter value. Rather, this value is loaded into timer 2 register when the timer starts or overflows
	// It is important to note reading from this register does NOT return the <reload> value.
	// instead it returns the current counter value (or recent/frozen counter value if timer 2 has been stopped)
	Counter2 = (*CounterReg)(unsafe.Pointer(memmap.IOAddr + 0x0108))

	// Counter3 is the controll register that sets the <reload> value for timer 3. Setting this register does not changes the
	// current counter value. Rather, this value is loaded into timer 3 register when the timer starts or overflows
	// It is important to note reading from this register does NOT return the <reload> value.
	// instead it returns the current counter value (or recent/frozen counter value if timer 3 has been stopped)
	Counter3 = (*CounterReg)(unsafe.Pointer(memmap.IOAddr + 0x010C))
)

// ControllReg is the type used for the timer controll register, see Controll# for mor information on using this type
type ControllReg uint16

var (
	// Controll0 is the controll register used to controll timer 0. It can also be used to start or stop
	// the timer. It has the following layout
	//
	// [0 - 1] Increment Frequency - modifies how often the Timer0 ticks
	//   - Freq1 - sets the timer to increment once every cpu cycle (55.59 ns)
	//   - Freq64 - sets the timer to increment once every 64 cpu cycles (3.815 μs)
	//   - Freq256 - sets the timer to increment once every 256 cpu cycles (15.26 μs)
	//   - Freq1024 - sett the timer to increment once every 1024 cpu cycles (61.04 μs)
	//
	// [2] Count-up Timing - This can not be using with this timer
	//
	// [6] Timer Interrupt - can be used to enable timer interrupts when the timer value overflows
	//   - IRQEnable - enables hardware intrrupts when timer 0 overflows
	//   - IRQDisable - disables hardware interrupts when timer 0 overflows
	//
	// [7] Timer Start - used to start and stop the timer
	//   - TmrStart - start timer 0. The value of Counter0 will be used as the starting point
	//   - TmrStop - stop/ freeze timmer 0
	Controll0 = (*ControllReg)(unsafe.Pointer(memmap.IOAddr + 0x0102))

	// Controll1 is the controll register used to controll timer 1. It can also be used to start or stop
	// the timer. It has the following layout
	//
	// [0 - 1] Increment Frequency - modifies how often the Timer0 ticks
	//   - Freq1 - sets the timer to increment once every cpu cycle (55.59 ns)
	//   - Freq64 - sets the timer to increment once every 64 cpu cycles (3.815 μs)
	//   - Freq256 - sets the timer to increment once every 256 cpu cycles (15.26 μs)
	//   - Freq1024 - sett the timer to increment once every 1024 cpu cycles (61.04 μs)
	//
	// [2] Count-up Timing - Ignore the Increment Frequency setting an instead increment the timer when
	//     timer 0 overflows
	//   - CountUpEnable - Enables the count up overflow
	//
	// [6] Timer Interrupt - can be used to enable timer interrupts when the timer value overflows
	//   - IRQEnable - enables hardware intrrupts when timer 0 overflows
	//   - IRQDisable - disables hardware interrupts when timer 0 overflows
	//
	// [7] Timer Start - used to start and stop the timer
	//   - TmrStart - start timer 1. The value of Counter1 will be used as the starting point
	//   - TmrStop - stop/ freeze timmer 1
	Controll1 = (*ControllReg)(unsafe.Pointer(memmap.IOAddr + 0x0106))

	// Controll2 is the controll register used to controll timer 1. It can also be used to start or stop
	// the timer. It has the following layout
	//
	// [0 - 1] Increment Frequency - modifies how often the Timer0 ticks
	//   - Freq1 - sets the timer to increment once every cpu cycle (55.59 ns)
	//   - Freq64 - sets the timer to increment once every 64 cpu cycles (3.815 μs)
	//   - Freq256 - sets the timer to increment once every 256 cpu cycles (15.26 μs)
	//   - Freq1024 - sett the timer to increment once every 1024 cpu cycles (61.04 μs)
	//
	// [2] Count-up Timing - Ignore the Increment Frequency setting an instead increment the timer when
	//     timer 1 overflows
	//   - CountUpEnable - Enables the count up overflow
	//
	// [6] Timer Interrupt - can be used to enable timer interrupts when the timer value overflows
	//   - IRQEnable - enables hardware intrrupts when timer 2 overflows
	//   - IRQDisable - disables hardware interrupts when timer 2 overflows
	//
	// [7] Timer Start - used to start and stop the timer
	//   - TmrStart - start timer 2. The value of Counter2 will be used as the starting point
	//   - TmrStop - stop/ freeze timmer 2
	Controll2 = (*ControllReg)(unsafe.Pointer(memmap.IOAddr + 0x010A))

	// Controll3 is the controll register used to controll timer 1. It can also be used to start or stop
	// the timer. It has the following layout
	//
	// [0 - 1] Increment Frequency - modifies how often the Timer0 ticks
	//   - Freq1 - sets the timer to increment once every cpu cycle (55.59 ns)
	//   - Freq64 - sets the timer to increment once every 64 cpu cycles (3.815 μs)
	//   - Freq256 - sets the timer to increment once every 256 cpu cycles (15.26 μs)
	//   - Freq1024 - sett the timer to increment once every 1024 cpu cycles (61.04 μs)
	//
	// [2] Count-up Timing - Ignore the Increment Frequency setting an instead increment the timer when
	//     timer 2 overflows
	//   - CountUpEnable - Enables the count up overflow
	//
	// [6] Timer Interrupt - can be used to enable timer interrupts when the timer value overflows
	//   - IRQEnable - enables hardware intrrupts when timer 3 overflows
	//   - IRQDisable - disables hardware interrupts when timer 3 overflows
	//
	// [7] Timer Start - used to start and stop the timer
	//   - TimerStart - start timer 3. The value of Counter3 will be used as the starting point
	//   - TimerStop - stop/ freeze timer 3
	Controll3 = (*ControllReg)(unsafe.Pointer(memmap.IOAddr + 0x010A))
)

const (
	// Freq1 is the default Increment Frequency. It sets the timer to increment every CPU cycles
	//
	// Cycles: 1
	// Frequency: 16.78 MHz
	// Period: 55.59 ns
	Freq1 ControllReg = 0x0000

	// Freq64 sets the timer to increment every 64 CPU cycles
	//
	// Cycles: 64
	// Frequency: 262.21 kHz
	// Period: 3.815 μs
	Freq64 ControllReg = 0x0001

	// Freq256 sets the timer to increment every 256 CPU cycles
	//
	// Cycles: 256
	// Frequency: 65.536 kHz
	// Period: 15.26 μs
	Freq256 ControllReg = 0x0002

	// Freq1024 sets the timer to increment every 1024 CPU cycles
	//
	// Cycles: 1024
	// Frequency: 16.384 kHz
	// Period: 61.04 μs
	Freq1024 ControllReg = 0x0003

	// CountUpEnable sets the timer to ignore the Increment Frequency and instead increment when the
	// previous timer overflows.
	//
	// NOTE: this can not be used with timer 0 as that is the first timer
	CountUpEnable ControllReg = 0x0004

	// IRQEnable enables timer hardware interrupts. The interrupt will be triggered when the timer overflows
	IRQEnable ControllReg = 0x0040

	// IRQDisable disables timer hardware interrupts.
	IRQDisable ControllReg = 0x0000

	// TimerStart starts the timer ticking
	TimerStart ControllReg = 0x0080

	// TimerStop stops/ freezes the timer from ticking
	TimerStop ControllReg = 0x0000
)
