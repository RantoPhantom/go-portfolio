package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/labstack/echo/v4"
)

type item struct {
	Item_number int
	Item_content string
}

var item_list []item
var max_item_id int
func to_do(c echo.Context) error {
	fetch_todo_db()
	return c.Render(http.StatusOK, "to-do.html", item_list)
}

func fetch_todo_db() {
	item_list = nil
	rows, query_err := db.Query("SELECT * FROM todo;")
	if query_err != nil {
		log.Print(query_err)
	}

	for rows.Next() {
		var id int
		var content string
		var item item

		rows.Scan(&id, &content)
		item.Item_number = id
		item.Item_content = content
		item_list = append(item_list, item)
	}
	query_err = db.QueryRow("SELECT ID FROM todo ORDER BY id DESC;").Scan(&max_item_id)
}

func debug_error() error {
	var output error
	output = fmt.Errorf("asdadsfasdfasdf")
	return output
}

func add_to_do(c echo.Context) error {
	var item_input string = c.FormValue("item_content")

	if item_input != "" {
		max_item_id += 1
		i := item{Item_number: max_item_id}
		i.Item_content = item_input

		query := fmt.Sprintf("insert into todo (content) values ('%s')", item_input)
		_, err := db.Exec(query)
		if err != nil {
			log.Print(err)
			return c.JSON(http.StatusNoContent, "cannot insert into db")
		}
		item_list = append(item_list, i)
	}
	return c.Render(http.StatusCreated, "form", item_list)
}

func reset(c echo.Context) error {
	item_list = nil
	return c.Render(http.StatusNoContent, "form", item_list)
}
