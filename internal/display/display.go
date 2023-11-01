package display

// Color represents a 15 bit RGB(grb) color
type Color uint16

// RGB15 converts the red green and blue channels into a valid 15 bit gba Color
// the red, green, and blue values MUST be in the range of 0-31
// this function will not check these values before constructing the color value
// so invalid values will cause the resulting color to be incorrect
func RGB15(red, green, blue uint) Color {
	return Color(red) | Color(green<<5) | Color(blue<<10)
}
