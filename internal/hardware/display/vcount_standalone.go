//go:build standalone

package display

var count int

// GetVCount returns the curret vertical scan line being drawn by the PPU
func GetVCount() int {
	// fake out the PPU drawing so we can get past VSyncWait
	count++
	count %= 165
	return count
}
