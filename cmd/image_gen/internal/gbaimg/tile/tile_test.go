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
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"github.com/go-test/deep"
)

var (
	white = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
	black = color.RGBA{}
	red   = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	green = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	blue  = color.RGBA{0x00, 0x00, 0xFF, 0xff}

	white16 = gbacol.NewRGB15(white)
	black16 = gbacol.NewRGB15(black)
	red16   = gbacol.NewRGB15(red)
	green16 = gbacol.NewRGB15(green)
	blue16  = gbacol.NewRGB15(blue)
)

func newImage(width, height int, pixels []color.Color) image.Image {
	img := gbaimg.NewRGB16(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, pixels[y*width+x])
		}
	}

	return img
}

func TestNewMeta(t *testing.T) {
	img8x8 := newImage(8, 8, []color.Color{
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
		red, white, red, white, red, white, red, white,
		white, red, white, red, white, red, white, red,
	})

	img16x16 := newImage(16, 16, []color.Color{
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

	imgA := newImage(16, 16, []color.Color{
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

	imgB := newImage(16, 16, []color.Color{
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
				NewMeta(imgA, pal, S8x8),
				NewMeta(imgB, pal, S8x8),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.args.tiles)
			fmt.Println("LEN: ", len(got))
			if diff := deep.Equal(got, tt.want); len(diff) > 0 {
				t.Errorf("Unique() diffs(%d): %s\n", len(diff), strings.Join(diff, "\n"))
			}
		})
	}
}

func TestMeta_Hash(t *testing.T) {
	img := newImage(8, 8, []color.Color{
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
