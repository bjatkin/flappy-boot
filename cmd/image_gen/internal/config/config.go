package config

import (
	"fmt"
	"image/color"
	"os"

	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/gbacol"
	"github.com/bjatkin/flappy_boot/cmd/image_gen/internal/gbaimg/tile"
	"gopkg.in/yaml.v2"
)

// Config is a config file for the command
type Config struct {
	Palettes       []Palette `yaml:"Paletts"`
	TileSets       []TileSet `yaml:"TileSets"`
	TileMaps       []TileMap `yaml:"TileMaps"`
	SetTransparent string    `yaml:"SetTransparent"`
}

// Palette is a named palette
type Palette struct {
	Name        string `yaml:"Name"`
	File        string `yaml:"File"`
	Description string `yaml:"Description"`
}

// TileSet is a named tile set
type TileSet struct {
	Name    string `yaml:"Name"`
	File    string `yaml:"File"`
	Palette string `yaml:"Palette"`
	Size    string `yaml:"Size"`
}

// TileMap is a named tile map
type TileMap struct {
	Name    string `yaml:"Name"`
	File    string `yaml:"File"`
	TileSet string `yaml:"TileSet"`
	Palette string `yaml:"Palette"`
}

// NewConfigFromFile reads in the yaml file at the provided file location and then marshalls it into a new config struct
func NewConfigFromFile(file string) (*Config, error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	fmt.Println("raw", string(raw))

	config := &Config{}
	err = yaml.Unmarshal(raw, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Validate ensure the configuration file has a valid format
func (c *Config) Validate() error {
	palettes := make(map[string]struct{})
	for _, pal := range c.Palettes {
		palettes[pal.Name] = struct{}{}
	}

	tileSets := make(map[string]struct{})
	for _, tileSet := range c.TileSets {
		tileSets[tileSet.Name] = struct{}{}
		if _, ok := palettes[tileSet.Palette]; !ok {
			return fmt.Errorf("palette %s does not exists", tileSet.Palette)
		}
		_, err := tile.NewSize(tileSet.Size)
		if err != nil {
			return err
		}
	}

	for _, tileMap := range c.TileMaps {
		if _, ok := tileSets[tileMap.TileSet]; !ok {
			return fmt.Errorf("tile set %s does not exists", tileMap.TileSet)
		}
	}

	_, err := ParseHexColor(c.SetTransparent)
	if err != nil {
		return err
	}

	return nil
}

// ParseHexColor parses a hex color into a valid gbacol.RGB15 color
func ParseHexColor(s string) (gbacol.RGB15, error) {
	c := color.RGBA{A: 0xFF}
	switch len(s) {
	case 7:
		_, err := fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
		if err != nil {
			return 0, fmt.Errorf("invaid color %s", err)
		}

	case 4:
		_, err := fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		if err != nil {
			return 0, fmt.Errorf("invalid color %s", err)
		}
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		return 0, fmt.Errorf("invalid color len %d", len(s))
	}

	return gbacol.NewRGB15(c), nil
}
