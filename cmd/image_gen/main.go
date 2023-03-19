package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/config"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/exit"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/generate"
)

func main() {
	// add this first so any other defers are run before we exit
	defer exit.Final()

	// look for an input yaml file
	if len(os.Args) != 2 {
		exit.Error(exit.InvalidArguments, fmt.Errorf("invalid usage, missing config file: %v", os.Args))
		return
	}

	file := os.Args[1]
	cfg, err := config.NewConfigFromFile(file)
	if err != nil {
		exit.Error(exit.InvalidConfig, fmt.Errorf("failed to read in config %w", err))
		return
	}

	err = cfg.Validate()
	if err != nil {
		exit.Error(exit.InvalidConfig, fmt.Errorf("failed to validate config %w", err))
		return
	}

	var setTransparent *gbacol.RGB15
	if cfg.SetTransparent != "" {
		tmp, err := config.ParseHexColor(cfg.SetTransparent)
		if err != nil {
			exit.Error(exit.InvalidConfig, fmt.Errorf("%s is not a valid hex color %w", cfg.SetTransparent, err))
			return
		}
		setTransparent = &tmp
	}

	palettes := make(map[string]*generate.PaletteData)
	for _, pal := range cfg.Palettes {
		palData, err := generate.NewPaletteData(pal, setTransparent)
		if err != nil {
			exit.Error(exit.InvalidPalette, fmt.Errorf("falied to generate new pallet %s | %w", pal.Name, err))
			return
		}
		palettes[pal.Name] = palData
	}

	tileSets := make(map[string]*generate.TileSetData)
	for _, tileSet := range cfg.TileSets {
		tileSetData, err := generate.NewTileSetData(tileSet, setTransparent, palettes)
		if err != nil {
			exit.Error(exit.InvalidTileSet, fmt.Errorf("failed to generate new tile set %s | %w", tileSet.Name, err))
			return
		}
		tileSets[tileSet.Name] = tileSetData
	}

	tileMaps := make(map[string]*generate.TileMapData)
	for _, tileMap := range cfg.TileMaps {
		tileMapData, err := generate.NewTileMapData(tileMap, setTransparent, tileSets, palettes)
		if err != nil {
			exit.Error(exit.InvalidTileMap, fmt.Errorf("failed to generate tile map %s | %w", tileMap.Name, err))
			return
		}
		tileMaps[tileMap.Name] = tileMapData
	}

	err = generate.WriteAssetFile(cfg.OutDir)
	if err != nil {
		exit.Error(exit.FileWriteFailed, fmt.Errorf("failed to write base asset.go file %w", err))
		return
	}

	for _, pal := range palettes {
		if pal.Shared <= 1 {
			continue
		}
		// this is a shared palette
		err := writeFiles(pal, cfg.OutDir, pal.Name+".pal4", pal.Name+"Palette.go")
		if err != nil {
			exit.Error(exit.InvalidPalette, err)
			return
		}
	}

	for _, tileSet := range tileSets {
		// TODO: how should I handle tile sets used for sprites?
		if tileSet.Shared <= 1 {
			continue
		}
		// this is a shared tile set
		err := writeFiles(tileSet, cfg.OutDir, tileSet.Name+".ts4", tileSet.Name+"TileSet.go")
		if err != nil {
			exit.Error(exit.InvalidTileSet, err)
			return
		}
	}

	for _, tileMap := range tileMaps {
		err := writeFiles(tileMap, cfg.OutDir, tileMap.Name+".tm4", tileMap.Name+"TileMap.go")
		if err != nil {
			exit.Error(exit.InvalidTileMap, err)
			return
		}
	}
}

func writeFiles(f generate.File, dir, rawName, goName string) error {
	rawFile := filepath.Join(dir, rawName)
	rawBytes, err := f.Raw()
	if err != nil {
		return err
	}

	err = os.WriteFile(rawFile, rawBytes, 0o0666)
	if err != nil {
		return fmt.Errorf("failed to save file %s | %w", rawFile, err)
	}

	goFile := filepath.Join(dir, goName)
	goRaw, err := f.Go()
	if err != nil {
		return fmt.Errorf("failed to generate go file %s | %w", goFile, err)
	}

	err = os.WriteFile(goFile, goRaw, 0o0666)
	if err != nil {
		return fmt.Errorf("failed to save go file %s | %w", goFile, err)
	}

	return nil
}
