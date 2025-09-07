CREATE TABLE IF NOT EXISTS todo_items (
	item_id INTEGER PRIMARY KEY AUTOINCREMENT,
	list_id INTEGER NOT NULL,
	content TEXT NOT NULL,
	date_created timestamp NOT NULL,
	is_done BOOL DEFAULT FALSE NOT NULL,
	foreign key(list_id) references lists(list_id)
);

CREATE TABLE IF NOT EXISTS user_info (
	password_hash TEXT NOT NULL,
	date_created timestamp NOT NULL
);

CREATE TABLE IF NOT EXISTS lists (
	list_id INTEGER PRIMARY KEY AUTOINCREMENT,
	list_name TEXT NOT NULL,
	icon_color TEXT NOT NULL,
	date_created timestamp NOT NULL
);
-- this is here because sqlc cannot handle multi schema for now
-- this should be an empty table
-- this should be synced with the sessions_db.sql 
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL,
	token TEXT NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	date_created TIMESTAMP NOT NULL
);
pragma journal_mode = wal
;

