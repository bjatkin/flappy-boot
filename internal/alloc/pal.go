package alloc

import "github.com/bjatkin/flappy_boot/internal/hardware/memmap"

// PMem is a section of palette memory
type PMem struct {
	Memory []memmap.PaletteValue
	Offset int
}

// Pal is an allocator that can be used with GBA palette memory
type Pal struct {
	meta   [8]bool
	memory []memmap.PaletteValue
}

// NewPal creates a new Pal allocator from a section of palette memory
func NewPal(memory []memmap.PaletteValue) *Pal {
	return &Pal{
		memory: memory,
	}
}

// Alloc returns a section of palette memory. If there are no more palettes availalbe
// an ErrOOM error will be returned
func (p *Pal) Alloc() (*PMem, error) {
	for i := range p.meta {
		if !p.meta[i] {
			p.meta[i] = true
			return &PMem{
				Memory: p.memory[i*8 : (i+1)*8],
				Offset: i,
			}, nil
		}
	}

	return nil, ErrOOM
}

func (p *Pal) Free(mem *PMem) {
	p.meta[mem.Offset] = false
}
