package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"log"

	"github.com/labstack/echo/v4"
	middleware "learning/go-portfolio/middleware"
	handlers "learning/go-portfolio/handlers"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	if !strings.Contains(name, ".html"){
		name += ".html"
	}
	return t.templates.ExecuteTemplate(writer, name, data)
}

func newTemplate() *Template {
	templates := template.New("")
	err := filepath.Walk("./views", func(path string, info os.FileInfo, err error) error{
		if strings.Contains(path, ".html"){
			_, err := templates.ParseFiles(path)
			if err != nil{
				log.Println(err)
			}
		}
		return err
	})

	if err != nil {
		panic(err)
	}
	return &Template{
		templates : templates,
	}
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	// set the static folder
	e.Static("/static", "static")

	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/to-do")
	})
	e.GET("/to-do", handlers.Todo)
	e.DELETE("/to-do/:id", handlers.Delete_item)
	e.PUT("/to-do/add-to-do", handlers.Add_to_do)

	// middleware
	middleware.SetupLogger(e)
	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
}
