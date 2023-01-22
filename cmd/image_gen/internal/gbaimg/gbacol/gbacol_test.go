package gbacol

import (
	"image/color"
	"reflect"
	"testing"
)

func TestRGB15_Bytes(t *testing.T) {
	tests := []struct {
		name string
		rgb  RGB15
		want []byte
	}{
		{
			"red",
			NewRGB15(color.RGBA{R: 0xFF}),
			[]byte{0b000_11111, 0b0_00000_00},
		},
		{
			"green",
			NewRGB15(color.RGBA{G: 0xFF}),
			[]byte{0b111_00000, 0b0_00000_11},
		},
		{
			"blue",
			NewRGB15(color.RGBA{B: 0xFF}),
			[]byte{0b000_00000, 0b0_11111_00},
		},
		{
			"gray",
			NewRGB15(color.RGBA{R: 0x0A, G: 0x0A, B: 0x0A}),
			[]byte{0b001_00001, 0b0_00001_00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rgb.Bytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RGB15.Bytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGB15_RGBA(t *testing.T) {
	tests := []struct {
		name  string
		rgb   RGB15
		wantR uint32
		wantG uint32
		wantB uint32
		wantA uint32
	}{
		{
			"red",
			NewRGB15(color.RGBA64{R: 0xFFFF, G: 0x0000, B: 0x0000, A: 0xFFFF}),
			0xFFFF,
			0x0000,
			0x0000,
			0xFFFF,
		},
		{
			"green",
			NewRGB15(color.RGBA64{R: 0x0000, G: 0xFFFF, B: 0x0000, A: 0xFFFF}),
			0x0000,
			0xFFFF,
			0x0000,
			0xFFFF,
		},
		{
			"blue",
			NewRGB15(color.RGBA64{R: 0x0000, G: 0x0000, B: 0xFFFF, A: 0xFFFF}),
			0x0000,
			0x0000,
			0xFFFF,
			0xFFFF,
		},
		{
			"gray",
			NewRGB15(color.RGBA64{R: 0x4210, G: 0x4210, B: 0x39CE, A: 0x0000}),
			0x4210,
			0x4210,
			0x39CE,
			0xFFFF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotG, gotB, gotA := tt.rgb.RGBA()
			if gotR != tt.wantR {
				t.Errorf("RGB15.RGBA() gotR = %X, want %X", gotR, tt.wantR)
			}
			if gotG != tt.wantG {
				t.Errorf("RGB15.RGBA() gotG = %X, want %X", gotG, tt.wantG)
			}
			if gotB != tt.wantB {
				t.Errorf("RGB15.RGBA() gotB = %X, want %X", gotB, tt.wantB)
			}
			if gotA != tt.wantA {
				t.Errorf("RGB15.RGBA() gotA = %v, want %v", gotA, tt.wantA)
			}
		})
	}
}
