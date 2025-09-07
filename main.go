package main

import (
	database "learning/go-portfolio/database"
	handlers "learning/go-portfolio/handlers"
	custom_middleware "learning/go-portfolio/middleware"
	utils "learning/go-portfolio/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	var err error

	defer database.CloseSessionDb()

	utils.LoadConfig("./.env.local")

	e := echo.New()
	e.Renderer = utils.NewTemplate()
	err = database.CreateSessionDB()
	if err != nil {
		panic(err)
	}

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
