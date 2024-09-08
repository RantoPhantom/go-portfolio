package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

var item_content []string

func to_do(c echo.Context) error {
	item_content = append(item_content, c.FormValue("item_content"))
	return c.Render(http.StatusOK, "to-do.html", item_content)
}

func reset(c echo.Context) error {
	item_content = nil
	return c.String(http.StatusAccepted, "reset done!")
}
