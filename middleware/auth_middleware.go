package middleware

import (
	"errors"
	"fmt"
	"learning/go-portfolio/custom_errors"
	"learning/go-portfolio/database"
	"net/http"

	"github.com/labstack/echo/v4"
	utils "learning/go-portfolio/utils"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(utils.AppConfig.SESSION_COOKIE_NAME)
		if cookie == nil {
			return c.Redirect(http.StatusSeeOther, "/auth/login")
		}

		// checks if session is valid then kicks if not
		if _, err := database.GetSessionFromToken(c.Request().Context(), cookie.Value); errors.Is(err, custom_errors.SessionExpired) {
			cookie.MaxAge = -1
			return c.Redirect(http.StatusSeeOther, "/auth/login")
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
