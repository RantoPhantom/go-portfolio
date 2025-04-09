package main

import (
	"context"
	"strconv"
	"learning/go-portfolio/database"
	"net/http"
	"github.com/labstack/echo/v4"
	"time"
)

var db *database.UserDb
func init_db() error {
	var username string = "test"
	var err error
	if err = database.CreateDB(username); err != nil { return err }
	if db, err = database.GetDB(username); err != nil { return err }
	return nil
}

var item_list []database.TodoItem

func todo(c echo.Context) error {
	if err := init_db(); err != nil { return err }
	if err := fetch_todo_db(); err != nil { return err }

	return c.Render(http.StatusOK, "to-do.html", item_list)
}

func fetch_todo_db() error {
	ctx := context.Background()
	var err error
	item_list, err = db.Queries.Get_all_items(ctx)
	if err != nil { return err }
	return nil
}

func add_to_do(c echo.Context) error {
	content := c.FormValue("item_content")
	ctx := context.Background()
	if err := db.Queries.Add_item(ctx, database.Add_itemParams{
		Content: content,
		DateCreated: time.Now(),
	}); err != nil { return err }
	if err := fetch_todo_db(); err != nil { return err }
	return c.Render(http.StatusOK, "form.html", item_list)
}

func delete_item(c echo.Context) error {
	var item_id int64
	var err error
	if item_id, err = strconv.ParseInt(c.Param("id"),10,64); err != nil { return err }
	ctx := context.Background()
	if err:= db.Queries.Delete_item(ctx, item_id); err != nil { return err }
	if err := fetch_todo_db(); err != nil { return err }
	return c.Render(http.StatusOK, "form.html", item_list)
}
