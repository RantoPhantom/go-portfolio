package middleware

import (
	"fmt"
	"learning/go-portfolio/custom_errors"
	"learning/go-portfolio/database"
	"net/http"
	"errors"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("username")
		if cookie == nil {
			return c.Redirect(http.StatusSeeOther,"/auth/login")
		}

		// check is user is in db, if not then kick them out
		if _,err := database.GetDB(cookie.Value); errors.Is(err, custom_errors.UserNotFound){
			cookie.MaxAge = -1
			return c.Redirect(http.StatusSeeOther,"/auth/login")
		}

		if err != nil {
			return echo.NewHTTPError(
				http.StatusInternalServerError,
				fmt.Errorf("something went wrong! %s", err.Error()),
			)
		}

		return next(c)
	}
}
