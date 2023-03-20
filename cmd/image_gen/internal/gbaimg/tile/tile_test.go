package tile

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg"
	"github.com/go-test/deep"
)

var (
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	black = color.RGBA{0x00, 0x00, 0x00, 0xFF}
	red   = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	blue  = color.RGBA{0x00, 0x00, 0xFF, 0xff}
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
	img := gbaimg.NewRGB16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, pixels[y*width+x])
		}
	}

	return img
}

func TestNewMeta(t *testing.T) {
	img8x8 := newImage16(8, 8, []color.Color{
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
	})

	img16x16 := newImage16(16, 16, []color.Color{
		red, white, red, white, red, white, red, white, white, blue, white, blue, white, blue, white, blue,
		white, red, white, red, white, red, white, red, blue, white, blue, white, blue, white, blue, white,
		red, white, red, white, red, white, red, white, white, blue, white, blue, white, blue, white, blue,
		white, red, white, red, white, red, white, red, blue, white, blue, white, blue, white, blue, white,
		red, white, red, white, red, white, red, white, white, blue, white, blue, white, blue, white, blue,
		white, red, white, red, white, red, white, red, blue, white, blue, white, blue, white, blue, white,
		red, white, red, white, red, white, red, white, white, blue, white, blue, white, blue, white, blue,
		white, red, white, red, white, red, white, red, blue, white, blue, white, blue, white, blue, white,

		green, black, green, black, green, black, green, black, white, black, white, black, white, black, white, black,
		black, green, black, green, black, green, black, green, black, white, black, white, black, white, black, white,
		green, black, green, black, green, black, green, black, white, black, white, black, white, black, white, black,
		black, green, black, green, black, green, black, green, black, white, black, white, black, white, black, white,
		green, black, green, black, green, black, green, black, white, black, white, black, white, black, white, black,
		black, green, black, green, black, green, black, green, black, white, black, white, black, white, black, white,
		green, black, green, black, green, black, green, black, white, black, white, black, white, black, white, black,
		black, green, black, green, black, green, black, green, black, white, black, white, black, white, black, white,
	})

	type args struct {
		img  image.Image
		pal  color.Palette
		size Size
	}
	tests := []struct {
		name string
		args args
		want *Meta
	}{
		{
			"8x8 meta tile",
			args{
				img:  img8x8,
				pal:  color.Palette{red, white},
				size: S8x8,
			},
			&Meta{
				Size: S8x8,
				Img:  img8x8,
				Pal:  color.Palette{red, white},
				Tiles: []image.Image{newImage(8, 8, []color.Color{
					red, white, red, white, red, white, red, white,
					white, red, white, red, white, red, white, red,
					red, white, red, white, red, white, red, white,
					white, red, white, red, white, red, white, red,
					red, white, red, white, red, white, red, white,
					white, red, white, red, white, red, white, red,
					red, white, red, white, red, white, red, white,
					white, red, white, red, white, red, white, red,
				})},
			},
		},
		{
			"16x16 meta tile",
			args{
				img:  img16x16,
				pal:  color.Palette{red, white, green, blue, black},
				size: S16x16,
			},
			&Meta{
				Size: S16x16,
				Img:  img16x16,
				Pal:  color.Palette{red, white, green, blue, black},
				Tiles: []image.Image{
					newImage(8, 8, []color.Color{
						red, white, red, white, red, white, red, white,
						white, red, white, red, white, red, white, red,
						red, white, red, white, red, white, red, white,
						white, red, white, red, white, red, white, red,
						red, white, red, white, red, white, red, white,
						white, red, white, red, white, red, white, red,
						red, white, red, white, red, white, red, white,
						white, red, white, red, white, red, white, red,
					}),
					newImage(8, 8, []color.Color{
						white, blue, white, blue, white, blue, white, blue,
						blue, white, blue, white, blue, white, blue, white,
						white, blue, white, blue, white, blue, white, blue,
						blue, white, blue, white, blue, white, blue, white,
						white, blue, white, blue, white, blue, white, blue,
						blue, white, blue, white, blue, white, blue, white,
						white, blue, white, blue, white, blue, white, blue,
						blue, white, blue, white, blue, white, blue, white,
					}),
					newImage(8, 8, []color.Color{
						green, black, green, black, green, black, green, black,
						black, green, black, green, black, green, black, green,
						green, black, green, black, green, black, green, black,
						black, green, black, green, black, green, black, green,
						green, black, green, black, green, black, green, black,
						black, green, black, green, black, green, black, green,
						green, black, green, black, green, black, green, black,
						black, green, black, green, black, green, black, green,
					}),
					newImage(8, 8, []color.Color{
						white, black, white, black, white, black, white, black,
						black, white, black, white, black, white, black, white,
						white, black, white, black, white, black, white, black,
						black, white, black, white, black, white, black, white,
						white, black, white, black, white, black, white, black,
						black, white, black, white, black, white, black, white,
						white, black, white, black, white, black, white, black,
						black, white, black, white, black, white, black, white,
					}),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMeta(tt.args.img, tt.args.pal, tt.args.size)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("NewMeta() diffs(%d): \n%s", len(diff), strings.Join(diff, "\n"))
			}
		})
	}
}

func TestUnique(t *testing.T) {
	pal := color.Palette{red, green, blue, white, black}

	imgA := newImage16(16, 16, []color.Color{
		red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white,
		green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red,
		blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green,
		white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue,
		red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white,
		green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red,
		blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green,
		white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue,
		red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white,
		green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red,
		blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green,
		white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue,
		red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white,
		green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red,
		blue, white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green,
		white, red, green, blue, white, red, green, blue, white, red, green, blue, white, red, green, blue,
	})

	imgB := newImage16(16, 16, []color.Color{
		red, white, black, white, black, white, black, white, black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black, white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white, black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black, white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white, black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black, white, black, white, black, white, black, white, black,
		black, white, black, white, black, white, black, white, black, white, black, white, black, white, black, white,
		white, black, white, black, white, black, white, black, white, black, white, black, white, black, white, black,
		red, blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue,
		blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue, red,
		red, blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue,
		blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue, red,
		red, blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue,
		blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue, red,
		red, blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue,
		blue, red, blue, red, blue, red, blue, green, blue, green, blue, green, blue, green, blue, red,
	})

	type args struct {
		tiles []*Meta
	}
	tests := []struct {
		name string
		args args
		want []*Meta
	}{
		{
			"duplicate tiles",
			args{
				tiles: []*Meta{
					NewMeta(gbaimg.Flip(imgA, true, false), pal, S8x8),
					NewMeta(gbaimg.Flip(imgA, false, true), pal, S8x8),
					NewMeta(gbaimg.Flip(imgA, true, true), pal, S8x8),
					NewMeta(imgA, pal, S8x8),
					NewMeta(gbaimg.Flip(imgB, true, false), pal, S8x8),
					NewMeta(gbaimg.Flip(imgB, false, true), pal, S8x8),
					NewMeta(gbaimg.Flip(imgB, true, true), pal, S8x8),
					NewMeta(imgB, pal, S8x8),
				},
			},
			[]*Meta{
				NewMeta(imgB, pal, S8x8),
				NewMeta(imgA, pal, S8x8),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.args.tiles)
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("Unique() diffs(%d): %s\n", len(diff), strings.Join(diff, "\n"))
			}
		})
	}
}

func TestMeta_Hash(t *testing.T) {
	img := newImage16(8, 8, []color.Color{
		white, blue, white, blue, red, green, red, green,
		blue, white, blue, white, green, red, green, red,
		white, blue, white, blue, red, green, red, green,
		blue, white, blue, white, green, red, green, red,
		white, black, white, black, white, blue, white, blue,
		black, white, black, white, blue, white, blue, white,
		white, black, white, black, white, blue, white, blue,
		black, white, black, white, blue, white, blue, white,
	})
	pal := color.Palette{red, green, blue, white, black}

	type fields struct {
		Size  Size
		Img   image.Image
		Pal   color.Palette
		Tiles []image.Image
	}
	tests := []struct {
		name   string
		fields fields
		want   [md5.Size]byte
	}{
		{
			"horizontal flipped",
			fields{
				Size:  S8x8,
				Img:   gbaimg.Flip(img, true, false),
				Pal:   pal,
				Tiles: []image.Image{img},
			},
			NewMeta(img, pal, S8x8).Hash(),
		},
		{
			"vertical flipped",
			fields{
				Size:  S8x8,
				Img:   gbaimg.Flip(img, false, true),
				Pal:   pal,
				Tiles: []image.Image{img},
			},
			NewMeta(img, pal, S8x8).Hash(),
		},
		{
			"horizontal and vertical flipped",
			fields{
				Size:  S8x8,
				Img:   gbaimg.Flip(img, true, true),
				Pal:   pal,
				Tiles: []image.Image{img},
			},
			NewMeta(img, pal, S8x8).Hash(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meta{
				Size:  tt.fields.Size,
				Img:   tt.fields.Img,
				Pal:   tt.fields.Pal,
				Tiles: tt.fields.Tiles,
			}
			if got := m.Hash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Meta.Hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMeta_Bytes(t *testing.T) {
	img := newImage16(8, 8, []color.Color{
		white, blue, white, blue, red, green, red, green,
		blue, white, blue, white, green, red, green, red,
		white, blue, white, blue, red, green, red, green,
		blue, white, blue, white, green, red, green, red,

		white, black, white, black, white, blue, white, blue,
		black, white, black, white, blue, white, blue, white,
		white, black, white, black, white, blue, white, blue,
		black, white, black, white, blue, white, blue, white,
	})
	pal := color.Palette{red, green, blue, white, black}

	type fields struct {
		Size  Size
		Img   image.Image
		Pal   color.Palette
		Tiles []image.Image
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			"8x8 meta tile",
			fields{
				Size:  S8x8,
				Img:   img,
				Pal:   pal,
				Tiles: []image.Image{img},
			},
			[]byte{
				0x23, 0x23, 0x10, 0x10,
				0x32, 0x32, 0x01, 0x01,
				0x23, 0x23, 0x10, 0x10,
				0x32, 0x32, 0x01, 0x01,
				0x43, 0x43, 0x23, 0x23,
				0x34, 0x34, 0x32, 0x32,
				0x43, 0x43, 0x23, 0x23,
				0x34, 0x34, 0x32, 0x32,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meta{
				Size:  tt.fields.Size,
				Img:   tt.fields.Img,
				Pal:   tt.fields.Pal,
				Tiles: tt.fields.Tiles,
			}
			if got := m.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Meta.Bytes() diff \n%s", hexDiff(got, tt.want))
			}
		})
	}
}

func hexDiff(a, b []byte) string {
	max := len(a)
	if len(b) > max {
		max = len(b)
	}

	var aStr, bStr, diffStr []string
	for i := 0; i < max; i++ {
		var aCheck, bCheck *byte
		if len(a) > i {
			aStr = append(aStr, fmt.Sprintf("0x%02X", a[i]))
			aCheck = &a[i]
		}
		if len(b) > i {
			bStr = append(bStr, fmt.Sprintf("0x%02X", b[i]))
			bCheck = &b[i]
		}

		var diff string
		switch {
		case aCheck == nil && bCheck != nil:
			diff = "[!!] "
		case aCheck != nil && bCheck == nil:
			diff = "[!!] "
		case *aCheck != *bCheck:
			diff = "[!!] "
		default:
			diff = "-----"
		}

		diffStr = append(diffStr, diff)
	}

	return fmt.Sprintf("%s\n%s\n%s",
		"  "+strings.Join(diffStr, " ")+"  ",
		"[ "+strings.Join(aStr, ", ")+" ]",
		"[ "+strings.Join(bStr, ", ")+" ]",
	)
}

func TestMeta_IsTransparent(t *testing.T) {
	type fields struct {
		Size  Size
		Img   image.Image
		Pal   color.Palette
		Tiles []image.Image
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			"transparent tile",
			fields{
				Size: S8x8,
				Img: newImage16(8, 8, []color.Color{
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
				}),
				Pal: color.Palette{white},
			},
			true,
		},
		{
			"sprite tile",
			fields{
				Size: S16x8,
				Img: newImage16(16, 8, []color.Color{
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white, white, white, white, white, white, white, white, white,
				}),
				Pal: color.Palette{white},
			},
			false,
		},
		{
			"colored tile",
			fields{
				Size: S8x8,
				Img: newImage16(8, 8, []color.Color{
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, black, black, white, white, white,
					white, white, white, black, black, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
					white, white, white, white, white, white, white, white,
				}),
				Pal: color.Palette{white, black},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Meta{
				Size:  tt.fields.Size,
				Img:   tt.fields.Img,
				Pal:   tt.fields.Pal,
				Tiles: tt.fields.Tiles,
			}
			if got := m.IsTransparent(); got != tt.want {
				t.Errorf("Meta.Transparent() = %v, want %v", got, tt.want)
			}
		})
	}
}
