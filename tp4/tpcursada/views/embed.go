package views

import (
	"embed"
	"html/template"
)

//go:embed *.templ
var tmplFS embed.FS

// ParseTemplates parsea todas las plantillas embebidas y devuelve *template.Template
func ParseTemplates() *template.Template {
	return template.Must(template.ParseFS(tmplFS, "*.templ"))
}
