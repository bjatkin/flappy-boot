# Flappy Boot
This is a simple flappy bird clone for the GBA.
You play as hermes boot trying to return to him.

# Project Structure
This project has the following structure.
    * cmd: tools used as part of game development
        * image_gen: conversion tool used to generate GBA compatible graphics from png image files
    * gameplay: all gameplay related code.
    * internal: internal engine code. The core logic that the game is built on top of.
        * hardware: GBA hardware related code, includes things like hardware registers and memory offsets
        * fix: the fixed point number type used extensivly through the project. It is has 24 whole number bits and 8 fractional bits

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


