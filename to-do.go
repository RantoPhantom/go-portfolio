package main

import (
	"fmt"
	"net/http"
	"log"
	"github.com/labstack/echo/v4"
)

type item struct {
	Item_number int
	Item_content string
}

var item_list []item

func to_do(c echo.Context) error {
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


	return c.Render(http.StatusOK, "to-do.html", item_list)
}

func debug_error() error {
	var output error
	output = fmt.Errorf("asdadsfasdfasdf")
	return output
}

func add_to_do(c echo.Context) error {
	var item_input string = c.FormValue("item_content")

	if item_input != "" {
		i := item{Item_number: len(item_list)}
		i.Item_content = item_input
		item_list = append(item_list, i)
	}
	return c.Render(http.StatusCreated, "form", item_list)
}

func reset(c echo.Context) error {
	item_list = nil
	return c.Render(http.StatusNoContent, "form", item_list)
}
