package ppu

import (
	"image"
	"reflect"
	"testing"
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
