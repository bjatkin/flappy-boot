package gbaimg

import (
	"image"
	"image/color"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
)

// RGB16 is an image that uses 15 bits per rgb color
// in order to support data aligment, these colors are expanded by 1 bit to be 16 bits each
type RGB16 struct {
	// Pix holds the image's pixels as RGB16 colors. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*2].
	Pix []gbacol.RGB15
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// NewRGB16 creates a new, empty RGB16 image which the dimentions defined by
// the image.Rectangle r
func NewRGB16(r image.Rectangle) *RGB16 {
	return &RGB16{
		Pix:    make([]gbacol.RGB15, r.Dx()*r.Dy()),
		Stride: r.Dx(),
		Rect:   r,
	}
}

// ColorModel can be used to convert any color into the RGB16 color space
func (p *RGB16) ColorModel() color.Model { return RGB15Model }

// Bounds returns the Rect of the RGB16 image
func (p *RGB16) Bounds() image.Rectangle { return p.Rect }

// At returns the color at coordinates x,y in the image
func (p *RGB16) At(x, y int) color.Color {
	return p.RGB15At(x, y)
}

// RGB15At returns the gbacol.RGB15 color at coordinates x, y in the image
func (p *RGB16) RGB15At(x, y int) gbacol.RGB15 {
	if !(image.Point{x, y}).In(p.Rect) {
		return gbacol.RGB15(0)
	}
	i := p.PixOffset(x, y)
	return p.Pix[i]
}

// PixOffset returns the index into Pix that map to the coordinates x, y
func (p *RGB16) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x - p.Rect.Min.X)
}

// Set sets the color in the image at coordinates x, y in the image
func (p *RGB16) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = RGB15Model.Convert(c).(gbacol.RGB15)
}

// SetRGB15 sets the gbacol.RGB15 color at coordinates x, y in the image
func (p *RGB16) SetRGB15(x, y int, c gbacol.RGB15) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	p.Pix[i] = c
}

// SubImage returns a new RGB16 image that uses the provided image.Rectangle r as it's Rect value
// and shares underlying pixel data with the calling image
func (p *RGB16) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB16{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB16{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}
