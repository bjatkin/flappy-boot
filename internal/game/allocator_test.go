package game

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// hexArray prints an int array as hex strings to make debugging easier
func hexArray(array []int) string {
	var parts []string
	for _, i := range array {
		parts = append(parts, fmt.Sprintf("0x%X", i))
	}
	return "[" + strings.Join(parts, ",") + "]"
}

func TestAllocator_Alloc(t *testing.T) {
	type args struct {
		size []int
	}
	tests := []struct {
		name      string
		a         Allocator
		args      args
		want      []int
		wantErr   []bool
		wantAlloc Allocator
	}{
		{
			"single allocation",
			NewAllocator(5),
			args{
				size: []int{3},
			},
			[]int{0},
			[]bool{false, false, false},
			Allocator{used | 3, 0, 0, 2, 0},
		},
		{
			"multiple allocations",
			NewAllocator(10),
			args{
				size: []int{3, 2, 3, 1},
			},
			[]int{0, 3, 5, 8},
			[]bool{false, false, false, false},
			Allocator{used | 3, 0, 0, used | 2, 0, used | 3, 0, 0, used | 1, 1},
		},
		{
			"oom error",
			NewAllocator(10),
			args{
				size: []int{5, 4, 5},
			},
			[]int{0, 5, 0},
			[]bool{false, false, true},
			Allocator{used | 5, 0, 0, 0, 0, used | 4, 0, 0, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.args.size {
				got, err := tt.a.Alloc(tt.args.size[i])
				if (err != nil) != tt.wantErr[i] {
					t.Errorf("Allocator.Alloc() error = %v, wantErr %v", err, tt.wantErr[i])
					return
				}
				if got != tt.want[i] {
					t.Errorf("Allocator.Alloc() = %v, want %v", got, tt.want[i])
				}
			}
			if !reflect.DeepEqual(tt.a, tt.wantAlloc) {
				t.Errorf("Allocator.Alloc() = \n%v, want \n%v", hexArray(tt.a), hexArray(tt.wantAlloc))
			}
		})
	}
}

func TestAllocator_Free(t *testing.T) {
	type args struct {
		addr []int
	}
	tests := []struct {
		name string
		a    Allocator
		args args
		want Allocator
	}{
		{
			"free first cell",
			Allocator{used | 3, 0, 0, 2, 0},
			args{
				addr: []int{0},
			},
			Allocator{5, 0, 0, 0, 0},
		},
		{
			"free last cell",
			Allocator{3, 0, 0, used | 2, 0},
			args{
				addr: []int{3},
			},
			Allocator{5, 0, 0, 0, 0},
		},
		{
			"free middle cell",
			Allocator{used | 3, 0, 0, used | 2, 0, used | 3, 0, 0, used | 1, 1},
			args{
				addr: []int{3},
			},
			Allocator{used | 3, 0, 0, 2, 0, used | 3, 0, 0, used | 1, 1},
		},
		{
			"free and merge middle",
			Allocator{3, 0, 0, used | 2, 0, used | 3, 0, 0, 1, used | 1},
			args{
				addr: []int{5, 3},
			},
			Allocator{9, 0, 0, 0, 0, 0, 0, 0, 0, used | 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := range tt.args.addr {
				tt.a.Free(tt.args.addr[i])
			}
			if !reflect.DeepEqual(tt.a, tt.want) {
				t.Errorf("Allocator.Free() = \n%v, want \n%v", hexArray(tt.a), hexArray(tt.want))
			}
		})
	}
}
