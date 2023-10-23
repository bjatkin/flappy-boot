package gameplay

import "testing"

func Test_scoreToSaveData(t *testing.T) {
	type args struct {
		score int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"low score",
			args{
				score: 10,
			},
		},
		{
			"high score",
			args{
				score: 0x7FFF,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := saveDataToScore(scoreToSaveData(tt.args.score)); got != tt.args.score {
				t.Errorf("scoreToSaveData() = %v, want %v", got, tt.args.score)
			}
		})
	}
}
