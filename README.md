# Flappy Boot
![flappy boot title](https://github.com/bjatkin/flappy-boot/blob/main/assets/title.png)
Oh No! Hermes, the Olympian god, seems to have dropped on of his winged boots from the heavens!
Better hurry and find your way back to him, but beware of the many Roman columns that stand in your way.

This is a flappy bird clone written from scratch for the GBA.
It is open source and fairly well commented so feel free to use it as a jumping off point for your own project.
If you would like to learn about this project check out [this presentation](https://youtu.be/mrWJZSVSRVQ?si=653ayqaEtqz5xB6o) on makeing GBA games in Go.


# Play Me
![flappy boot gameplay](https://github.com/bjatkin/flappy-boot/blob/main/assets/gameplay.gif)
You can play this game on [itch.io](https://aanval.itch.io/flappy-boot-advance).
itch.io will allow you do download the gba rom to play on any gba emulator.
It also contains a web build so you can play it in your browser.
The web player is a customer emulator built using [ebitengine](https://ebitengine.org/).
You can look at the emulator code [here](https://github.com/bjatkin/flappy-boot/tree/main/internal/emu).

# Project Structure
This project has the following structure.
* assets: png assets and mockups for the game
* cmd: tools used as part of game development
    * image_gen: conversion tool used to generate GBA compatible graphics from png image files.
    * lut: look up table generation for the sin function.
* gameplay: all gameplay related code.
* internal: internal engine code. The core logic that the game is built on top of.
    * fix: the fixed point number type used extensivly through the project. It is has 24 whole number bits and 8 fractional bits.
    * alloc: memory allocators for the gba's VRAM and Paletts memory.
    * assets: generated assets that are used directly in the engine.
    * display: display and color related code for the engine.
    * emu/ppu: a simple ppu emulator that allows standalone and web builds.
    * key: key codes for input handling.
    * lut: look up tables for the sin function.
    * math: some simple math focused utilities.
    * game: the code for the game engine.
    * hardware: GBA hardware related code, includes things like hardware registers and memory offsets.
        * audio: some of the basic audio registers. (unused)
        * display: display related registers.
        * dma: registers for direct memory access. (unused)
        * key: input related registers.
        * memmap: gba memory layout and register access. 
        * sprite: oam and palette memory.
        * timer: some of the basic gba timer registers. (unused)
        * save: registers and memory related to save data on the GBA. It specifically supports FRAM style hardware.
* config.yaml: configuration for the image_gen tool.
* wasm: all code related to the frontend web build.

# Running Flappy Boot
Flappy Boot can be run in 3 different modes:
1) as a standalone game using an emulated PPU
2) as a wasm build using wasm and npm
3) as a GBA ROM inside an emulator/ or on actual hardware

You can build all these files for all these targets using the `./build` script.
Note that you should only run `./build` from the root directory of the repo otherwise it will fail.

### Standalone
The Standalone build can be built using the normal go build tool.
You will need to include the `standalone` and `local` build tags however.
```sh
go build -tags=standalone,local .
```
when run this game will create a `flappy_boot_stand.sav` file, which contains the high score save data.

### Web
You can build the flappy bird file for `.wasm` using the following command.
```sh
env GOOS=js GOARCH=wasm go build -tags=standalone,web -o wasm/flappy_boot.wasm github.com/bjatkin/flappy_boot
```

If this is the first time building the wasm file you will nee to run
```sh
npm init
```

You can play this in a browser by entering the `wasm` directory and then running
```sh
npm run dev
```

note that the PPU emulator doesn't quite performe as well as the standalone or emulated versions of the game.
For the best experience, you should play one of the other verions.

### GBA ROM
First ensure you have the [tiny-go complier](https://tinygo.org/getting-started/install/) installed.
This is the complier that this project uses and you will not be able to complie without it.

Next you'll likely want to install the `mgba-qt` eumlator.
This is not required however, and any GBA emulator can be used to run flappy boot.

If you have `mgba-qt` available to run from your terminal you can simply run the `./run` shell file.
This will test all the code, re-generate all assets files, complie the code, and then start the game in `mgba-qt`.

If you do not have this emulator installed you'll need to run the `./build` shell file instead
This performs all the same steps as `./run` however, it will not attempt to start the game.
Instead it will create a file called `flappy_boot.gba`.
This `.gba` file can then be run with any sutiable GBA emulator.

# References
This project was made possible because of the awesome [tiny go complier](https://tinygo.org/),
as well as those who worked to get support for the [GBA compile target](https://tinygo.org/docs/reference/microcontrollers/gameboy-advance/).

Flappy boot was built with the help of a couple of really excelent GBA programming resources.
These are both tailored for C/C++ development on the GBA but still provied great general knowledge.
* [TONIC](https://www.coranac.com/tonc/text/toc.htm)
* [GBATEK](https://problemkaputt.de/gbatek.htm)

The GBxCart RW was used to flash flappy boot to a phisical repoduction cart.
* [GBxCart RW](https://www.gbxcart.com/)

There are also several talks related to the development of this project.
* [Go West 2023](https://www.youtube.com/live/5qvfAc1C2Kg?si=mJF-05vCzNOkdn4z&t=5955)
* [Forge Utah Go Users Group](https://youtu.be/mrWJZSVSRVQ?si=SfUO3th8HURi0uJj)
