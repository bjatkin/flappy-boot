package math

import (
	"testing"
)

func TestLerp(t *testing.T) {
	type args struct {
		a Fix8
		b Fix8
		t Fix8
	}
	tests := []struct {
		name string
		args args
		want Fix8
	}{
		{
			"initial",
			args{
				a: FixOne * 2,
				b: FixOne * 5,
				t: 0,
			},
			FixOne * 2,
		},
		{
			"final",
			args{
				a: FixOne,
				b: FixOne * 10,
				t: FixOne,
			},
			FixOne * 10,
		},
		{
			"50%",
			args{
				a: FixOne * 4,
				b: FixOne * 6,
				t: FixHalf,
			},
			FixOne * 5,
		},
		{
			"33%",
			args{
				a: FixOne * 2,
				b: FixOne * 3,
				t: FixThird,
			},
			FixOne*2 + FixThird,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lerp(tt.args.a, tt.args.b, tt.args.t); got != tt.want {
				t.Errorf("Lerp() = %v, want %v", got, tt.want)
			}
		})
	}
}
