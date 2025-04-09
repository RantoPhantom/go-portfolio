-- name: Get_all_items :many
SELECT * FROM todo_items;

-- name: Delete_item :exec
DELETE FROM todo_items
WHERE id = ?;

-- name: Update_item :exec
UPDATE todo_items
SET is_done = ?;

-- name: Add_item :exec
INSERT INTO todo_items(
	content,
	item_number,
	date_created
)
VALUES (?,?,?);

-- name: Get_item_count :one
SELECT CAST(IFNULL(MAX(id), 0) as INTEGER) as COUNT FROM todo_items;
