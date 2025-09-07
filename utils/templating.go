package utils

import (
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(writer io.Writer, name string, data any, context echo.Context) error {
	if !strings.Contains(name, ".html") {
		name += ".html"
	}
	return t.templates.ExecuteTemplate(writer, name, data)
}

func NewTemplate() *Template {
	templates := template.New("")
	err := filepath.Walk("./views", func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".html") {
			_, err := templates.ParseFiles(path)
			if err != nil {
				log.Println(err)
			}
		}
		return err
	})
	if err != nil {
		panic(err)
	}
	return &Template{
		templates: templates,
	}
}
