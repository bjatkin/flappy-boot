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

	dirty bool
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
			p.dirty = true

			// TODO: this could be come a source of lots of garbage, it should be cleaned up
			return &PMem{
				Memory: p.memory[i*memmap.PaletteOffset : (i+1)*memmap.PaletteOffset],
				Offset: i,
			}, nil
		}
	}

	return nil, ErrOOM
}

// Free marks the memory associated with the provided allocation as free
func (p *Pal) Free(mem *PMem) {
	p.meta[mem.Offset] = false
}

// IsDirty returns true if the allocator has made any new allocations since the palette was last marked clean
func (p *Pal) IsDirty() bool {
	return p.dirty
}

// MarkClean marks the allocator as clean
func (p *Pal) MarkClean() {
	p.dirty = false
}
