package math

// Rect represents a rectangle with 4 points at each corner
type Rect struct {
	X1, Y1 int
	X2, Y2 int
}

// V2 is a 2 element vector
type V2 struct {
	X, Y Fix8
}

// Lerp linerally interpolates between the values a, and b. t determins the percentage of the interpolation
// t is not clamped to the range (0-1) so you must do that yourself.
func Lerp(a, b, t Fix8) Fix8 {
	return a + (t*(b-a))>>8
}

// Clamp clamps the value to a maximum possible value
func Clamp(i, min Fix8, max Fix8) Fix8 {
	if i > max {
		return max
	}
	if i < min {
		return min
	}

	return i
}
