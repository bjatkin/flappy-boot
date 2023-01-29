package gbaimg

import (
	"image"
	"image/color"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// RGB15Model is a color model that can be used to covert any color.Color into
// the gbacol.RGB15 color space
var RGB15Model = color.ModelFunc(rgb15Model)

func rgb15Model(c color.Color) color.Color {
	if _, ok := c.(gbacol.RGB15); ok {
		return c
	}
	return gbacol.NewRGB15(c)
}

// Walk walks the provided image from top, left to bottom, right
// the provided function fn will be called at each pixel
func Walk(img image.Image, fn func(x, y int)) {
	WalkN(img, image.Point{X: 1, Y: 1}, fn)
}

// WalkN walks the provided image from top, left to bottom, right
// at each horizontal step stride.X pixles are skipped and at each vertical step
// stride.Y pixels are skipped, thus the function fn will be called for every
// strided pixel in the image
func WalkN(img image.Image, stride image.Point, fn func(x, y int)) {
	minX, minY := img.Bounds().Min.X, img.Bounds().Min.Y
	dx, dy := img.Bounds().Dx(), img.Bounds().Dy()

	for y := minY; y < dy+minY; y += stride.Y {
		for x := minX; x < dx+minX; x += stride.X {
			fn(x, y)
		}
	}
}

// MutImage is a mutable image
// It is the same as the image.Image interface but it also includes the Set(x, y int, c color.Color) method
type MutImage interface {
	image.Image
	Set(x, y int, c color.Color)
}

// Copy coppies an image from the src to the destination
// Due to the fact that dest is being modified it must be a MutImage
func Copy(src image.Image, dest MutImage) {
	destMin := dest.Bounds().Min
	srcMin := src.Bounds().Min
	Walk(src, func(x, y int) {
		dest.Set((x-srcMin.X)+destMin.X, (y-srcMin.Y)+destMin.Y, src.At(x, y))
	})
}

// SubImage can be used to get a copy of an subsction of an image.
// the original image and the returned sub image are not realated and do not
// share any underlying data
func SubImage(img image.Image, r image.Rectangle) image.Image {
	min := r.Bounds().Min
	cpy := image.NewRGBA(image.Rect(0, 0, r.Bounds().Dx(), r.Bounds().Dy()))

	Walk(cpy, func(x, y int) {
		cpy.Set(x, y, img.At(min.X+x, min.Y+y))
	})

	return cpy
}

// Flip returns a fliped version of the provided image
// the returned image is a copy and does not share any data with the original image
// the image can be flipped horizontally, vertically or both
func Flip(img image.Image, h, v bool) image.Image {
	cpy := image.NewRGBA(img.Bounds())
	Copy(img, cpy)

	switch {
	case h && v:
		Walk(img, func(x, y int) {
			newX := img.Bounds().Max.X - x - 1
			newY := img.Bounds().Max.Y - y - 1
			cpy.Set(newX, newY, img.At(x, y))
		})
	case h:
		Walk(img, func(x, y int) {
			newX := img.Bounds().Max.X - x - 1
			cpy.Set(newX, y, img.At(x, y))
		})
	case v:
		Walk(img, func(x, y int) {
			newY := img.Bounds().Max.Y - y - 1
			cpy.Set(x, newY, img.At(x, y))
		})
	}

	return cpy
}

// NewPal creates a new color.Palette, colors are converted into the gbacol.RGB15 color space
// when the palette is created so the resulting palette may have a reduced color range
func NewPal(m image.Image) color.Palette {
	colorMap := make(map[gbacol.RGB15]struct{}, 256)
	var transparent gbacol.RGB15
	Walk(m, func(x, y int) {
		c := m.At(x, y)
		rgb15 := rgb15Model(c).(gbacol.RGB15)
		if x == 0 && y == 0 {
			transparent = rgb15
		}
		colorMap[rgb15] = struct{}{}
	})

	pal := maps.Keys(colorMap)
	for i, color := range pal {
		if color == transparent {
			pal[0], pal[i] = pal[i], pal[0]
		}
	}
	slices.Sort(pal[1:])

	var ret color.Palette
	for _, color := range pal {
		ret = append(ret, color)
	}

	return ret
}

// Match returns true if the color data of the underlying images is equivilant
func Match(a, b image.Image) bool {
	if a.Bounds().Dx() != b.Bounds().Dx() || a.Bounds().Dy() != b.Bounds().Dy() {
		return false
	}

	match := true
	Walk(a, func(x, y int) {
		c1 := a.At(x, y)
		c2 := b.At(x, y)
		rgb1 := color.RGBAModel.Convert(c1).(color.RGBA)
		rgb2 := color.RGBAModel.Convert(c2).(color.RGBA)

		match = match && rgb1.R == rgb2.R && rgb1.G == rgb2.G && rgb1.B == rgb2.B
	})

	return match
}
