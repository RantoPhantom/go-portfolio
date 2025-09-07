package handlers

import (
	"fmt"
	"learning/go-portfolio/database"
	"net/http"
	"strconv"
	"time"

	custom_middleware "learning/go-portfolio/middleware"

	"github.com/labstack/echo/v4"
)

const ITEMS_PER_PAGE int = 15

func TodoRouter(e *echo.Echo) {
	main_group := e.Group("/to-do")
	main_group.Use(custom_middleware.AuthMiddleware)
	main_group.GET("/", redirect_first_list)
	main_group.GET("/lists/", listUI)
	main_group.GET("/lists/:list_id/", todoUI)
	main_group.POST("/lists/:list_id/", itemsUI)
	main_group.DELETE("/lists/:list_id/", delete_list)
	main_group.GET("/lists/new/", add_list_ui)
	main_group.POST("/lists/", add_list)
	main_group.PATCH("/lists/:list_id/", change_list_name)
	main_group.PUT("/lists/:list_id/", add_to_do)
	main_group.DELETE("/lists/:list_id/:item_id/", delete_to_do)
}

func redirect_first_list(c echo.Context) error {
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}

	lists, err := db.Queries.Get_lists(c.Request().Context())
	if err != nil {
		return err
	}

	if len(lists) == 0 {
		id, err := db.Queries.Insert_list(c.Request().Context(), database.Insert_listParams{
			ListName:    "Today",
			IconColor:   "#00ff00",
			DateCreated: time.Now(),
		})
		if err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("./lists/%d/", id))
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("./lists/%d/", lists[0].ListID))
}

func user_retrieval_helper(c echo.Context) (*database.DB_Connection, error) {
	cookie, err := c.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	token := cookie.Value
	session, err := database.GetSessionFromToken(c.Request().Context(), token)
	if err != nil {
		return nil, err
	}

	db, err := database.GetDB(session.Username)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func todoUI(c echo.Context) error {
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}
	list, err := db.Queries.Get_list_info(c.Request().Context(), int64(list_id))
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "to-do", list)
}

type ItemsUIHydrate struct {
	Items     []database.TodoItem
	Next_page int
	Has_more  bool
}

func itemsUI(c echo.Context) error {
	has_more := true
	current_page_str := c.QueryParam("current_page")
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if current_page_str == "" {
		current_page_str = "0"
	}
	current_page, err := strconv.Atoi(current_page_str)
	if err != nil {
		return err
	}
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}
	items, err := db.Queries.Get_paginated_items(c.Request().Context(), database.Get_paginated_itemsParams{
		ListID: int64(list_id),
		Limit:  int64(ITEMS_PER_PAGE),
		Offset: int64(ITEMS_PER_PAGE * current_page),
	})
	if len(items) < ITEMS_PER_PAGE {
		has_more = false
	}
	var items_ui_hydrate *ItemsUIHydrate = &ItemsUIHydrate{
		Items:     items,
		Next_page: current_page + 1,
		Has_more:  has_more,
	}
	fmt.Println(has_more)
	return c.Render(http.StatusOK, "items", items_ui_hydrate)
}

func listUI(c echo.Context) error {
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}

	lists, err := db.Queries.Get_lists(c.Request().Context())
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "list", lists)
}

func add_list_ui(c echo.Context) error {
	return c.Render(http.StatusOK, "new_list", nil)
}

func add_list(c echo.Context) error {
	list_name := c.FormValue("listName")
	list_icon := c.FormValue("iconColor")
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}
	id, err := db.Queries.Insert_list(c.Request().Context(), database.Insert_listParams{
		ListName:    list_name,
		IconColor:   list_icon,
		DateCreated: time.Now(),
	})
	if err != nil {
		return err
	}
	c.Response().Header().Add("HX-Redirect", fmt.Sprintf("/to-do/lists/%d/", id))
	return c.NoContent(http.StatusNoContent)
}

func change_list_name(c echo.Context) error {
	// CHANGE THIS LATER
	new_name := "test"
	list_id_str := c.Param("list_id")
	list_id, err := strconv.Atoi(list_id_str)
	if err != nil {
		return err
	}
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}

	err = db.Queries.Rename_list(c.Request().Context(), database.Rename_listParams{
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
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}

	err = db.Queries.Remove_item(c.Request().Context(), database.Remove_itemParams{
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
	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}
	err = db.Queries.Remove_list(c.Request().Context(), int64(list_id))
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

	db, err := user_retrieval_helper(c)
	if err != nil {
		return err
	}
	err = db.Queries.Insert_item(
		c.Request().Context(),
		database.Insert_itemParams{
			ListID:      int64(list_id),
			Content:     content,
			DateCreated: time.Now(),
		},
	)
	if err != nil {
		return err
	}
	return itemsUI(c)
}
