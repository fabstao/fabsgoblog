package views

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// ComunesSitio - common rendering values
type ComunesSitio struct {
	Title  string
	Header string
	Footer string
}

// Comunes - compartir info común a todas las páginas
var Comunes = ComunesSitio{
	Title:  "FabsBlog",
	Header: "Fabs Blog",
	Footer: "By Fabs",
}

// Template struct for views
type Template struct {
	Templates *template.Template
}

// Render Templates
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
