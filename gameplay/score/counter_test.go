package score

import (
	"reflect"
	"testing"
)

func TestCounter_Set(t *testing.T) {
	type args struct {
		score int
	}
	tests := []struct {
		name string
		args args
		want [4]int
	}{
		{
			"max",
			args{
				score: 9999,
			},
			[4]int{9, 9, 9, 9},
		},
		{
			"min",
			args{
				score: 0,
			},
			[4]int{0, 0, 0, 0},
		},
		{
			"321",
			args{
				score: 321,
			},
			[4]int{0, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Counter{}
			c.Set(tt.args.score)

			if !reflect.DeepEqual(c.score, tt.want) {
				t.Errorf("Counter.Set() digits = %v, wanted = %v", c.score, tt.want)
			}
		})
	}
}
