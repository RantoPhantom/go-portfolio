package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	middleware "learning/go-portfolio/middleware"
	handlers "learning/go-portfolio/handlers"
)

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	main_group := e.Group("/to-do")
	main_group.Use(middleware.AuthMiddleware)


	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/login")
	})
	e.GET("/login", func (c echo.Context) error {
		return c.Render(http.StatusOK, "form.html", nil)
	})
	// set the static folder
	e.Static("/static", "static")
	main_group.GET("/", handlers.Todo)
	main_group.DELETE("/:id", handlers.Delete_item)
	main_group.PUT("/add-to-do", handlers.Add_to_do)

	// middleware
	middleware.SetupLogger(e)
	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
}
