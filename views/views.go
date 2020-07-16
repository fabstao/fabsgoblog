package views

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

// Template struct for views
type Template struct {
	Templates *template.Template
}

// Render Templates
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.Templates.ExecuteTemplate(w, name, data)
}
