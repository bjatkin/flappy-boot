package alloc

import (
	"reflect"
	"testing"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

func TestVRAM_Alloc(t *testing.T) {
	var memBlock []memmap.VRAMValue
	// fill memBlock with data so deep equal can
	// check if slice sections are in fact equal
	for i := 0; i < 100; i++ {
		memBlock = append(memBlock, memmap.VRAMValue(i))
	}

	type args struct {
		size []int
	}
	tests := []struct {
		name     string
		m        VRAM
		args     args
		want     []*VMem
		wantErr  []bool
		wantMeta VRAM
	}{
		{
			"single allocation",
			NewVRAM(memBlock[:20], 4),
			args{
				size: []int{3},
			},
			[]*VMem{{
				Memory: memBlock[:12],
				Offset: 0,
			}},
			[]bool{false, false, false},
			VRAM{
				meta:     []int{used | 3, 0, 0, 2, 0},
				memory:   memBlock[:20],
				cellSize: 4,
			},
		},
		{
			"multiple allocations",
			NewVRAM(memBlock, 10),
			args{
				size: []int{3, 2, 3, 1},
			},
			[]*VMem{
				{
					Memory: memBlock[:30],
					Offset: 0,
				},
				{
					Memory: memBlock[30:50],
					Offset: 3,
				},
				{
					Memory: memBlock[50:80],
					Offset: 5,
				},
				{
					Memory: memBlock[80:90],
					Offset: 8,
				},
			},
			[]bool{false, false, false, false},
			VRAM{
				meta:     []int{used | 3, 0, 0, used | 2, 0, used | 3, 0, 0, used | 1, 1},
				memory:   memBlock,
				cellSize: 10,
			},
		},
		{
			"oom error",
			NewVRAM(memBlock[:20], 2),
			args{
				size: []int{5, 4, 5},
			},
			[]*VMem{
				{
					Memory: memBlock[:10],
					Offset: 0,
				},
				{
					Memory: memBlock[10:18],
					Offset: 5,
				},
				nil,
			},
			[]bool{false, false, true},
			VRAM{
				meta:     []int{used | 5, 0, 0, 0, 0, used | 4, 0, 0, 0, 1},
				memory:   memBlock[:20],
				cellSize: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.args.size {
				got, err := tt.m.Alloc(tt.args.size[i])
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("VRAM.Alloc() error = %v, wantErr %v (%d)", err, tt.wantErr[i], i)
					return
				}
				if !reflect.DeepEqual(got, tt.want[i]) {
					t.Errorf("VRAM.Alloc() = %v, want %v (%d)", got, tt.want[i], i)
				}
			}
			if !reflect.DeepEqual(tt.m, tt.wantMeta) {
				t.Errorf("VRAM.Alloc() = \n%v, want \n%v", tt.m, tt.wantMeta)
			}
		})
	}
}

func TestVRAM_Free(t *testing.T) {
	var memBlock []memmap.VRAMValue
	// fill memBlock with data so deep equal can
	// check if slice sections are in fact equal
	for i := 0; i < 100; i++ {
		memBlock = append(memBlock, memmap.VRAMValue(i))
	}

	type args struct {
		mem []*VMem
	}
	tests := []struct {
		name string
		v    VRAM
		args args
		want VRAM
	}{
		{
			"free first cell",
			VRAM{
				meta:     []int{used | 3, 0, 0, 2, 0},
				memory:   memBlock[:50],
				cellSize: 10,
			},
			args{
				mem: []*VMem{{
					Memory: memBlock[0:30],
					Offset: 0,
				}},
			},
			VRAM{
				meta:     []int{5, 0, 0, 0, 0},
				memory:   memBlock[:50],
				cellSize: 10,
			},
		},
		{
			"free last cell",
			VRAM{
				meta:     []int{3, 0, 0, used | 2, 0},
				memory:   memBlock[:50],
				cellSize: 10,
			},
			args{
				mem: []*VMem{{
					Memory: memBlock[30:50],
					Offset: 3,
				}},
			},
			VRAM{
				meta:     []int{5, 0, 0, 0, 0},
				memory:   memBlock[:50],
				cellSize: 10,
			},
		},
		{
			"free middle cell",
			VRAM{
				meta:     []int{used | 3, 0, 0, used | 2, 0, used | 3, 0, 0, used | 1, 1},
				memory:   memBlock,
				cellSize: 10,
			},
			args{
				mem: []*VMem{{
					Memory: memBlock[30:50],
					Offset: 3,
				}},
			},
			VRAM{
				meta:     []int{used | 3, 0, 0, 2, 0, used | 3, 0, 0, used | 1, 1},
				memory:   memBlock,
				cellSize: 10,
			},
		},
		{
			"free and merge middle",
			VRAM{
				meta:     []int{3, 0, 0, used | 2, 0, used | 3, 0, 0, 1, used | 1},
				memory:   memBlock,
				cellSize: 10,
			},
			args{
				mem: []*VMem{
					{
						Memory: memBlock[50:80],
						Offset: 5,
					},
					{
						Memory: memBlock[30:50],
						Offset: 3,
					},
				},
			},
			VRAM{
				meta:     []int{9, 0, 0, 0, 0, 0, 0, 0, 0, used | 1},
				memory:   memBlock,
				cellSize: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.args.mem {
				tt.v.Free(tt.args.mem[i])
			}
			if !reflect.DeepEqual(tt.v, tt.want) {
				t.Errorf("Allocator.Free() = \n%v, want \n%v", tt.v, tt.want)
			}
		})
	}
}
