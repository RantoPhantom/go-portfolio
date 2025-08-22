package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func InformationRouter(e *echo.Echo) {
	main_group := e.Group("/information")
	main_group.GET("/", indexUI)
}

func indexUI(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}
