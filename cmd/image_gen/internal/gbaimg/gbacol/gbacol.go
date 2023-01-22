package gbacol

import (
	"image/color"
	"math"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/byteconv"
)

// RGB15 is a 15 bit color. It has the same format as color data in the GBA VRAM
type RGB15 uint16

// NewRGB15 creates an RGB15 color from a color.Color
func NewRGB15(c color.Color) RGB15 {
	rgbac := color.RGBAModel.Convert(c)
	r, g, b, _ := rgbac.RGBA()

	conv := func(c uint32) uint16 {
		return uint16(math.Round(float64(c) / float64(0xFFFF) * 0b11111))
	}
	r16 := conv(r)
	g16 := conv(g)
	b16 := conv(b)

	return RGB15((r16) + (g16 << 5) + (b16 << 10))
}

// RGBA returns the colors red, green, blue, and alpha channels.
// each channel is in the range 0-0xFFFF
func (rgb RGB15) RGBA() (r, g, b, a uint32) {
	conv := func(c RGB15) uint32 {
		b := byte(c) & 0b11111
		return uint32(math.Round(float64(b) / float64(0b11111) * float64(0xFFFF)))
	}
	r = conv(rgb)
	g = conv((rgb >> 5))
	b = conv((rgb >> 10))
	a = 0xFFFF
	return r, g, b, a
}

// Bytes converts an RGB15 color into a little endian byte slice
func (rgb RGB15) Bytes() []byte {
	return byteconv.Itoa(uint16(rgb))
}
