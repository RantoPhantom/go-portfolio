-- -----items-------
-- name: Get_paginated_items :many
select *
from todo_items
where list_id = ?
order by item_id asc
limit ?
offset ?
;

-- name: Get_items :many
select *
from todo_items
where list_id = ?
;

-- name: Insert_item :exec
INSERT INTO todo_items(list_id, content, date_created)
VALUES(?,?,?);

-- name: Remove_item :exec
delete from todo_items
where item_id = ? and list_id = ?
;

-- -----lists-------
-- name: Get_lists :many
select *
from lists
;

-- name: Get_list_info :one
select *
from lists
where list_id = ?
;

-- name: Insert_list :one
insert into lists(list_name, icon_color, date_created)
values (?,?,?)
returning list_id;

-- name: Remove_list :exec
delete from lists
where list_id = ?
;

-- name: Rename_list :exec
update lists
set list_name=?
where list_id=?;

-- ----user_info-----
-- name: Get_password :one
select password_hash
from user_info
;

-- name: Insert_user_info :exec
insert into user_info(password_hash, date_created)
values(?,?);

