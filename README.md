# Flappy Boot
Oh No! Hermes, the Olympian god, seems to have dropped on of his winged boots from the heavens!
Better hurry and find your way back to him, but beware of the many Roman columns that stand in your way.

This is a flappy bird clone written from scratch for the GBA.
It is open source and fairly well commented so feel free to use it as a jumping off point for your own project.
If you would like to learn about this project check out [this presentation](https://youtu.be/mrWJZSVSRVQ?si=653ayqaEtqz5xB6o) on makeing GBA games in Go.

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
* config.yaml: configuration for the image_gen tool.

# Running Flappy Boot
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
This game with built with the help of a couple of really excelent  GBA programming resources.
These are both tailored for C/C++ development on the GBA but still provied great general knowledge

TONIC: https://www.coranac.com/tonc/text/toc.htm

GBATEK: https://problemkaputt.de/gbatek.htm


