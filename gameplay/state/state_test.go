package state

import (
	"reflect"
	"testing"

	"github.com/bjatkin/flappy_boot/internal/math"
)

func TestTracker_Frac(t *testing.T) {
	type fields struct {
		SceneFrames map[State]int
		frame       int
	}
	tests := []struct {
		name   string
		fields fields
		want   math.Fix8
	}{
		{
			"zero",
			fields{
				SceneFrames: map[State]int{
					A: 30,
				},
				frame: 0,
			},
			0,
		},
		{
			"one",
			fields{
				SceneFrames: map[State]int{
					A: 30,
				},
				frame: 30,
			},
			math.FixOne,
		},
		{
			"half",
			fields{
				SceneFrames: map[State]int{
					A: 30,
				},
				frame: 15,
			},
			math.FixHalf,
		},
		{
			"one quarter",
			fields{
				SceneFrames: map[State]int{
					A: 20,
				},
				frame: 5,
			},
			math.FixQuarter,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Tracker{
				SceneFrames: tt.fields.SceneFrames,
				frame:       tt.fields.frame,
				state:       A,
			}
			if got := tr.Frac(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tracker.Frac() = %v, want %v", got, tt.want)
			}
		})
	}
}
