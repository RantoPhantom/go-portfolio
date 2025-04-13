package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("username")
		if cookie == nil {
			cookie = new(http.Cookie)
			cookie.Name = "username"
			cookie.Value = "test"
			cookie.Expires = time.Now().Add(12 * time.Hour)
			cookie.HttpOnly = true
			c.SetCookie(cookie)
		}
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		return next(c)
	}
}
