package alloc

import (
	"errors"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

// ErrOOM is returned if the allocator can't find space for the requested allocation
var ErrOOM = errors.New("out of memory")

const used = 0x70000000

// VMem is a section of VRAM memory
type VMem struct {
	Memory []memmap.VRAMValue
	Offset int
}

// VRAM is an allocator that can be used with GBA vram memory
type VRAM struct {
	meta     []int
	memory   []memmap.VRAMValue
	cellSize int
}

// NewVRAM creates a new VRAM allocator from a secontion of vram memory. cellSize is the minimum chunk of
// memory that can be allocated
func NewVRAM(memory []memmap.VRAMValue, cellSize int) *VRAM {
	size := len(memory) / cellSize
	meta := make([]int, size)
	meta[0] = size

	return &VRAM{
		meta:     meta,
		memory:   memory,
		cellSize: cellSize,
	}
}

// Alloc allocates a section of VRAM memory if the requested size is too large for the current VRAM allocator
// and ErrOOM will be returned
func (v *VRAM) Alloc(size int) (*VMem, error) {
	var i int
	for {
		if i >= len(v.meta) {
			return nil, ErrOOM
		}

		cellSize := (v.meta[i] & ^used)
		// check if the current cell is free
		if !v.isFree(i) {
			i += cellSize
			continue
		}

		diff := cellSize - size
		switch {
		case diff == 0:
			v.meta[i] = used | size
			// TODO: this could become a source of lots of garbage, it should be cleaned up
			return &VMem{
				Memory: v.memory[i*v.cellSize : (i+size)*v.cellSize],
				Offset: i,
			}, nil
		case diff > 0:
			v.meta[i] = used | size
			v.meta[i+size] = diff
			// TODO: this could become a source of lots of garbage, it should be cleaned up
			return &VMem{
				Memory: v.memory[i*v.cellSize : (i+size)*v.cellSize],
				Offset: i,
			}, nil
		default:
			i += cellSize
		}
	}
}

// VRAM free frees a section of VRAM and marks it as available for other users
func (v *VRAM) Free(mem *VMem) {
	// find the previous cell
	prev := mem.Offset - 1
	for ; prev >= 0; prev-- {
		if v.meta[prev] != 0 {
			break
		}
	}
	next := mem.Offset + (v.meta[mem.Offset] & ^used)

	switch {
	case v.isFree(prev) && v.isFree(next):
		v.meta[prev] = v.meta[prev] + (v.meta[mem.Offset] & ^used) + v.meta[next]
		v.meta[mem.Offset] = 0
		v.meta[next] = 0
	case v.isFree(next):
		v.meta[mem.Offset] = (v.meta[mem.Offset] & ^used) + v.meta[next]
		v.meta[next] = 0
	case v.isFree(prev):
		v.meta[prev] = v.meta[prev] + (v.meta[mem.Offset] & ^used)
		v.meta[mem.Offset] = 0
	default:
		v.meta[mem.Offset] = (v.meta[mem.Offset] & ^used)
	}
}

// isFree returns true if the specified cell is currently free
func (v *VRAM) isFree(i int) bool {
	if i < 0 || i >= len(v.meta) {
		return false
	}

	return v.meta[i] < used
}
