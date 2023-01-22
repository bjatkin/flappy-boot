package byteconv

type Ints interface {
	int8 | int16 | int32 | int64 |
		uint8 | uint16 | uint32 | uint64
}

// Itoa converts an integer into a byte slice
func Itoa[T Ints](i T) []byte {
	switch i := any(i).(type) {
	case int8:
		return []byte{byte(i)}
	case int16:
		return []byte{byte(i), byte(i >> 8)}
	case int32:
		return []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
	case int64:
		return []byte{
			byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24),
			byte(i >> 32), byte(i >> 40), byte(i >> 48), byte(i >> 56),
		}
	case uint8:
		return []byte{byte(i)}
	case uint16:
		return []byte{byte(i), byte(i >> 8)}
	case uint32:
		return []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
	case uint64:
		return []byte{
			byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24),
			byte(i >> 32), byte(i >> 40), byte(i >> 48), byte(i >> 56),
		}
	}
	return nil
}

// Atoi converts a byte slice back into an integer
func Atoi(bytes []byte) int64 {
	l := len(bytes)
	switch {
	case l >= 8:
		return (int64(bytes[7]) << 56) | (int64(bytes[6]) << 48) |
			(int64(bytes[5]) << 40) | (int64(bytes[4]) << 32) |
			(int64(bytes[3]) << 24) | (int64(bytes[2]) << 16) |
			(int64(bytes[1]) << 8) | int64(bytes[0])
	case l >= 4:
		return (int64(bytes[3]) << 24) | (int64(bytes[2]) << 16) |
			(int64(bytes[1]) << 8) | int64(bytes[0])
	case l >= 2:
		return (int64(bytes[1]) << 8) | int64(bytes[0])
	case l >= 1:
		return int64(bytes[0])
	default:
		return 0
	}
}

// Atoi converts a byte slice back into a unsigned integer
func Atou(bytes []byte) uint64 {
	l := len(bytes)
	switch {
	case l >= 8:
		return (uint64(bytes[7]) << 56) | (uint64(bytes[6]) << 48) |
			(uint64(bytes[5]) << 40) | (uint64(bytes[4]) << 32) |
			(uint64(bytes[3]) << 24) | (uint64(bytes[2]) << 16) |
			(uint64(bytes[1]) << 8) | uint64(bytes[0])
	case l >= 4:
		return (uint64(bytes[3]) << 24) | (uint64(bytes[2]) << 16) |
			(uint64(bytes[1]) << 8) | uint64(bytes[0])
	case l >= 2:
		return (uint64(bytes[1]) << 8) | uint64(bytes[0])
	case l >= 1:
		return uint64(bytes[0])
	default:
		return 0
	}

}
