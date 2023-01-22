package gb16

import (
	"fmt"
	"image"
	"io"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/byteconv"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
)

// Encode writes the image as a valid 16 bit rgb image (.gb16) to the provided writer
// the provided image has a maximum width and height of 0xFFFF(65,535) pixels
// any larger and the function will return an error
func Encode(w io.Writer, m image.Image) error {
	dx, dy := m.Bounds().Dx(), m.Bounds().Dy()
	if dx > 0xFFFF || dy > 0xFFFF {
		return fmt.Errorf("image must be no larger than 65,535x65,535 [%dx%d]", dx, dy)
	}

	// 2 bytes for each pixel
	pixelSize := dx * dy * 2
	// 2 bytes for width and height
	headerSize := 4

	raw := make([]byte, 0, pixelSize+headerSize)
	raw = append(raw, byteconv.Itoa(uint16(dx))...)
	raw = append(raw, byteconv.Itoa(uint16(dy))...)

	gbaimg.Walk(m, func(x, y int) {
		c := m.At(x, y)
		c15 := gbaimg.RGB15Model.Convert(c).(gbacol.RGB15)
		raw = append(raw, c15.Bytes()...)
	})

	_, err := w.Write(raw)
	return err
}

// Decode can be used to decode a valid 15 bit rgb image (.gb16) from the given reader
func Decode(r io.Reader) (image.Image, string, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	dx := int(byteconv.Atoi(raw[0:2]))
	dy := int(byteconv.Atoi(raw[2:4]))
	img := gbaimg.NewRGB16(image.Rect(0, 0, dx, dy))
	i := 2
	gbaimg.Walk(img, func(x, y int) {
		c := gbacol.RGB15(byteconv.Atoi(raw[i : i+2]))
		img.Set(x, y, c)

		i += 2
	})

	return img, "gb16", nil
}
