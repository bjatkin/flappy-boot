package main

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/config"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/generate"
)

func main() {
	// look for an input yaml file
	if len(os.Args) != 2 {
		fmt.Println("invalie usage, missing config file", os.Args)
		return
	}

	file := os.Args[1]
	cfg, err := config.NewConfigFromFile(file)
	if err != nil {
		fmt.Printf("failed to read in config %s\n", err)
		return
	}

	err = cfg.Validate()
	if err != nil {
		fmt.Printf("failed to validate config %s\n", err)
		return
	}

	palettes := make(map[string]*generate.PaletteData)
	for _, pal := range cfg.Palettes {
		palData, err := generate.NewPaletteData(pal)
		if err != nil {
			fmt.Printf("falied to generate new pallet %s | %s\n", pal.Name, err)
			return
		}
		palettes[pal.Name] = palData
	}

	tileSets := make(map[string]*generate.TileSetData)
	for _, tileSet := range cfg.TileSets {
		tileSetData, err := generate.NewTileSetData(tileSet, palettes)
		if err != nil {
			fmt.Printf("failed to generate new tile set %s | %s\n", tileSet.Name, err)
			return
		}
		tileSets[tileSet.Name] = tileSetData
	}

	tileMaps := make(map[string]*generate.TileMapData)
	for _, tileMap := range cfg.TileMaps {
		tileMapData, err := generate.NewTileMapData(tileMap, tileSets, palettes)
		if err != nil {
			fmt.Printf("failed to generate tile map %s | %s\n", tileMap.Name, err)
			return
		}
		tileMaps[tileMap.Name] = tileMapData
	}

	err = generate.WriteAssetFile(cfg.OutDir)
	if err != nil {
		fmt.Printf("failed to write base asset.go file %s\n", err)
		return
	}

	for _, pal := range palettes {
		if pal.Shared <= 1 {
			continue
		}
		// this is a shared palette
		err := writeFiles(pal, cfg.OutDir, pal.Name+".pal4", pal.Name+"Palette.go")
		if err != nil {
			fmt.Println(err)
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
			fmt.Println(err)
			return
		}
	}

	for _, tileMap := range tileMaps {
		err := writeFiles(tileMap, cfg.OutDir, tileMap.Name+".tm4", tileMap.Name+"TileMap.go")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func writeFiles(f generate.File, dir, rawName, goName string) error {
	rawFile := filepath.Join(dir, rawName)
	err := os.WriteFile(rawFile, f.Raw(), 0o0666)
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
