package config

import (
	"errors"
	"fmt"
	"image/color"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gba/tile"
)

// Config is a config file for the command
type Config struct {
	Palettes       []Palette `yaml:"Paletts"`
	TileSets       []TileSet `yaml:"TileSets"`
	TileMaps       []TileMap `yaml:"TileMaps"`
	OutDir         string    `yaml:"OutDir"`
	SetTransparent string    `yaml:"SetTransparent"`
}

// Palette is a named palette
type Palette struct {
	Name        string `yaml:"Name"`
	File        string `yaml:"File"`
	Description string `yaml:"Description"`
	Transparent string `yaml:"Transparent"`
}

// TileSet is a named tile set
type TileSet struct {
	Name        string `yaml:"Name"`
	File        string `yaml:"File"`
	Palette     string `yaml:"Palette"`
	Size        string `yaml:"Size"`
	Transparent string `yaml:"Transparent"`
}

// TileMap is a named tile map
type TileMap struct {
	Name        string `yaml:"Name"`
	File        string `yaml:"File"`
	TileSet     string `yaml:"TileSet"`
	Palette     string `yaml:"Palette"`
	Transparent string `yaml:"Transparent"`
}

// NewConfigFromFile reads in the yaml file at the provided file location and then marshalls it into a new config struct
func NewConfigFromFile(file string) (*Config, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Validate ensure the configuration file has a valid format
func (c *Config) Validate() error {
	if c.OutDir == "" {
		return errors.New("output directory must be set")
	}

	if _, err := os.Stat(c.OutDir); err != nil {
		return fmt.Errorf("failed to access output directory %s | %w", c.OutDir, err)
	}

	palettes := make(map[string]Palette)
	for _, pal := range c.Palettes {
		palettes[pal.Name] = pal

		err := validateColor(pal.Transparent)
		if err != nil {
			return fmt.Errorf("invalid palette color %s | %w", pal.Transparent, err)
		}
	}

	tileSets := make(map[string]TileSet)
	for _, tileSet := range c.TileSets {
		tileSets[tileSet.Name] = tileSet
		_, err := tile.NewSize(tileSet.Size)
		if err != nil {
			return err
		}

		if tileSet.Palette != "" && tileSet.Transparent != "" {
			return fmt.Errorf("can not set transparent color when palette %s is set", tileSet.Palette)
		}

		err = validateColor(tileSet.Transparent)
		if err != nil {
			return err
		}

		if tileSet.Palette == "" {
			continue
		}

		if _, ok := palettes[tileSet.Palette]; !ok {
			return fmt.Errorf("palette %s does not exists", tileSet.Palette)
		}
	}

	for _, tileMap := range c.TileMaps {
		if tileMap.TileSet == "" && tileMap.Palette == "" {
			continue
		}

		if tileMap.TileSet != "" && tileMap.Palette != "" {
			return fmt.Errorf("tile set %s and palette %s can not both be set", tileMap.TileSet, tileMap.Palette)
		}

		if tileMap.TileSet != "" && tileMap.Transparent != "" {
			return fmt.Errorf("can not set transparent color when tile set %s is set", tileMap.TileSet)
		}

		err := validateColor(tileMap.Transparent)
		if err != nil {
			return fmt.Errorf("could not validate transparent color %s | %w", tileMap.Transparent, err)
		}

		if tileMap.TileSet != "" {
			tileSet, ok := tileSets[tileMap.TileSet]
			if !ok {
				return fmt.Errorf("tile set %s does not exists", tileMap.TileSet)
			}
			if tileSet.Size != "8x8" {
				return fmt.Errorf("tile sets used by tile maps must be 8x8 but it's actually %s", tileSet.Size)
			}
		}

		if tileMap.Palette != "" {
			_, ok := palettes[tileMap.Palette]
			if !ok {
				return fmt.Errorf("palette %s does not exist", tileMap.Palette)
			}
		}
	}

	err := validateColor(c.SetTransparent)
	if err != nil {
		return err
	}

	return nil
}

func validateColor(hex string) error {
	if hex == "" {
		return nil
	}

	_, err := ParseHexColor(hex)
	if err != nil {
		return err
	}

	return nil
}

// ParseHexColor parses a hex color into a valid gbacol.RGB15 color
func ParseHexColor(hex string) (*gbacol.RGB15, error) {
	c := color.RGBA{A: 0xFF}
	switch len(hex) {
	case 7:
		_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		if err != nil {
			return nil, fmt.Errorf("invaid color %s", err)
		}

	case 4:
		_, err := fmt.Sscanf(hex, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		if err != nil {
			return nil, fmt.Errorf("invalid color %s", err)
		}
		// Double the hex digits:
		c.R *= 0x11
		c.G *= 0x11
		c.B *= 0x11
	default:
		return nil, fmt.Errorf("invalid color len %d", len(hex))
	}

	ret := gbacol.NewRGB15(c)
	return &ret, nil
}
