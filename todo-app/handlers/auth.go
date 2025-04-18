package handlers

import (
	"errors"
	"learning/go-portfolio/custom_errors"
	"learning/go-portfolio/database"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)


func LoginUI(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}
func SignupUI(c echo.Context) error {
	return c.Render(http.StatusOK, "signup.html", nil)
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	db, err := database.GetDB(username)
	if errors.Is(err, custom_errors.UserNotFound){
		c.Response().Header().Add("HX-Retarget", "#username_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, err.Error())
	} else if err != nil { return err }

	ctx := c.Request().Context()
	password_hash, err := db.Queries.Get_password(ctx)
	if err != nil { return err }
	if err := bcrypt.CompareHashAndPassword([]byte(password_hash), []byte(password)); err != nil {
		c.Response().Header().Add("HX-Retarget", "#password_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, custom_errors.InvalidCredentials.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = username
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(12 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	c.Response().Header().Add("HX-Redirect", "/to-do/")
	return c.NoContent(http.StatusOK)
}

func Logout(c echo.Context) error {
	old_cookie, err := c.Cookie("username")
	if err != nil { return err }
	old_cookie.Value = "asdf"
	old_cookie.Path = "/"
	old_cookie.MaxAge = -1
	c.SetCookie(old_cookie)
	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func Signup(c echo.Context) error {
	username := c.FormValue("username")
	password := []byte(c.FormValue("password"))


	// bcrypt only accepts 72 bytes and lower
	if len(password) > 71 {
		c.Response().Header().Add("HX-Retarget", "#password_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, custom_errors.PasswordTooLong.Error())
	}

	err := database.CheckUserExists(username)
	if errors.Is(err, custom_errors.UserDbExists){
		c.Response().Header().Add("HX-Retarget", "#username_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, err.Error())
	}
	password_hash, err := bcrypt.GenerateFromPassword(password, 4)
	if err != nil { return err }

	// create the db
	err = database.CreateDB(username, string(password_hash))
	if errors.Is(err, custom_errors.InvalidUsername) {
		c.Response().Header().Add("HX-Retarget", "#username_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, err.Error())
	} else if err != nil { return err }

	c.Response().Header().Add("HX-Redirect", "./login")
	return c.NoContent(http.StatusOK)
}
