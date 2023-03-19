package generate

import (
	"embed"
	"strings"
	"text/template"
	"unicode"
)

//go:embed templates
var templates embed.FS

var goTemplates = template.Must(
	template.New("go_templates").
		Funcs(
			map[string]any{
				"private": private,
				"public":  public,
				"add":     add,
			}).
		ParseFS(templates, "templates/*.tmpl"),
)

func public(name string) string {
	runes := []rune(name)
	if unicode.IsLower(runes[0]) {
		return strings.ToUpper(string(runes[0])) + string(runes[1:])
	}

	return name
}

func private(name string) string {
	runes := []rune(name)
	if unicode.IsUpper(runes[0]) {
		return strings.ToLower(string(runes[0])) + string(runes[1:])
	}

	return name
}

func add(args ...int) int {
	var total int
	for _, a := range args {
		total += a
	}
	return total
}
