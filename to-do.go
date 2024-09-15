package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
)

type item struct {
	Item_number int
	Item_content string
}

var item_list []item

func to_do(c echo.Context) error {
	return c.Render(http.StatusOK, "to-do.html", item_list)
}

func get_to_do(c echo.Context) error {
	html_response := create_item_list_response()
	return c.HTML(http.StatusOK, html_response)
}

func add_to_do(c echo.Context) error {
	var item_input string = c.FormValue("item_content")

	if item_input != "" {
		i := item{Item_number: len(item_list)}
		i.Item_content = item_input
		item_list = append(item_list, i)
	}
	c.Response().Header().Add("HX-Refresh", "true")
	return c.NoContent(http.StatusOK)
}

func reset(c echo.Context) error {
	item_list = nil
	return c.NoContent(http.StatusNoContent)
}

func create_item_list_response() string {
	var html_response string = "<div id='item_list' hx-swap-oob='#item_list'>"

	for i:=0; i < len(item_list); i++{
		html_response += fmt.Sprintf("<ul id='item'>%d,%s</ul>", item_list[i].Item_number, item_list[i].Item_content)
	}
	html_response += "</div>"
	return html_response
}
