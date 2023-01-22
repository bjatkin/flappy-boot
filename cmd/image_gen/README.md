# ImageGen

ImageGen is a generation tool that can be used by go:generate to convert image files into a go package.
By default it will create an `assets` package in the `internal/assets` directory but this can be confgured using the `-o` option.
Each invocation of `image_gen` will create a new file in the package.
This means it's safe to invoke the command multiple times as part of your compliation step.
By default files will be named based on the image files they are converting but again this can be altered using the `-o` option.

# Arguments

ImageGen supports the following arguments.
These arguments can be repeated to convert more that one image at a time.
The resulting image data will then share the same color palette.

* Image - the image to convert into tile data.
    This can be a `.png`,`.jpeg`, or a `.gif` image 

* Size - the size of the tiles to use when converting the image data.
    The GBA supports the following tiles sizes
    * 8x8
    * 16x8
    * 8x16
    * 16x16
    * 32x8
    * 8x32
    * 32x32
    * 32x16
    * 16x32
    * 64x64
    * 64x32
    * 32x64

# Flags

ImageGen supports the following optional flags

* transparent (t) - By default, the first color encountered in a image will be set as the transparent color in the palette.
    This flage over-rides that bevaior and allows the transparent color to be specified as an input. 
    It supports the following color formats:
    * 3 digit hex color `#F0F`
    * 6 digit hex color `#FF00FF`
    * RGB color with 8 bit color channels `(r, g, b)`

* p256 - Tells ImageGen to generate image data with a 256 color palette rather than the default 16 color palette