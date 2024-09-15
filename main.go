package main

import (
	"io"
	"os"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"html/template"
	"github.com/rs/zerolog"
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
	// set the static folder
	e.Static("/static", "static")
	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/to-do")
	})

	e.GET("/to-do", to_do)
	e.POST("/reset-list", reset)
	e.POST("/add-to-do", add_to_do)
	e.POST("/get-to-do", get_to_do)

	// middleware

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
			Str("URI", v.URI).
			Int("status", v.Status).
			Msg("request")

			return nil
		},
	}))	

	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
}
