package ppu

import (
	"image"
	"reflect"
	"testing"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

func TestBackground_getTileView(t *testing.T) {
	type fields struct {
		Size v2
	}
	type args struct {
		tile int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   image.Rectangle
	}{
		{
			"zero",
			fields{v2{X: 1, Y: 1}},
			args{tile: 0},
			image.Rect(0, 0, 8, 8),
		},
		{
			"simple",
			fields{
				Size: v2{X: 1, Y: 1},
			},
			args{
				tile: 10,
			},
			image.Rect(80, 0, 88, 8),
		},
		{
			"wide screen",
			fields{
				Size: v2{X: 2, Y: 1},
			},
			args{
				tile: 1024,
			},
			image.Rect(256, 0, 264, 8),
		},
		{
			"tall screen",
			fields{
				Size: v2{X: 1, Y: 2},
			},
			args{
				tile: 1024,
			},
			image.Rect(0, 256, 8, 264),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Background{
				Size: tt.fields.Size,
			}
			if got := b.getTileView(tt.args.tile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Background.getTileView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGB15_Set(t *testing.T) {
	type fields struct {
		colors []memmap.PaletteValue
		width  int
		height int
	}
	type args struct {
		x int
		y int
		c memmap.PaletteValue
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []memmap.PaletteValue
	}{
		{
			"out of bounds",
			fields{
				colors: make([]memmap.PaletteValue, 3*3),
				width:  3,
				height: 3,
			},
			args{
				x: -10,
				y: -10,
				c: 0x00FF,
			},
			[]memmap.PaletteValue{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
		},
		{
			"set color",
			fields{
				colors: make([]memmap.PaletteValue, 3*3),
				width:  3,
				height: 3,
			},
			args{
				x: 1,
				y: 1,
				c: 0x0007,
			},
			[]memmap.PaletteValue{
				0, 0, 0,
				0, 0x007, 0,
				0, 0, 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &RGB15{
				colors: tt.fields.colors,
				width:  tt.fields.width,
				height: tt.fields.height,
			}
			i.Set(tt.args.x, tt.args.y, tt.args.c)

			if !reflect.DeepEqual(i.colors, tt.want) {
				t.Errorf("RGB15.Set() colors do not match:\n%v\n%v", i.colors, tt.want)
			}
		})
	}
}
