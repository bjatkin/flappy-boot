# ImageGen

ImageGen is the support command used by this project to generate GBA compatable graphics.
It can be used to generate both sprite sheets and background tile maps.
The final output of the program will be both raw data files (.pal4, .tm4 and .ts4) as well as associated go files.
Generated go files will use the `assets` pack.
For an example view the `config.yaml` file in the base directory of this repo

## config
ImageGen takes a config file as it's only argument.
This config file controlls the output of ImageGen and supports the following attributes

#### OutDir
This configures the output directory where generated asset and go files will be placed.
This directory can be an absolute directory or a relative one.
Relative directories are considered reletive to the location where image gen was run, not the location of the file iteslf.

#### SetTransparent
This attribute must be a valid 3 or 6 digit hex color.
After asset conversion is complete this settings overwrites the assets transparent color.
This setting is usefull when you want to controll the "clear" color of the screen.

#### TileSets
This is a list of tile sets and their associated attributes

* Name: the name of the tile set
* File: the image file associated with the tile set
* Size: the tile size to use for the tile set, valid GBA tile sizes are
    * 8x8
    * 8x16
    * 16x8
    * 16x16
    * 32x8
    * 8x32
    * 32x32
    * 32x16
    * 16x32
    * 64x64
    * 64x32
    * 32x64
* Description: a description of the tile set, this will be added to the generated code
* Palette: the name of a palette defined in the config. This palette will be used when converting the tile set
* Transparent: the hex color to use as the transparent color in the asset. This can not be used if Palette is also set

#### TileMaps
this is a list of the tile maps and their associated attributes.
Note that tile maps always use 8x8 tiles.

* Name: the name of the tile map
* File: the image file assocated with the tile set
* TileSet: the name of a tile set defined in the config. This tile set will be used when converting the tile map
* Palette: the name of a palette defined in the config. This palette will be used when converting the tile map. This can not be used if the TileSet is set.
* Description: a description of the tile map, this will be added to the generated code
* Transparent: the hex color to use as the transparent color in the asset. This can not be used if either TileSet or Palette are also set.
