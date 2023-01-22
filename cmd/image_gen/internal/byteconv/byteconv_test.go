package byteconv

import (
	"reflect"
	"testing"
)

func TestItoa(t *testing.T) {
	type args struct {
		i8  int8
		i16 int16
		i32 int32
		i64 int64
		u8  uint8
		u16 uint16
		u32 uint32
		u64 uint64
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			"convert int8",
			args{
				i8: 0x2F,
			},
			[]byte{0x2F},
		},
		{
			"convert int16",
			args{
				i16: 0x1A_2B,
			},
			[]byte{0x2B, 0x1A},
		},
		{
			"convert int32",
			args{
				i32: 0x1A_2B_3C_4D,
			},
			[]byte{0x4D, 0x3C, 0x2B, 0x1A},
		},
		{
			"convert int64",
			args{
				i64: 0x1A_2B_3C_4D_5E_6F_7A_8B,
			},
			[]byte{0x8B, 0x7A, 0x6F, 0x5E, 0x4D, 0x3C, 0x2B, 0x1A},
		},
		{
			"convert unt8",
			args{
				u8: 0xFF,
			},
			[]byte{0xFF},
		},
		{
			"convert unt16",
			args{
				u16: 0xAA_BB,
			},
			[]byte{0xBB, 0xAA},
		},
		{
			"convert unt32",
			args{
				u32: 0xAA_BB_CC_DD,
			},
			[]byte{0xDD, 0xCC, 0xBB, 0xAA},
		},
		{
			"convert unt64",
			args{
				u64: 0xAA_BB_CC_DD_EE_FF_AA_BB,
			},
			[]byte{0xBB, 0xAA, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []byte
			switch {
			case tt.args.i8 > 0:
				got = Itoa(tt.args.i8)
			case tt.args.i16 > 0:
				got = Itoa(tt.args.i16)
			case tt.args.i32 > 0:
				got = Itoa(tt.args.i32)
			case tt.args.i64 > 0:
				got = Itoa(tt.args.i64)
			case tt.args.u8 > 0:
				got = Itoa(tt.args.u8)
			case tt.args.u16 > 0:
				got = Itoa(tt.args.u16)
			case tt.args.u32 > 0:
				got = Itoa(tt.args.u32)
			case tt.args.u64 > 0:
				got = Itoa(tt.args.u64)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Itoa() = 0x%X, want 0x%X", got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"convert int8",
			args{
				bytes: []byte{0x1A},
			},
			0x1A,
		},
		{
			"convert int16",
			args{
				bytes: []byte{0x2B, 0x1A},
			},
			0x1A_2B,
		},
		{
			"convert int32",
			args{
				bytes: []byte{0x4D, 0x3C, 0x2B, 0x1A},
			},
			0x1A_2B_3C_4D,
		},
		{
			"conver int64",
			args{
				bytes: []byte{0x8B, 0x7A, 0x6F, 0x5E, 0x4D, 0x3C, 0x2B, 0x1A},
			},
			0x1A_2B_3C_4D_5E_6F_7A_8B,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Atoi(tt.args.bytes); got != tt.want {
				t.Errorf("Atoi() = 0x%X, want 0x%X", got, tt.want)
			}
		})
	}
}

func TestAtou(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			"convert uint8",
			args{
				bytes: []byte{0xAA},
			},
			0xAA,
		},
		{
			"convert uint16",
			args{
				bytes: []byte{0xBB, 0xAA},
			},
			0xAA_BB,
		},
		{
			"convert uint32",
			args{
				bytes: []byte{0xDD, 0xCC, 0xBB, 0xAA},
			},
			0xAA_BB_CC_DD,
		},
		{
			"conver uint64",
			args{
				bytes: []byte{0xBB, 0xAA, 0xFF, 0xEE, 0xDD, 0xCC, 0xBB, 0xAA},
			},
			0xAA_BB_CC_DD_EE_FF_AA_BB,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Atou(tt.args.bytes); got != tt.want {
				t.Errorf("Atou() = %v, want %v", got, tt.want)
			}
		})
	}
}

func FuzzConvertUint8(f *testing.F) {
	f.Add(uint8(0xFF))
	f.Add(uint8(0x00))
	f.Fuzz(func(t *testing.T, i uint8) {
		if got := Atou(Itoa(i)); uint8(got) != i {
			t.Errorf("DoubleConvertUint8 = 0x%X, want 0x%X", got, i)
		}
	})

}

func FuzzConvertUint16(f *testing.F) {
	f.Add(uint16(0xFF_FF))
	f.Add(uint16(0x00_00))
	f.Fuzz(func(t *testing.T, i uint16) {
		if got := Atou(Itoa(i)); uint16(got) != i {
			t.Errorf("DoubleConvertUint16 = 0x%X, want 0x%X", got, i)
		}
	})

}

func FuzzConvertUint32(f *testing.F) {
	f.Add(uint32(0xFF_FF_FF_FF))
	f.Add(uint32(0x00_00_00_00))
	f.Fuzz(func(t *testing.T, i uint32) {
		if got := Atou(Itoa(i)); uint32(got) != i {
			t.Errorf("DoubleConvertUint32 = 0x%X, want 0x%X", got, i)
		}
	})
}

func FuzzConvertUint64(f *testing.F) {
	f.Add(uint64(0xFF_FF_FF_FF_FF_FF_FF_FF))
	f.Add(uint64(0x00_00_00_00_00_00_00_00))
	f.Fuzz(func(t *testing.T, i uint64) {
		if got := Atou(Itoa(i)); uint64(got) != i {
			t.Errorf("DoubleConvertUint64 = 0x%X, want 0x%X", got, i)
		}
	})
}

func FuzzConvertInt8(f *testing.F) {
	f.Add(int8(0x7F))
	f.Add(int8(-0x7F))
	f.Add(int8(0x00))
	f.Fuzz(func(t *testing.T, i int8) {
		if got := Atoi(Itoa(i)); int8(got) != i {
			t.Errorf("DoubleConvertInt8 = 0x%X, want 0x%X", got, i)
		}
	})

}

func FuzzConvertInt16(f *testing.F) {
	f.Add(int16(0x7F_FF))
	f.Add(int16(-0x7F_FF))
	f.Add(int16(0x00_00))
	f.Fuzz(func(t *testing.T, i int16) {
		if got := Atoi(Itoa(i)); int16(got) != i {
			t.Errorf("DoubleConvertInt16 = 0x%X, want 0x%X", got, i)
		}
	})

}

func FuzzConvertInt32(f *testing.F) {
	f.Add(int32(0x7F_FF_FF_FF))
	f.Add(int32(-0x7F_FF_FF_FF))
	f.Add(int32(0x00_00_00_00))
	f.Fuzz(func(t *testing.T, i int32) {
		if got := Atoi(Itoa(i)); int32(got) != i {
			t.Errorf("DoubleConvertInt32 = 0x%X, want 0x%X", got, i)
		}
	})
}

func FuzzConvertInt64(f *testing.F) {
	f.Add(int64(0x7F_FF_FF_FF_FF_FF_FF_FF))
	f.Add(int64(-0x7F_FF_FF_FF_FF_FF_FF_FF))
	f.Add(int64(0x00_00_00_00_00_00_00_00))
	f.Fuzz(func(t *testing.T, i int64) {
		if got := Atoi(Itoa(i)); got != i {
			t.Errorf("DoubleConvertUint64 = 0x%X, want 0x%X", got, i)
		}
	})
}
