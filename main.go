package main

import (
	"database/sql"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"log"

	chat "learning/go-portfolio/chat"
	db "learning/go-portfolio/db"
	todo "learning/go-portfolio/todo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
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

var Db *sql.DB
func main() {
	var err error
	e := echo.New()
	e.Renderer = newTemplate()

	// set the static folder
	e.Static("/static", "static")

	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/chat")
	})
	e.GET("/to-do", todo.To_do)
	e.DELETE("/to-do/:id", todo.Delete_item)
	e.PUT("/add-to-do", todo.Add_to_do)

	e.GET("/chat", chat.Chat)

	// middleware

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout})
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogError: true,
		LogURI:    true,
		LogStatus: true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info().
				Timestamp().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")
			} else {
				logger.Error().
				Timestamp().
				Err(v.Error).
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("err:")
			}
			return nil
		},
	}))


	Db, err = db.GetDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
	defer Db.Close()
}
