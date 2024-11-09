package chat

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func Chat(c echo.Context) error{
	return c.Render(http.StatusOK, "chat", nil)
}

