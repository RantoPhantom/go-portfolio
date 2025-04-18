package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
	custom_middleware "learning/go-portfolio/middleware"
	handlers "learning/go-portfolio/handlers"
)

func main() {
	e := echo.New()
	e.Renderer = newTemplate()
	main_group := e.Group("/to-do")
	main_group.Use(custom_middleware.AuthMiddleware)

	auth_group := e.Group("/auth")
	auth_group.GET("/login", handlers.LoginUI)
	auth_group.GET("/signup", handlers.SignupUI)
	auth_group.GET("/logout", handlers.Logout)
	auth_group.POST("/login-request", handlers.Login)
	auth_group.POST("/signup-request", handlers.Signup)

	// router
	e.GET("/", func (c echo.Context) error {
		return c.Redirect(http.StatusSeeOther, "/to-do/")
	})
	// set the static folder
	e.Static("/static", "static")
	main_group.GET("/", handlers.Todo)
	main_group.DELETE("/:id", handlers.Delete_item)
	main_group.PUT("/add-to-do", handlers.Add_to_do)

	// middleware
	custom_middleware.SetupLogger(e)
	e.Debug = true
	e.Logger.Fatal(e.Start(":6969"))
}
