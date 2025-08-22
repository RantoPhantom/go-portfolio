package main

import (
	handlers "learning/go-portfolio/handlers"
	custom_middleware "learning/go-portfolio/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Renderer = newTemplate()

	// set the static folder
	e.Static("/static", "static")

	// middleware
	custom_middleware.SetupLogger(e)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/information/")
	})

	// routers
	handlers.AuthRouter(e)
	handlers.TodoRouter(e)
	handlers.InformationRouter(e)

	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
}
