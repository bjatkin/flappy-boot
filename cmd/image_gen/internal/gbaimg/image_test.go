package gbaimg

import (
	"image"
	"image/color"
	"reflect"
	"testing"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
)

var (
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	black = color.RGBA{}
	red   = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	blue  = color.RGBA{0x00, 0x00, 0xFF, 0xff}

	black16 = gbacol.NewRGB15(black)
	red16   = gbacol.NewRGB15(red)
	green16 = gbacol.NewRGB15(green)
	blue16  = gbacol.NewRGB15(blue)
)

func newImage(width, height int, pixels []color.Color) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, pixels[y*width+x])
		}
	}

	return img
}

func newImage16(width, height int, pixels []color.Color) image.Image {
	img := NewRGB16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, pixels[y*width+x])
		}
	}

	return img
}

func TestSubImage(t *testing.T) {
	type args struct {
		img image.Image
		r   image.Rectangle
	}
	tests := []struct {
		name string
		args args
		want image.Image
	}{
		{
			"2x2",
			args{
				img: newImage(4, 4, []color.Color{
					white, black, black, white,
					black, red, green, black,
					black, blue, white, black,
					white, black, black, white,
				}),
				r: image.Rect(1, 1, 3, 3),
			},
			newImage(2, 2, []color.Color{
				red, green,
				blue, white,
			}),
		},
		{
			"1x1",
			args{
				img: newImage(3, 3, []color.Color{
					red, blue, red,
					blue, white, blue,
					red, blue, red,
				}),
				r: image.Rect(1, 1, 2, 2),
			},
			newImage(1, 1, []color.Color{white}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SubImage(tt.args.img, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SubImage() = \n%#v, want \n%#v", got, tt.want)
			}
		})
	}
}

func TestFlip(t *testing.T) {
	type args struct {
		img image.Image
		h   bool
		v   bool
	}
	tests := []struct {
		name string
		args args
		want image.Image
	}{
		{
			"2x2 horizontal",
			args{
				img: newImage(2, 2, []color.Color{
					red, green,
					blue, white,
				}),
				h: true,
				v: false,
			},
			newImage(2, 2, []color.Color{
				green, red,
				white, blue,
			}),
		},
		{
			"2x2 vertical",
			args{
				img: newImage(2, 2, []color.Color{
					red, green,
					blue, white,
				}),
				h: false,
				v: true,
			},
			newImage(2, 2, []color.Color{
				blue, white,
				red, green,
			}),
		},
		{
			"2x2 horizontal and vertical",
			args{
				img: newImage(2, 2, []color.Color{
					red, green,
					blue, white,
				}),
				h: true,
				v: true,
			},
			newImage(2, 2, []color.Color{
				white, blue,
				green, red,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Flip(tt.args.img, tt.args.h, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCopy(t *testing.T) {
	type args struct {
		src  image.Image
		dest MutImage
	}
	tests := []struct {
		name string
		args args
		want image.Image
	}{
		{
			"2x2",
			args{
				src: newImage(2, 2, []color.Color{
					red, green,
					blue, white,
				}),
				dest: newImage(2, 2, []color.Color{
					black, black,
					black, black,
				}).(*image.RGBA),
			},
			newImage(2, 2, []color.Color{
				red, green,
				blue, white,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Copy(tt.args.src, tt.args.dest)

			if !reflect.DeepEqual(tt.args.dest, tt.want) {
				t.Errorf("TestCopy() = %v, want %v", tt.args.dest, tt.want)
			}
		})
	}
}

func TestNewPal(t *testing.T) {
	type args struct {
		m image.Image
	}
	tests := []struct {
		name string
		args args
		want color.Palette
	}{
		{
			"small palette",
			args{
				m: newImage(2, 2, []color.Color{
					red, black,
					blue, green,
				}),
			},
			color.Palette{
				red16,
				black16,
				green16,
				blue16,
			},
		},
		{
			"16 color palette",
			args{
				m: newImage16(4, 4, []color.Color{
					gbacol.RGB15(0xFF),
					gbacol.RGB15(0x70),
					gbacol.RGB15(0xE0),
					gbacol.RGB15(0x00),
					gbacol.RGB15(0x60),
					gbacol.RGB15(0x20),
					gbacol.RGB15(0x30),
					gbacol.RGB15(0x10),
					gbacol.RGB15(0x40),
					gbacol.RGB15(0x90),
					gbacol.RGB15(0xA0),
					gbacol.RGB15(0xB0),
					gbacol.RGB15(0xF0),
					gbacol.RGB15(0xC0),
					gbacol.RGB15(0xD0),
					gbacol.RGB15(0x50),
				}),
			},
			color.Palette{
				gbacol.RGB15(0xFF),
				gbacol.RGB15(0x00),
				gbacol.RGB15(0x10),
				gbacol.RGB15(0x20),
				gbacol.RGB15(0x30),
				gbacol.RGB15(0x40),
				gbacol.RGB15(0x50),
				gbacol.RGB15(0x60),
				gbacol.RGB15(0x70),
				gbacol.RGB15(0x90),
				gbacol.RGB15(0xA0),
				gbacol.RGB15(0xB0),
				gbacol.RGB15(0xC0),
				gbacol.RGB15(0xD0),
				gbacol.RGB15(0xE0),
				gbacol.RGB15(0xF0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPal(tt.args.m); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPal() = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func TestMatch(t *testing.T) {
	type args struct {
		a image.Image
		b image.Image
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"2x2 match",
			args{
				a: newImage(2, 2, []color.Color{
					white, black,
					red, blue,
				}),
				b: newImage(2, 2, []color.Color{
					white, black,
					red, blue,
				}),
			},
			true,
		},
		{
			"size miss match",
			args{
				a: newImage(2, 2, []color.Color{
					white, black,
					red, blue,
				}),
				b: newImage(3, 3, []color.Color{
					white, black, red,
					red, blue, black,
					blue, black, white,
				}),
			},
			false,
		},
		{
			"2x2 color miss match",
			args{
				a: newImage(2, 2, []color.Color{
					white, black,
					red, blue,
				}),
				b: newImage(2, 2, []color.Color{
					white, black,
					red, red,
				}),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Match(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
