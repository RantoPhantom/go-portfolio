package handlers

import (
	"learning/go-portfolio/database"
	"net/http"
	"strconv"
	"time"

	custom_middleware "learning/go-portfolio/middleware"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
)

func TodoRouter(e *echo.Echo) {
	main_group := e.Group("/to-do")
	main_group.Use(custom_middleware.AuthMiddleware)
	main_group.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "./lists/0")
	})
	main_group.GET("/lists/", listUI)
	main_group.GET("/lists/:list_id", todoUI)
	main_group.DELETE("/lists/:list_id", delete_list)
	main_group.PUT("/lists/:list_id", add_list)
	main_group.PATCH("/lists/:list_id", change_list_name)
	main_group.POST("/lists/:list_id", add_to_do)
	main_group.DELETE("/lists/:list_id/:item_id", delete_to_do)
}

type TodoUIHydrate struct {
	Items []database.TodoItem
	List  database.List
}

func todoUI(c echo.Context) error {
	list_id_str := c.Param("list_id")
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	username := cookie.Value

	list_id, err := strconv.Atoi(list_id_str)
	list, err := get_list_info(username, int64(list_id))
	items, err := get_items(username, int64(list_id))

	hydration := &TodoUIHydrate{
		Items: items,
		List:  *list,
	}
	return c.Render(http.StatusOK, "to-do", hydration)
}

func listUI(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}

	username := cookie.Value
	db, err := database.GetDB(username)
	if err != nil {
		return err
	}

	lists, err := db.Queries.Get_lists(context.Background())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "list", lists)
}

func add_list(c echo.Context) error {
	return nil
}

func change_list_name(c echo.Context) error {
	// CHANGE THIS LATER
	new_name := "test"
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	username := cookie.Value
	db, err := database.GetDB(username)
	if err != nil {
		return err
	}

	err = db.Queries.Rename_list(context.Background(), database.Rename_listParams{
		ListName: new_name,
		ListID:   int64(list_id),
	})
	return c.NoContent(http.StatusNoContent)
}

func delete_to_do(c echo.Context) error {
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}
	item_id_str := c.Param("item_id")
	item_id, err := strconv.Atoi(item_id_str)
	if err != nil {
		return err
	}
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	username := cookie.Value
	db, err := database.GetDB(username)
	if err != nil {
		return err
	}

	err = db.Queries.Remove_item(context.Background(), database.Remove_itemParams{
		ListID: int64(list_id),
		ItemID: int64(item_id),
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func delete_list(c echo.Context) error {
	var list_id_str string = c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	username := cookie.Value
	db, err := database.GetDB(username)
	if err != nil {
		return err
	}
	err = db.Queries.Remove_list(context.Background(), int64(list_id))
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func add_to_do(c echo.Context) error {
	content := c.FormValue("item_content")
	if content == "" {
		c.Response().Header().Add("HX-Retarget", "#input_error")
		c.Response().Header().Add("HX-Reswap", "innerHTML")
		return c.HTML(http.StatusBadRequest, "the input must not be empty")
	}
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}

	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}
	username := cookie.Value
	db, err := database.GetDB(username)
	if err != nil {
		return err
	}
	err = db.Queries.Insert_item(
		context.Background(),
		database.Insert_itemParams{
			ListID:      int64(list_id),
			Content:     content,
			DateCreated: time.Now(),
		},
	)
	if err != nil {
		return err
	}
	items, err := get_items(username, int64(list_id))
	return c.Render(http.StatusOK, "form", items)
}

func get_items(username string, list_id int64) ([]database.TodoItem, error) {
	db, err := database.GetDB(username)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	items, err := db.Queries.Get_items(ctx, list_id)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func get_list_info(username string, list_id int64) (*database.List, error) {
	db, err := database.GetDB(username)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	list, err := db.Queries.Get_list_info(ctx, list_id)
	return &list, nil
}
