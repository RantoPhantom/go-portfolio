package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type item struct {
	Item_number int
	Item_content string
	Date_created string
	Is_Done bool
}

var item_list []item
var max_item_id int

func to_do(c echo.Context) error {
	err := fetch_todo_db()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "to-do.html", item_list)

}

func fetch_todo_db() error {
	item_list = nil
	rows, query_err := db.Query("SELECT id, todo_content, datetime(date_created), is_done FROM todo;")
	if query_err != nil {
		log.Print(query_err)
		return query_err
	}

	for rows.Next() {
		var item item
		rows.Scan(&item.Item_number, &item.Item_content, &item.Date_created, &item.Is_Done)
		item_list = append(item_list, item)
	}

	query_err = db.QueryRow("SELECT COUNT(*) FROM todo").Scan(&max_item_id)
	if query_err != nil {
		return query_err
	}

	return nil
}

func delete_item(c echo.Context) error {
	id := c.Param("id")

	query := fmt.Sprintf("DELETE FROM todo WHERE id=%s",id)
	_,err := db.Exec(query)
	if err != nil {
		log.Print(err)
		return err
	}

	return c.HTML(http.StatusOK, "item deleted")
}

func add_to_do(c echo.Context) error {
	var item_input string = c.FormValue("item_content")

	max_item_id += 1
	i := item{Item_number: max_item_id}
	i.Item_content = strings.TrimSpace(item_input)
	i.Date_created = time.Now().Format(time.RFC3339)

	query := fmt.Sprintf(`INSERT INTO todo (id, todo_content, date_created)
	VALUES (%d, '%s', '%s');`,
	i.Item_number, 
	i.Item_content, 
	i.Date_created,)

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	err = fetch_todo_db()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "form", item_list)
}
