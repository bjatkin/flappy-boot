package alloc

import (
	"reflect"
	"testing"

	"github.com/bjatkin/flappy_boot/internal/hardware/memmap"
)

func TestPal_Alloc(t *testing.T) {
	var memBlock []memmap.PaletteValue
	// fill memBlock with data so deep equal can
	// check if slice sections are in fact equal
	for i := 0; i < 16*16; i++ {
		memBlock = append(memBlock, memmap.PaletteValue(i))
	}

	tests := []struct {
		name    string
		p       *Pal
		want    *PMem
		wantErr bool
	}{
		{
			"success",
			&Pal{
				meta:   [8]bool{true, true},
				memory: memBlock,
			},
			&PMem{
				Memory: memBlock[32:48],
				Offset: 2,
			},
			false,
		},
		{
			"error OOM",
			&Pal{
				meta:   [8]bool{true, true, true, true, true, true, true, true},
				memory: memBlock,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Alloc()
			if (err != nil) != tt.wantErr {
				t.Errorf("Pal.Alloc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pal.Alloc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPal_Free(t *testing.T) {
	type fields struct {
		meta   [8]bool
		memory []memmap.PaletteValue
	}
	type args struct {
		mem *PMem
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Pal
	}{
		{
			"success",
			fields{
				meta: [8]bool{false, false, true, true},
			},
			args{
				&PMem{
					Offset: 2,
				},
			},
			&Pal{
				meta: [8]bool{false, false, false, true},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Pal{
				meta:   tt.fields.meta,
				memory: tt.fields.memory,
			}
			p.Free(tt.args.mem)
			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Pal.Free() = %v, want %v", p, tt.want)
			}
		})
	}
}
