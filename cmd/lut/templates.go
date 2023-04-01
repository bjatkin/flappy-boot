package main

import (
	"embed"
	"fmt"
	"math"
	"strings"
	"text/template"
)

//go:embed templates
var templates embed.FS

var goTemplates = template.Must(
	template.New("go_templates").
		Funcs(map[string]any{
			"floatRef":   floatRef,
			"hex":        hex,
			"fixToFloat": fixToFloat,
			"floatToFix": floatToFix,
		}).
		ParseFS(templates, "templates/*.tmpl"),
)

func floatRef(i int) string {
	f := float64(i<<2) / float64(1<<8)
	return fmt.Sprintf("%04f", f)
}

func hex(i int) string {
	s := fmt.Sprintf("%04X", i)
	prefix := "0x"
	if strings.HasPrefix(s, "-") {
		prefix = "-0x0"
		s = strings.TrimLeft(s, "-")
	}

	return prefix + s
}

func floatToFix(f float64) int {
	conv := math.Round(f * float64(1<<8))
	return int(conv)
}

func fixToFloat(f int) float64 {
	return float64(f) / float64(1<<8)
}
