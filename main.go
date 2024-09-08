package main

import (
	"io"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	return t.templates.ExecuteTemplate(writer, name, data)
}
func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
}

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	e.Static("/static", "static")
	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/to-do")
	})

	e.GET("/to-do", to_do)
	e.POST("/to-do", to_do)
	e.GET("/reset-list", reset)

	// middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Logger.Fatal(e.Start(":6969"))
}
