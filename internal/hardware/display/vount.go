//go:build !standalone

package display

import "github.com/bjatkin/flappy_boot/internal/hardware/memmap"

// GetVCount returns the curret vertical scan line being drawn by the PPU
func GetVCount() int {
	return int(memmap.GetReg(VCount))
}
