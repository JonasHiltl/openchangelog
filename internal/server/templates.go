package server

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	t := template.Must(template.ParseGlob("web/views/*.html"))
	template.Must(t.ParseGlob("web/views/**/*.html"))

	return &Templates{
		templates: t,
	}
}
