package game

import (
	"errors"
)

const used = 0x70000000

var ErrOOM = errors.New("out of memory")

type Allocator []int

func NewAllocator(size int) Allocator {
	alloc := make(Allocator, size)
	alloc[0] = size

	return alloc
}

func (a Allocator) Alloc(size int) (int, error) {
	var i int
	for {
		if i >= len(a) {
			return 0, ErrOOM
		}

		cellSize := (a[i] & ^used)
		// check if the current cell is free
		if !a.isFree(i) {
			i += cellSize
			continue
		}

		diff := cellSize - size
		switch {
		case diff == 0:
			a[i] = used | size
			return i, nil
		case diff > 0:
			a[i] = used | size
			a[i+size] = diff
			return i, nil
		default:
			i += cellSize
		}
	}
}

func (a Allocator) Free(addr int) {
	// find the previous cell
	prev := addr - 1
	for ; prev >= 0; prev-- {
		if a[prev] != 0 {
			break
		}
	}
	next := addr + (a[addr] & ^used)

	switch {
	case a.isFree(prev) && a.isFree(next):
		a[prev] = a[prev] + (a[addr] & ^used) + a[next]
		a[addr] = 0
		a[next] = 0
	case a.isFree(next):
		a[addr] = (a[addr] & ^used) + a[next]
		a[next] = 0
	case a.isFree(prev):
		a[prev] = a[prev] + (a[addr] & ^used)
		a[addr] = 0
	default:
		a[addr] = (a[addr] & ^used)
	}
}

func (a Allocator) isFree(i int) bool {
	if i < 0 || i >= len(a) {
		return false
	}

	return a[i] < used
}
