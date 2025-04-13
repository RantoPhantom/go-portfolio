package handlers

import (
	"context"
	"learning/go-portfolio/custom_errors"
	"learning/go-portfolio/database"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// init db based on echo context
func init_db(c echo.Context) (*database.UserDb, error) {
	var db *database.UserDb
	cookie, err := c.Cookie("username")
	if err != nil { return nil, err }
	username := cookie.Value

	db, err = database.GetDB(username)
	if err != nil { return nil, err }
	return db, nil
}

var item_list []database.TodoItem

func Todo(c echo.Context) error {
	db, err := init_db(c)
	if err != nil { return err }
	if err := fetch_all_todo(db); err != nil { return err }

	return c.Render(http.StatusOK, "to-do.html", item_list)
}

func fetch_all_todo(db *database.UserDb) error {
	ctx := context.Background()
	var err error
	item_list, err = db.Queries.Get_all_items(ctx)
	if err != nil { return err }
	return nil
}


func Add_to_do(c echo.Context) error {
	db, err := init_db(c)
	if err != nil { return err }

	ctx := c.Request().Context()
	content := c.FormValue("item_content")
	// return an error if content is empty
	if content == "" {
		c.Response().Header().Set("HX-Reswap", "innerHtml")
		c.Response().Header().Set("HX-Retarget", "#error")
		return c.HTML(http.StatusOK, custom_errors.InvalidInput.Error())
	}
	// add to db
	if err := db.Queries.Add_item(ctx, database.Add_itemParams{
		Content: content,
		DateCreated: time.Now(),
	}); err != nil { return err }
	// refetch
	if err := fetch_all_todo(db); err != nil { return err }
	return c.Render(http.StatusOK, "form.html", item_list)
}

func Delete_item(c echo.Context) error {
	db, err := init_db(c)
	if err != nil { return err }

	ctx := c.Request().Context()
	item_id, err := strconv.ParseInt(c.Param("id"),10,64); if err != nil { return err }
	if err:= db.Queries.Delete_item(ctx, item_id); err != nil { return err }

	if err := fetch_all_todo(db); err != nil { return err }
	return c.Render(http.StatusOK, "form.html", item_list)
}
