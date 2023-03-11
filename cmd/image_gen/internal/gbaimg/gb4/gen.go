package gb4

import (
	"embed"
	"fmt"
	"image"
	"os"
	"text/template"
)

//go:embed templates
var templates embed.FS

// Palette is a named palette
type Palette struct {
	Name string
	File string

	img image.Image
}

var paletteTemplate = template.Must(template.ParseFS(templates))

func GeneratePalettes(palette Palette) error {
	p, err := NewPal16(palette.img, nil)
	if err != nil {
		return fmt.Errorf("failed to create new 16 color palette %s", err)
	}

	err = os.WriteFile(palette.Name+"Palette.pal4", rawPalette(p), 0o0666)
	if err != nil {
		return fmt.Errorf("failed to write %s", err)
	}

	f, err := os.Create(palette.Name + "tilemap.go")
	if err != nil {
		return fmt.Errorf("failed to create tilemap.go file %s", err)
	}

	err = paletteTemplate.ExecuteTemplate(f, "tilemap.go.tmpl", palette)
	if err != nil {
		return fmt.Errorf("")
	}
}
