#!/bin/bash
set -euxo pipefail

./build.sh
mgba-qt -4 flappy_boot.gba